// Package api defines the Greenhead HTTP API.
package api

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/version"
)

// Config defines the configuration of the API.
//
// It is included as part of the Runner config.
//
// Fields beginning with "App" are passed to the Fiber app.
type Config struct {

	// General configuration:
	ListenAddress string `toml:"listen_address"` // Address for serving, e.g. ":3000"
	LogFiber      bool   `toml:"log_fiber"`      // Use default Fiber logger for requests.

	// Access control:
	Roles      []*Role `toml:"roles"`       // Roles defining access.
	Keys       []*Key  `toml:"keys"`        // Keys mapping to roles by name.
	AccessFile string  `toml:"access_file"` // TOML file for (more) Roles and Keys.
	RawKeys    bool    `toml:"raw_keys"`    // Use raw, unencoded API keys.
	NoKeys     bool    `toml:"no_keys"`     // DO NOT require API keys.
	NoUI       bool    `toml:"no_ui"`       // DO NOT expose the web UI.

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
	logger       *slog.Logger
	sourceAgents map[string]*agent.Agent
	activeAgents map[string]*agent.Agent
	access       *Access
	defaultKey   string
}

// NewAPI creates an API instance.
func NewAPI(cfg *Config, agents []*agent.Agent) (*API, error) {

	// TODO: consider having no agents, only what you define on the API.
	// Seems like a bad idea to not have at least one available agent though.
	if len(agents) == 0 {
		return nil, fmt.Errorf("at least one agent must be defined")
	}

	// Set up access, unless we don't.
	var encoder func(string) string
	if cfg.RawKeys {
		encoder = NotEncodeAuthKey
	} else {
		encoder = EncodeAuthKey
	}
	var access *Access
	var err error
	var default_key string
	if !cfg.NoKeys {
		// If there is nothing, use the default.
		if len(cfg.Roles) == 0 && len(cfg.Keys) == 0 {
			access, default_key = DefaultAccess(encoder)
		} else {
			access, err = NewAccess(cfg.Roles, cfg.Keys, encoder)
			if err != nil {
				return nil, fmt.Errorf("access setup error: %w", err)
			}
		}
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
	api := &API{
		ident:        ident,
		config:       cfg,
		app:          app,
		logger:       slog.Default(),
		sourceAgents: sourceAgents,
		activeAgents: map[string]*agent.Agent{},
		access:       access,
		defaultKey:   default_key,
	}
	// Set up app routes and middleware. NB: ORDER MATTERS.
	app.Use(api.KeyAccess())
	if cfg.LogFiber {
		app.Use(logger.New())
	} else {
		// TODO: useful filter to exclude the instrumentation, which is
		// annoying AF when doing development and testing.
		app.Use(slogfiber.New(slog.Default()))
	}
	app.Use(recover.New()) // TODO: why exactly?
	api.setRoutes()

	return api, nil

}

var DefaultListenAddress = ":3030"

// Serve runs the API server on the configured ApiListenAddress.
func (api *API) Listen() error {
	adrs := api.config.ListenAddress
	if adrs == "" {
		adrs = DefaultListenAddress
	}
	if api.defaultKey != "" {
		fmt.Println("**")
		fmt.Println("** ALL-ACCESS DEFAULT API KEY:", api.access.keyEncoder(api.defaultKey))
		fmt.Println("**")
	}

	return api.app.Listen(adrs)
}

// GetKey calls GetKey on the underlying Access of the API.
func (api *API) GetKey(auth_key string) *Key {
	return api.access.GetKey(auth_key)
}

// AgentNames returns the names of the API's agents available to the given
// key based on access.
//
// If NoKeys is configured, returns the names of all agents.
func (api *API) AgentNames(key *Key) []string {
	names := []string{}
	for _, a := range api.sourceAgents {
		if api.config.NoKeys || api.access.AgentAllowed(key, a.Name) {
			names = append(names, a.Name)
		}
	}
	return names

}
