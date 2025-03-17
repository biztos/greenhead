// Package agent contains code for creating and operating individual Agents:
// LLM clients with tool-execution abilities.
package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
	"github.com/biztos/greenhead/utils"
)

// AgentType supported by the built-in ApiClients.
const (
	AgentTypeOpenAi = "openai"
)

// Config describes the configuration of an Agent, and is usually supplied in
// a file.  Note that CLI options may override some configs
type Config struct {
	Name     string         `json:"name"`               // name of the agent
	Type     string         `json:"type"`               // type, e.g. AgentTypeOpenAi
	Model    string         `json:"model,omitempty"`    // model for the LLM, if applicable
	Endpoint string         `json:"endpoint,omitempty"` // endpoint if not default
	Tools    []string       `json:"tools"`              // allowed tools
	Stream   bool           `json:"stream"`             // stream responses to STDOUT
	Color    string         `json:"color"`              // color for console output
	BgColor  string         `json:"bg_color"`           // background color for console output
	Context  []*ContextItem `json:"context,omitempty"`  // context window for client
}

// ContextItem is a high-level representation of a message to add to the
// context window.  Note that it does *not* at this point include ToolCall or
// ToolResult.
type ContextItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ToolCall is a high-level representation of a tool call from the LLM.
type ToolCall struct {
	Id   string
	Name string
	Args string
}

// ToolResult holds the result of a tool call.
type ToolResult struct {
	Id     string
	Output any
}

// RawCompletion represents a single round-trip completion request and
// response in its original types, which can be postprocessed with type
// assertions in the Agent.  It *must* be JSON serializable but it need not
// be deserializable from the Agent's point of view.
type RawCompletion struct {
	Request  any
	Response any
}

// CompletionRequest is a high-level representation of a message to the LLM
// from the "user."  If ToolResults are included, it is normal for Content to
// be empty.
type CompletionRequest struct {
	Content     string
	ToolResults []*ToolResult
}

// CompletionResponse is a high-level representation of a single-choice
// completion response.
//
// As of now it supports raw messages and tool calls.
// TODO: support files, images, whatever else we can create.
// (However this TODO is not super high priority -- we are mostly concerned
// with calling our own functions, not with non-text generation.)
type CompletionResponse struct {
	FinishReason   string // TODO: consider not including this...
	Content        string
	ToolCalls      []*ToolCall
	Usage          *Usage
	RawCompletions []*RawCompletion
}

// Usage is a high-level representation of token usage.  Note that the meaning
// of this, and its real-world cost, depends on the API being used. The main
// purpose here is to allow configured token limits.
//
// TODO: possibly support training!
// TODO: possibly support internal-tool calls (how are they reported?)
// TODO: support audio tokens?  possible?
// TODO: support reasoning tokens how exactly?
type Usage struct {
	Input       int
	CachedInput int
	Output      int
	Reasoning   int
	Total       int // nb: Total is just whatever was reported as total.
}

// ApiClient abstracts the API client itself, allowing the use of different
// clients for the various LLM APIs -- or in some cases, using the same client
// package differently.
//
// The ApiClient is responsible for maintaining its own LLM context over its
// lifetime; the Add* functions are for setting initial context.
type ApiClient interface {

	// SetLogger sets the logger that the ApiClient should use for all log
	// calls.
	SetLogger(*slog.Logger)

	// SetPreFunc sets the preprocessor function that can manipulate the
	// outgoing request in place before it is sent to the LLM.
	//
	// This can be used to limit the amount of context sent to the LLM, for
	// example.  Tool results are included in the request, so they can also
	// be addressed here.
	//
	// Implementations should include an example function if supported.
	SetPreFunc(func(any) error)

	// SetPostFunc sets the postprocessor function that can manipulate the
	// incoming response in place before it is sent to the Agent.
	//
	// Implementations should include an example function if supported.
	SetPostFunc(func(any) error)

	// AddContext adds a prompt or response to the LLM context.
	AddContext(*ContextItem)

	// RunCompletion runs a completion and returns
	RunCompletion(context.Context, *CompletionRequest) (*CompletionResponse, error)

	// Check validates the underlying client by making a (presumably) no-cost
	// round-trip to the configured API endpoint, e.g.
	//
	//     https://api.openai.com/v1/models
	Check(context.Context) error
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
	Id     ulid.ULID
	Config *Config
	Client ApiClient
	Tools  map[string]tools.Tooler

	logger *slog.Logger
}

// Ident for the Agent combines Id and the configured Type and optional Name.
func (a *Agent) Ident() string {
	s := fmt.Sprintf("%s:%s", a.Id, a.Config.Type)
	if a.Config.Name != "" {
		s += ":" + a.Config.Name
	}
	return s
}

// String provides a hopefully-useful stringification of the agent.
func (a *Agent) String() string {
	return fmt.Sprintf("<Agent %s>", a.Ident())
}

// SetLogger overrides the logger in the Agent and its ApiClient.
//
// Note that this does *not* add the agent=<ident> attribute that is used by
// default.  The caller should add that or its equivalent if desired.
func (a *Agent) SetLogger(logger *slog.Logger) {
	a.logger = logger
	a.Client.SetLogger(logger)
}

// NewAgent returns an agent initialized for use based on cfg.  If any of the
// configured Tools are not registered, an error is returned.
//
// The logger is a default slog JSON logger to Stderr with an "agent" attr
// defined as the agent's Ident value.
func NewAgent(cfg *Config) (*Agent, error) {
	toolmap := map[string]tools.Tooler{}
	for _, name := range cfg.Tools {
		tool := registry.Get(name)
		if tool == nil {
			return nil, fmt.Errorf("tool not registered: %s", name)
		}
		toolmap[name] = tool
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
	for _, c := range cfg.Context {
		client.AddContext(c)
	}
	agent := &Agent{
		Id:     ulid.Make(),
		Config: cfg,
		Client: client,
		Tools:  toolmap,
	}
	agent.SetLogger(slog.New(slog.NewJSONHandler(os.Stderr, nil)).With(
		"agent",
		agent.Ident(),
	))

	return agent, nil

}

// RunCompletion runs a completion request for the given prompt, returning its
// CompletionResult after handling any tool calls in the responses and sending
// them back for new completions.  The *final* completion in such a chain is
// returned, but its RawCompletions field includes all round-trips.
func (a *Agent) RunCompletion(ctx context.Context, prompt string) (*CompletionResponse, error) {

	raws := []*RawCompletion{}
	req := &CompletionRequest{Content: prompt}
	res, err := a.Client.RunCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error running completion: %w", err)
	}
	raws = append(raws, res.RawCompletions...)
	// yeah this isn't what I think.
	did := false
	for len(res.ToolCalls) > 0 {
		if did {
			a.logger.Warn("AGAIN")
			a.logger.Info("FFS", "toolcalls", utils.MustJsonString(res.ToolCalls))
			time.Sleep(10 * time.Second)
		}
		did = true
		// Run all the tool calls, keeping their responses.
		//
		// FOR NOW we just bail on bad calls, but IRL maybe we should just
		// give the LLM back an error?  Nah, CONFIG THIS.
		// TODO: figure out best approach for this, probably config AllowToolError
		results := make([]*ToolResult, len(res.ToolCalls))
		for idx, call := range res.ToolCalls {
			tool := a.Tools[call.Name]
			if tool == nil {
				return nil, fmt.Errorf("no such tool for agent: %s", call.Name)
			}
			// TODO: config whether this logs args, might leak private info
			// into the logs when we don't want to.  Or debug level?
			a.logger.Info("calling tool", "tool", call.Name, "args", call.Args)
			output, err := tool.Exec(context.Background(), call.Args)
			if err != nil {
				output = map[string]string{"error": err.Error()}
			}
			results[idx] = &ToolResult{
				Id:     call.Id,
				Output: output,
			}
		}
		// Get a new reponse from that.
		req := &CompletionRequest{ToolResults: results}
		res, err = a.Client.RunCompletion(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("error running tool-result completion: %w", err)
		}
		raws = append(raws, res.RawCompletions...)
		// TODO: limit loops on tools!
	}

	// Any tool calls have completed and we have a result plus a set of raw
	// completions that override the current one.
	return &CompletionResponse{
		Content:        res.Content,
		RawCompletions: raws,
	}, nil

}

func init() {
	RegisterNewClientFunc("openai", NewOpenAiClient)
}
