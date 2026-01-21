package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/biztos/greenhead/ghd/assets"
)

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

	api.app.Post("/v1/agents/:agent_id/chat", func(c *fiber.Ctx) error {
		return api.HandleAgentsChat(c)
	})

	api.app.Post("/v1/agents/:agent_id/completion", func(c *fiber.Ctx) error {
		return api.HandleAgentsCompletion(c)
	})

	api.app.Post("/v1/agents/:agent_id/end", func(c *fiber.Ctx) error {
		return api.HandleAgentsEnd(c)
	})

	if !api.config.NoUI {
		api.app.Get("/v1/ui", func(c *fiber.Ctx) error {
			return api.HandleUI(c)
		})

		api.app.Post("/v1/ui", func(c *fiber.Ctx) error {
			return api.HandleUI(c)
		})

		// Serve a favicon because the requests are annoying.
		api.app.Get("/favicon.png", func(c *fiber.Ctx) error {
			c.Type("png")
			return c.Send(assets.MustAsset("webui/favicon.png"))
		})

		// Also serve it for .ico requests because effing browsers.
		api.app.Get("/favicon.ico", func(c *fiber.Ctx) error {
			c.Type("png")
			return c.Send(assets.MustAsset("webui/favicon.png"))
		})
	}

}
