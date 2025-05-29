// api/handlers.go

package api

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/assets"
	"github.com/biztos/greenhead/utils"
)

type RequestPayloadAgent struct {
	Agent string `json:"agent"`
}

type RequestPayloadChat struct {
	Prompt string `json:"prompt"`
}

// HandleRoot is a handler for the root ("/") response.
//
// It serves the UI root page for (apparent) browsers, or a simple text
// message.
func (api *API) HandleRoot(c *fiber.Ctx) error {

	accept := c.Get("Accept")

	if !api.config.NoUI && strings.Contains(accept, "text/html") {
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

	if !api.config.NoKeys {
		key := c.Locals("access_key").(*Key)
		if !api.access.AgentAllowed(key, payload.Agent) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "agent not allowed",
			})
		}
	}

	// TODO: map api keys to available agents by name.
	// (make the api key thingy first of course)
	src_agent := api.sourceAgents[payload.Agent]
	if src_agent == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "agent not available",
		})
	}

	spawn, err := src_agent.SpawnSilent()
	if err != nil {
		// TODO: log error here!
		api.logger.Error("failed to spawn agent", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to spawn agent",
		})
	}
	api.activeAgents[spawn.ULID.String()] = spawn
	api.logger.Info("spawned new agent", "agent", spawn.Ident())
	res := fiber.Map{
		"id":          spawn.ULID,
		"name":        spawn.Name,
		"description": spawn.Description,
	}

	return c.JSON(res)

}

// ChatResponse represents a simple chat completion response.
type ChatResponse struct {
	Content   string            `json:"content"`
	ToolCalls []*agent.ToolCall `json:"tool_calls"`
}

// HandleAgentsChat is a handler for executing a simple chat request.
//
// On success, a ChatResponse will be sent.
func (api *API) HandleAgentsChat(c *fiber.Ctx) error {

	res, err := api.runAgentCompletion(c)
	if err != nil {
		return err
	}

	// TODO: limit access to tool_calls either by config or per-user
	chat_res := &ChatResponse{
		Content:   res.Content,
		ToolCalls: res.ToolCalls,
	}
	return c.JSON(chat_res)

}

// HandleAgentsCompletion is a handler for executing a completion request and
// returning the full response.
//
// On success, a full agent.CompletionResponse will be sent.
func (api *API) HandleAgentsCompletion(c *fiber.Ctx) error {

	res, err := api.runAgentCompletion(c)
	if err != nil {
		return err
	}
	return c.JSON(res)

}

// shared logic for chat handlers.
func (api *API) runAgentCompletion(c *fiber.Ctx) (*agent.CompletionResponse, error) {

	id := c.Params("agent_id")
	active_agent := api.activeAgents[id]
	if active_agent == nil {
		return nil, fiber.ErrNotFound
	}

	// Set up a context that *should* (depending on the underlying client
	// setup) cancel the LLM request when Fiber times out or detects that
	// the client has disconnected.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-c.Context().Done() // fasthttp.RequestCtx Done()
		cancel()
	}()

	var payload RequestPayloadChat
	if err := c.BodyParser(&payload); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON payload",
		})
	}
	if strings.TrimSpace(payload.Prompt) == "" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "empty prompt",
		})
	}

	req := &agent.CompletionRequest{Content: payload.Prompt}
	res, err := active_agent.RunCompletion(ctx, req)
	if err != nil {
		// TODO: sniff out user vs agent vs llm errors
		// TODO: handle non-error errors appropriately e.g. finished...
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return res, nil
}

type UiPageUser struct {
	ApiKey     string   `json:"api_key"`
	Name       string   `json:"name"`
	AgentNames []string `json:"agent_names"`
}
type UiPageConfig struct {
	User UiPageUser `json:"user"`
}

// HandleUI is the handler for the simple chat UI.
func (api *API) HandleUI(c *fiber.Ctx) error {

	c.Type("html", "utf-8")

	if c.Method() == fiber.MethodGet {
		return c.Send(assets.MustAsset("webui/root.html"))
	}

	// For no-key validation, use fake values.
	var auth_key string
	var key_name string
	var agent_names []string
	if api.config.NoKeys {
		api.logger.Warn("serving UI in no-keys mode")
		auth_key = "no-keys"
		key_name = "Anonymous User"
		agent_names = api.AgentNames(nil)
	} else {
		// Validate the key; don't send a full UI if it's bogus.
		auth_key = strings.TrimSpace(c.FormValue("auth_key"))
		key := api.access.GetKey(auth_key)
		if key == nil {
			api.logger.Debug("auth_key not found")
			// wtf then?!
			return fiber.ErrUnauthorized
		}
		key_name = key.Name

		// The key must be allowed to use the UI.
		if !api.access.EndpointAllowed(key, "/v1/ui") {
			api.logger.Debug("endpoint not allowed")
			return fiber.ErrUnauthorized
		}

		// Get the agent names. If no agents, you are also disallowed here.
		agent_names = api.AgentNames(key)
		if len(agent_names) == 0 {
			api.logger.Debug("no agents allowed")
			return fiber.ErrUnauthorized
		}
	}

	// Serve our SPA with the values we need stuck on the end.
	// (Because templates are not worth it for something this simple).
	page := assets.MustAssetString("webui/app.html")
	page_config := UiPageConfig{
		User: UiPageUser{
			ApiKey:     auth_key,
			Name:       key_name,
			AgentNames: agent_names,
		},
	}
	// GPT can't decide if this is an XSS risk or not, seems to me NOT.
	injection := fmt.Sprintf("<script>window.__CONFIG__ = %s</script>",
		utils.MustJsonString(page_config))
	html := strings.Replace(page, "</html>", injection+"</html>", 1)

	return c.SendString(html)
}

// HandleAgentsEnd is a handler for ending interaction with an agent.
//
// The agent is removed from operation and will be Not Found for futher
// requests.
//
// Note that future work may involve freezing/thawing agents, but at this time
// deletion is permanent on the server.
//
// TODO: prove that completion requests in flight are unaffected.
func (api *API) HandleAgentsEnd(c *fiber.Ctx) error {

	id := c.Params("agent_id")
	active_agent := api.activeAgents[id]
	if active_agent == nil {
		return fiber.ErrNotFound
	}

	delete(api.activeAgents, id)

	return c.JSON(fiber.Map{"success": true})

}
