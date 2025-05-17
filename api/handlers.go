// Package api defines the Greenhead HTTP API.
package api

import (
	"fmt"
	"html"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/biztos/greenhead/assets"
)

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

	spawn, err := src_agent.SpawnSilent()
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
	// TODO: come up iwtha
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
