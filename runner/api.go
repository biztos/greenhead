package runner

import (
	"fmt"
	"io"
	"sort"

	"github.com/gofiber/fiber/v2"

	"github.com/biztos/greenhead/agent"
)

// CheckAPI instantiates an API and writes a log to its logger, printing OK
// to w if successful.
func (r *Runner) CheckAPI(w io.Writer) error {

	api, err := NewAPI(r)
	if err != nil {
		return err
	}
	// TODO: logs!
	fmt.Fprintln(w, "agents", len(api.agents))
	fmt.Fprintln(w, "OK")
	return nil

}

// ServeAPI instantiates an API and calls Listen.
func (r *Runner) ServeAPI(w io.Writer) error {

	api, err := NewAPI(r)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, "Staring server on", api.address)
	return api.Listen()

}

// API represents an instance of the Greenhead HTTP API.
type API struct {
	runner  *Runner
	app     *fiber.App
	agents  map[string]*agent.Agent
	address string
}

// NewAPI creates and API based on r.
func NewAPI(r *Runner) (*API, error) {

	// TODO: consider having no agents, only what you define on the API.
	// Seems like a bad idea to not have at least one available agent though.
	if len(r.Agents) == 0 {
		return nil, fmt.Errorf("at least one agent must be defined")
	}
	address := r.Config.ApiListenAddress
	if address == "" {
		address = ":3000" // fiber default.
	}

	// TODO: get some of the API-specific stuff from the config, w/API section
	// https://pkg.go.dev/github.com/gofiber/fiber/v2#Config
	fiber_cfg := fiber.Config{
		AppName: "GREENHEAD HTTP API",
	}
	app := fiber.New(fiber_cfg)
	agents := map[string]*agent.Agent{}
	for _, a := range r.Agents {
		if agents[a.Name] != nil {
			return nil, fmt.Errorf("duplicate agent by name: %q", a.Name)
		}
		clone, err := a.Clone() // don't accidentally mess w/agents outside
		if err != nil {
			return nil, fmt.Errorf("agent clone error for %q: %w", a.Name, err)
		}
		agents[a.Name] = clone
	}
	api := &API{
		runner:  r,
		app:     app,
		agents:  agents,
		address: address,
	}
	api.setRoutes()
	return api, nil

}

// set up the routing for the fiber app, with access to the agents et al.
func (api *API) setRoutes() {

	api.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("去做什么?") // What to do? (or similar it seems)
	})

	api.app.Get("/v1/agents/list", func(c *fiber.Ctx) error {
		list := []string{}
		for k := range api.agents {
			list = append(list, k)
		}
		sort.Strings(list)
		res := fiber.Map{"agents": list}

		return c.JSON(res)
	})

	// Posting a chat will be harder!

}

// Serve runs the API server on the configured ApiListenAddress.
func (api *API) Listen() error {
	return api.app.Listen(api.address)
}
