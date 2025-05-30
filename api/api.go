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

	// You really *should* have agents, but you don't *have* to have them.
	// At some future point we will support loading and/or creating agents
	// after the API is running.

	// Set up access, unless we don't.
	var encoder func(string) string
	if cfg.RawKeys {
		encoder = NotEncodeAuthKey
	} else {
		encoder = EncodeAuthKey
	}
	var access *Access
	var err error
	var default_auth_key string
	if !cfg.NoKeys {
		// If there is nothing, use the default.
		roles := cfg.Roles
		keys := cfg.Keys
		if len(roles) == 0 && len(keys) == 0 && cfg.AccessFile == "" {
			roles = DefaultRoles
			keys = DefaultKeys
			default_auth_key = keys[0].AuthKey
		}

		access, err = NewAccess(roles, keys, cfg.AccessFile, encoder)
		if err != nil {
			return nil, fmt.Errorf("access setup error: %w", err)
		}

	}

	ident := fmt.Sprintf("GREENHEAD %s HTTP API %s",
		version.Version, version.ApiVersion)
	fiber_cfg := cfg.FiberConfig(ident)
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
		defaultKey:   default_auth_key,
	}
	// Set up app routes and middleware. NB: ORDER MATTERS.
	if cfg.LogFiber {
		app.Use(logger.New())
	} else {
		// TODO: useful filter to exclude the instrumentation, which is
		// annoying AF when doing development and testing.
		app.Use(slogfiber.New(slog.Default()))
	}
	app.Use(api.KeyAccess())
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
		// Warn that we're running with the default key, but wait a bit to
		// get the log after the startup message.  My OCD impulses, cheeeze...
		time.AfterFunc(time.Second, func() {
			api.logger.Warn("ALL-ACCESS DEFAULT KEY IN USE",
				"auth_key", api.access.keyEncoder(api.defaultKey))
		})
	}
	if len(api.sourceAgents) == 0 {
		// Likewise warn if the API isn't much usable due to lack of "agency."
		time.AfterFunc(time.Second, func() {
			api.logger.Warn("NO AGENTS DEFINED")
		})
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
