// api/middleware.go

package api

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
)

func (api *API) KeyAccess() fiber.Handler {

	return func(c *fiber.Ctx) error {

		// Root and UI root don't use keys.
		if c.Path() == "/" || c.Path() == "/v1/ui" {
			return c.Next()
		}

		// Special case for annoying favicon.
		if c.Path() == "/favicon.png" || c.Path() == "/favicon.ico" {
			return c.Next()
		}

		// If we are not using keys, just ignore the header and stash a dummy
		// user.
		if api.config.NoKeys {
			slogfiber.AddCustomAttributes(c, slog.String("access_key", "no-keys"))
			return c.Next()
		}

		// First check that we have a valid key.
		hdr := c.Get("Authorization")
		if hdr == "" || !strings.HasPrefix(hdr, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}
		key := api.access.GetKey(strings.TrimPrefix(hdr, "Bearer "))
		if key == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unknown auth key",
			})
		}
		// Now "authz" it for the endpoint.
		if !api.access.EndpointAllowed(key, c.Path()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Endpoint not allowed",
			})
		}

		// If we have an agent ID, authz for that too.
		//
		// Why not bind the agent to the key?  First, it seems overkill, since
		// you could also just "steal" the AuthKey if security is the concern.
		// But you might also have some workflow in which one key creates
		// agents for another key to use, that is less exotic than it sounds
		// at first: imagine you have a set of workers and they are allowed
		// one agent each, but you have a master worker assigning them.
		//
		// In any case, here we check that the type (name) of agent is allowed
		// and thus also handle the (even farther-fetched?) case of revoking
		// access to a type of agent for a key.  Once reloading keys works.
		id := c.Params("agent_id")
		a := api.activeAgents[id]
		if a != nil && !api.access.AgentAllowed(key, a.Name) {
			// NB: nil-agent will get caught downstream
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Agent not allowed",
			})
		}

		// All good!
		slogfiber.AddCustomAttributes(c, slog.String("access", key.Name))
		c.Locals("access_key", key)
		return c.Next()

	}

}
