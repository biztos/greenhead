// Package api defines the Greenhead HTTP API.
package api

import (
	"fmt"
	"html"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/assets"
	"github.com/biztos/greenhead/rgxp"
	"github.com/biztos/greenhead/version"
)

// Role defines a set of permissions for API Keys.
type Role struct {
	Name      string               `toml:"name"`   // Name of role.
	Endpoints []*rgxp.OptionalRgxp `toml:"paths"`  // Endpoint access.
	Agents    []*rgxp.OptionalRgxp `toml:"agents"` // Agents access.
}

// Key defines an API Key that is attached to a role.
type Key struct {
	AuthKey string   `toml:"auth_key"` // Key string for client auth.
	Name    string   `toml:"name"`     // Name of the key user for logs.
	Roles   []string `toml:"roles"`    // Roles available.
}

// Config defines the configuration of the API.
//
// It is included as part of the Runner config.
//
// Fields beginning with "App" are passed to the Fiber app.
type Config struct {

	// General configuration:
	ListenAddress    string `toml:"listen_address"` // Address for serving, e.g. ":3000"
	UnstructuredLogs bool   `toml:"default_logger"` // Use default Fiber logger instead of structured.

	// Access control:
	Roles      []*Role `toml:"roles"`       // Roles defining access.
	Keys       []*Key  `toml:"keys"`        // Keys mapping to roles by name.
	AccessFile string  `toml:"access_file"` // TOML file for (more) Roles and Keys.
	NoKeys     bool    `toml:"no_keys"`     // DO NOT require API keys.

	// Fiber app config; see Fiber docs for specifics:
	AppPrefork                 bool          `toml:"app_prefork"`
	AppServerHeader            string        `toml:"app_server_header"`
	AppBodyLimit               int           `toml:"app_body_limit"`
	AppConcurrency             int           `toml:"app_concurrency"`
	AppReadTimeout             time.Duration `toml:"app_read_timeout"`
	AppWriteTimeout            time.Duration `toml:"app_write_timeout"`
	AppIdleTimeout             time.Duration `toml:"app_idle_timeout"`
	AppProxyHeader             string        `toml:"app_proxy_header"`
	AppDisableStartupMessage   bool          `toml:"app_disable_startup_message"`
	AppEnableTrustedProxyCheck bool          `toml:"app_enable_trusted_proxy_check"`
	AppTrustedProxies          []string      `toml:"app_trusted_proxies"`
}

// API represents an instance of the Greenhead HTTP API.
type API struct {
	config       *Config
	ident        string
	app          *fiber.App
	sourceAgents map[string]*agent.Agent
	activeAgents map[string]*agent.Agent
	agentMutex   *sync.Mutex
}

// NewAPI creates an API instance.
func NewAPI(cfg *Config, agents []*agent.Agent) (*API, error) {

	// TODO: consider having no agents, only what you define on the API.
	// Seems like a bad idea to not have at least one available agent though.
	if len(agents) == 0 {
		return nil, fmt.Errorf("at least one agent must be defined")
	}

	// Some app configs can be set by the user, to tune the server for their
	// needs.
	// https://pkg.go.dev/github.com/gofiber/fiber/v2#Config
	ident := fmt.Sprintf("GREENHEAD %s HTTP API %s",
		version.Version, version.ApiVersion)
	fiber_cfg := fiber.Config{
		AppName:                 ident,
		Prefork:                 cfg.AppPrefork,
		ServerHeader:            cfg.AppServerHeader,
		BodyLimit:               cfg.AppBodyLimit,
		Concurrency:             cfg.AppConcurrency,
		ReadTimeout:             cfg.AppReadTimeout,
		WriteTimeout:            cfg.AppWriteTimeout,
		IdleTimeout:             cfg.AppIdleTimeout,
		ProxyHeader:             cfg.AppProxyHeader,
		DisableStartupMessage:   cfg.AppDisableStartupMessage,
		EnableTrustedProxyCheck: cfg.AppEnableTrustedProxyCheck,
		TrustedProxies:          cfg.AppTrustedProxies,
		// EnablePrintRoutes ?
	}
	app := fiber.New(fiber_cfg)
	sourceAgents := map[string]*agent.Agent{}
	for _, a := range agents {
		if sourceAgents[a.Name] != nil {
			return nil, fmt.Errorf("duplicate agent by name: %q", a.Name)
		}
		sourceAgents[a.Name] = a
	}
	// TODO: logging setup!
	api := &API{
		ident:        ident,
		config:       cfg,
		app:          app,
		sourceAgents: sourceAgents,
		activeAgents: map[string]*agent.Agent{},
		agentMutex:   &sync.Mutex{},
	}
	api.setRoutes()
	return api, nil

}

type RequestPayloadAgent struct {
	Agent string `json:"agent"`
}

type RequestPayloadChat struct {
	Prompt string `json:"prompt"`
}

type RequestPayloadRemove struct {
	Reason string `json:"reason"`
}

// Set up the routing for the Fiber app, with access to the agents et al.
// Middleware will handle the auth and logging.
func (api *API) setRoutes() {

	api.app.Get("/", func(c *fiber.Ctx) error {
		return api.HandleRoot(c)
	})

	api.app.Get("/v1/agents/list", func(c *fiber.Ctx) error {
		return api.HandleAgentsList(c)
	})

	api.app.Post("/v1/agents/new", func(c *fiber.Ctx) error {
		return api.HandleAgentsNew(c)
	})

	api.app.Post("/v1/agents/:id/chat", func(c *fiber.Ctx) error {
		return api.HandleAgentsChat(c)
	})

	api.app.Get("/v1/ui", func(c *fiber.Ctx) error {
		return api.HandleUI(c)
	})

	// ok -- want to do the UI as / if accept is set to HTML.
	api.app.Post("/v1/ui", func(c *fiber.Ctx) error {
		return api.HandleUI(c)
	})

	// Serve a favicon because the requests are annoying.
	api.app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		c.Type("svg", "utf-8")
		return c.Send(assets.MustAsset("webui/favicon.svg"))
	})
}

var DefaultListenAddress = ":3030"

// Serve runs the API server on the configured ApiListenAddress.
func (api *API) Listen() error {
	adrs := api.config.ListenAddress
	if adrs == "" {
		adrs = DefaultListenAddress
	}

	return api.app.Listen(adrs)
}

// HandleRoot is a handler for the root ("/") response.
//
// It serves the UI root page for (apparent) browsers, or a simple text
// message.
func (api *API) HandleRoot(c *fiber.Ctx) error {

	accept := c.Get("Accept")

	if strings.Contains(accept, "text/html") {
		return api.HandleUI(c)
	}

	return c.SendString(api.ident)
}

// HandleAgentList is a handler for listing the available agents.
func (api *API) HandleAgentsList(c *fiber.Ctx) error {

	list := []string{}
	for k := range api.sourceAgents {
		list = append(list, k)
	}
	sort.Strings(list)
	res := fiber.Map{"agents": list}

	return c.JSON(res)

}

// HandleAgentsNew is a handler for spawning a new agent for use.
func (api *API) HandleAgentsNew(c *fiber.Ctx) error {

	var payload RequestPayloadAgent
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON payload",
		})
	}

	// TODO: map api keys to available agents by name.
	// (make the api key thingy first of course)
	src_agent := api.sourceAgents[payload.Agent]
	if src_agent == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "requested agent not available",
		})
	}

	spawn, err := src_agent.Spawn()
	if err != nil {
		// TODO: log error here!
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to spawn agent",
		})
	}
	api.activeAgents[spawn.ULID.String()] = spawn
	res := fiber.Map{
		"id":          spawn.ULID,
		"name":        spawn.Name,
		"description": spawn.Description,
	}

	return c.JSON(res)

}

// HandleAgentsChat is a handler for spawning a new agent for use.
func (api *API) HandleAgentsChat(c *fiber.Ctx) error {

	id := c.Params("id")
	agent := api.activeAgents[id]
	if agent == nil {
		return fiber.ErrNotFound
	}

	var payload RequestPayloadChat
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON payload",
		})
	}
	if strings.TrimSpace(payload.Prompt) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "empty prompt",
		})
	}
	completion, err := agent.RunCompletionPrompt(payload.Prompt)
	if err != nil {
		// TODO: sniff out user vs agent vs llm errors
		// TODO: handle non-error errors appropriately e.g. finished...
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON payload",
		})
	}
	// TODO: want tool calls so we can see them!
	res := fiber.Map{
		"completion": completion,
	}
	return c.JSON(res)

}

// HandleUI is the handler for the simple chat UI.
func (api *API) HandleUI(c *fiber.Ctx) error {

	c.Type("html", "utf-8")

	if c.Method() == fiber.MethodGet {
		return c.Send(assets.MustAsset("webui/root.html"))
	}

	// Validate the key, and get the agents for that key.
	api_key := strings.TrimSpace(c.FormValue("api_key"))
	user_name := "Anonymous User"
	if api.config.NoKeys {
		api_key = "" // don't take a chance on weird keys breaking JS.
	} else {
		// TODO: look up, error with 404 err-badkey if not found.
		// assign name if we have it, or default to anon above.
	}
	agent_names := api.KeyAgentNames(api_key)

	// No agents for the key?  Nothing to do then.
	if len(agent_names) == 0 {
		return c.Status(fiber.StatusServiceUnavailable).Send(
			assets.MustAsset("webui/err-noagents.html"))
	}

	// Serve our SPA with the values we need in the DOM.
	// (Because templates are not worth it for something this simple).
	agent_opts := make([]string, len(agent_names))
	for i, n := range agent_names {
		agent_opts[i] = fmt.Sprintf(`<option value="%s">`, html.EscapeString(n))
	}
	page := assets.MustAssetString("webui/app.html")
	f := `<!-- user definition -->
<form id="user-session" class="hidden">
<input type="hidden" id="user-api-key" value="%s">
<input type="hidden" id="user-name" value="%s">
<select class="hidden" id="user-agent-name">
%s
</select>
</form>`
	page += fmt.Sprintf(f,
		html.EscapeString(api_key),
		html.EscapeString(user_name),
		strings.Join(agent_opts, "\n"))

	return c.SendString(page)
}

// KeyAgentNames returns the names of the API's agents available to the given
// api_key base on configured access.
//
// If NoKeys is configured, returns the names of all agents.
func (api *API) KeyAgentNames(api_key string) []string {
	names := []string{}
	// TODO: for real
	for _, a := range api.sourceAgents {
		names = append(names, a.Name)
	}
	return names

}
