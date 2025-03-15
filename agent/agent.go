// Package agent contains code for creating and operating individual Agents:
// LLM clients with tool-execution abilities.
package agent

import (
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/biztos/greenhead/registry"
)

// AgentType
const (
	AgentTypeOpenAi = "openai"
)

// Config describes the configuration of an Agent.
type Config struct {
	Name     string   `json:"name"`               // name of the agent
	Type     string   `json:"type"`               // type, e.g. AgentTypeOpenAi
	Endpoint string   `json:"endpoint,omitempty"` // endpoint if not default
	Tools    []string `json:"tools"`              // allowed tools
}

// ToolCall is an abstract representation of a tool call from the LLM.
type ToolCall struct {
	Id   string
	Name string
	Args string
}

// ToolResult holds the result of a tool call.
type ToolResult struct {
	Id     string
	Result any
	Err    error
}

// ApiClient abstracts the API client itself.
type ApiClient interface {
	// SetLogName sets the logging identifier, which should match the Agent's.
	SetLogName(name string)

	// AddSystemPrompt adds a system prompt to the LLM context.
	AddSystemPrompt(content string)

	// AddUserPrompt adds a user prompt to the LLM context.
	AddUserPrompt(content string)

	// AddAssistantResponse adds an assistant response to the LLM context.
	AddAssistantResponse(content string)

	// TODO: figure out how exactly to handle tool calls.
	//
	// IIRC we can't actually add them to context, without valid IDs the LLM
	// freaks out (at least OpenAI, at least IIRC).
	//
	// Anyway the ApiClient is responsible for maintaining context, right? So
	// we should not have to explicitly add the tool calls, assistant resp,
	// and so on.  That should be done internally.
	//
	// AddToolResults adds a set of tool call-results to the LLM context.
	AddToolResults([]*ToolResult)

	RunCompletion() error

	// Check validates the underlying client by making a (presumably) no-cost
	// round-trip to the configured API endpoint, e.g.
	//
	//     https://api.openai.com/v1/models
	Check() error

	// TBD... make the concrete type first then see how it goes.
	// GetToolCalls() []
}

var newClientFunc = map[string]func(*Config) (ApiClient, error){}

// RegisterNewClientFunc registers an agent type with a function returning a
// client for that type.  It is normally called in the init() function of a
// package.  Later registrations of the same name take precedence.
func RegisterNewClientFunc(agent_type string, f func(*Config) (ApiClient, error)) {
	newClientFunc[agent_type] = f
}

// Agent is a single "agentic" (tool-executing) LLM client.
type Agent struct {
	Id      ulid.ULID
	Config  *Config
	Client  ApiClient
	LogName string
}

// NewAgent returns an agent initialized for use based on cfg.  If any of the
// configured Tools are not registered, an error is returned.
func NewAgent(cfg *Config) (*Agent, error) {
	for _, tool := range cfg.Tools {
		if registry.Get(tool) == nil {
			return nil, fmt.Errorf("tool not registered: %s", tool)
		}
	}
	cfunc := newClientFunc[cfg.Type]
	if cfunc == nil {
		return nil, fmt.Errorf("no client for type %q", cfg.Type)
	}
	client, err := cfunc(cfg)
	if err != nil {
		return nil, fmt.Errorf("error initializing client for type %q: %w",
			cfg.Type, err)
	}
	id := ulid.Make()
	log_name := fmt.Sprintf("agent-%s", id)
	if cfg.Name != "" {
		log_name += ":" + cfg.Name
	}
	client.SetLogName(log_name)

	return &Agent{
		Id:     id,
		Config: cfg,
		Client: client,
	}, nil
}

func init() {
	RegisterNewClientFunc("openai", NewOpenAiClient)
}
