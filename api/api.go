// Package api defines the Greenhead HTTP API.
package api

import (
	"fmt"
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
	}
	api.setRoutes()
	return api, nil

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

	api.app.Post("/v1/agents/:id/completion", func(c *fiber.Ctx) error {
		return api.HandleAgentsCompletion(c)
	})

	api.app.Post("/v1/agents/:id/end", func(c *fiber.Ctx) error {
		return api.HandleAgentsEnd(c)
	})

	api.app.Get("/v1/ui", func(c *fiber.Ctx) error {
		return api.HandleUI(c)
	})

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
