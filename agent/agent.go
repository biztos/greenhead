// Package agent contains code for creating and operating individual Agents:
// LLM clients with tool-execution abilities.
package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/oklog/ulid/v2"

	"github.com/biztos/greenhead/registry"
)

// AgentType supported by the built-in ApiClients.
const (
	AgentTypeOpenAi = "openai"
)

var DefaultPrintFunc = func(a ...any) { fmt.Print(a...) }

// Config describes the configuration of an Agent, and is usually supplied in
// a file.
//
// Note that in normal operation, runner configs will take precedence over
// agent configs.
type Config struct {
	Name     string   `toml:"name"`               // Name of the agent.
	Type     string   `toml:"type"`               // Type, e.g. AgentTypeOpenAi.
	Model    string   `toml:"model,omitempty"`    // Model for the LLM, if applicable.
	Endpoint string   `toml:"endpoint,omitempty"` // Endpoint if not default.
	Tools    []string `toml:"tools"`              // Allowed tools by name.

	Context []ContextItem `toml:"context,omitempty"` // Context window for client.

	// Safety and limits:  (Zero generally means "no limit.")
	MaxCompletionTokens int  `toml:"max_completion_tokens"`      // Max completion tokens *per completion* (may truncate responses).
	MaxCompletions      int  `toml:"max_completions,omitempty"`  // Max number of completions to run.
	MaxTokens           int  `toml:"max_tokens,omitempty"`       // Max number of total tokens for all operations.
	MaxToolChain        int  `toml:"max_tool_chain,omitempty"`   // Max number of tool call responses allowed in a row.
	AbortOnRefusal      bool `toml:"abort_on_refusal,omitempty"` // Abort if a completion is refused by an LLM.

	// Output control:
	Color     string `toml:"color"`             // Color for console output.
	BgColor   string `toml:"bg_color"`          // Background color for console output.
	Stream    bool   `toml:"stream"`            // Stream responses; if streaming not possible, print them.
	ShowCalls bool   `toml:"stream_tool_calls"` // Show tool calls in output (experimental; can leak data).
	Silent    bool   `toml:"silent"`            // Suppress responses unless streamed.

	// Logging and debugging:
	Debug       bool   `toml:"debug"`                   // Log at DEBUG level instead of INFO.
	LogFile     string `toml:"log_file,omitempty"`      // Log to a file.
	DumpDir     string `toml:"dump_dir,omitempty"`      // Dump completions to this dir (can leak data).
	LogToolArgs bool   `toml:"log_tool_args,omitempty"` // Log tool args (can leak data).

}

// Copy returns a deep copy of c.
func (c *Config) Copy() *Config {
	n := *c
	n.Tools = make([]string, len(c.Tools))
	copy(n.Tools, c.Tools)
	n.Context = make([]ContextItem, len(c.Context))
	copy(n.Context, c.Context)
	return &n
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
// lifetime; the Add* functions are for setting initial context.  However, the
// agent *may* do this itself with the ClearContext function.
type ApiClient interface {

	// SetLogger sets the logger that the ApiClient should use for all log
	// calls.
	SetLogger(*slog.Logger)

	// SetStreaming sets whether responses should be streamed to Stdout as
	// they are received.  If streaming is not supported, responses should be
	// printed when they are received.  In both cases, the print function from
	// SetPrintFunc should be used for printing output.
	SetStreaming(bool)

	// SetShowCalls sets whether tool calls should be streamed to Stdout
	// the same as content responses.  This is experimental and could leak
	// data into the session that you would rather keep private.
	SetShowCalls(bool)

	// SetPrintFunc sets the function used to print output.
	SetPrintFunc(func(v ...any))

	// SetPreFunc sets a function that processes the raw request before it
	// is sent to the LLM.
	//
	// This allows customization of ApiClients without additional types, e.g.
	// for content filtering.
	//
	// It is up to the implementation to call the pre- and post-functions in
	// RunCompletion. If this is not supported, an error should be returned.
	SetPreFunc(func(ApiClient, any) error) error

	// SetPostFunc sets a function that processes the raw response when it is
	// received from the LLM.
	//
	// It is up to the implementation to call the pre- and post-functions in
	// RunCompletion. If this is not supported, an error should be returned.
	SetPostFunc(func(ApiClient, any) error) error

	// SetModel sets the model used by the API.  If the model is not supported
	// it should return an error.
	SetModel(string) error

	// SetTools sets the tools that will be described to the LLM as callable.
	//
	// These must be available in the registry when SetTools is called.
	SetTools([]string) error

	// SetMaxCompletionTokens sets the maximum number of tokens the LLM should
	// produce on the *next* completion.
	SetMaxCompletionTokens(int)

	// ClearContext clears any existing LLM context.
	ClearContext()

	// AddContextItem adds a prompt or response to the LLM context.
	AddContextItem(ContextItem)

	// RunCompletion runs a completion and returns
	RunCompletion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	// Check validates the underlying client by making a (presumably) no-cost
	// round-trip to the configured API endpoint, e.g.
	//
	//     https://api.openai.com/v1/models
	Check(context.Context) error
}

var newApiClientFunc = map[string]func() (ApiClient, error){}

// RegisterNewApiClientFunc registers an agent type with a function returning
// an ApiClient for that type.  It is normally called in the init() function
// of the package defining that type.
//
// Later registrations of the same name take precedence.
func RegisterNewApiClientFunc(agent_type string, f func() (ApiClient, error)) {
	newApiClientFunc[agent_type] = f
}

// Agent is a single "agentic" (tool-executing) LLM client.
type Agent struct {
	ULID      ulid.ULID
	client    ApiClient
	config    *Config
	printFunc func(a ...any)
	logger    *slog.Logger
}

// Id returns the Agent's ULID as a string.
func (a *Agent) Id() string {
	return a.ULID.String()
}

// Name returns the Agent's configured name, which is optional..
//
// Note that there is no guarantee of uniqueness for the name.
func (a *Agent) Name() string {
	return a.config.Name
}

// Ident returns an informative, uniquely identifying string.
func (a *Agent) Ident() string {
	s := fmt.Sprintf("%s:%s", a.ULID, a.config.Type)
	if a.config.Name != "" {
		s = fmt.Sprintf("%s:%s", s, a.config.Name)
	}
	return s
}

// String provides a standard stringification of the agent.
func (a *Agent) String() string {
	return fmt.Sprintf("<Agent %s>", a.Ident())
}

// Logger returns the logger set with SetLogger.
//
// This is useful for logging things "as" the agent, i.e. with its ident
// component.
func (a *Agent) Logger() *slog.Logger {
	return a.logger
}

// SetLogger overrides the logger in the Agent and its ApiClient.
//
// Note that this does *not* add the agent=<ident> attribute that is used by
// default.  The caller should add that or its equivalent if desired, as does
// NewAgent.
func (a *Agent) SetLogger(logger *slog.Logger) {
	a.logger = logger
	a.client.SetLogger(logger)
}

// Check calls the ApiClient's Check function with ctx.
func (a *Agent) Check(ctx context.Context) error {
	return a.client.Check(ctx)
}

// NewAgent returns an agent initialized for use based on cfg.  If any of the
// configured Tools are not registered, an error is returned.
//
// TODO: consider the possibility of runtime tool registrations, in which case
// what do we do to keep the agent up to date?
func NewAgent(cfg *Config) (*Agent, error) {

	// Start with basics:
	a := &Agent{
		ULID:   ulid.Make(),
		config: cfg.Copy(),
	}

	// Get an ApiClient to set up:
	cfunc := newApiClientFunc[cfg.Type]
	if cfunc == nil {
		return nil, fmt.Errorf("no client for type %q", cfg.Type)
	}
	client, err := cfunc()
	if err != nil {
		return nil, fmt.Errorf("error initializing client for type %q: %w",
			cfg.Type, err)
	}
	a.client = client
	client.SetStreaming(cfg.Stream)
	client.SetModel(cfg.Model)
	client.SetShowCalls(cfg.ShowCalls)
	client.SetMaxCompletionTokens(cfg.MaxCompletionTokens)

	// Add any configured context to the ApiClient.  Note that we do *not*
	// clear the context here: if the newClientFunc wants to include premade
	// context, we leave that alone.
	for _, c := range cfg.Context {
		client.AddContextItem(c)
	}

	// Add tools to the ApiClient, checking for validity first:
	// TODO: support tool names like "foo*" but only here, client gets
	// foo_bar, foo_boo
	for _, name := range cfg.Tools {
		_, err := registry.Get(name)
		if err != nil {
			return nil, err
		}
	}
	client.SetTools(cfg.Tools)

	// Set up the streaming and color printing:
	pfunc, err := ColorPrintFunc(cfg.Color, cfg.BgColor)
	if err != nil {
		return nil, fmt.Errorf("error with print colors: %w", err)
	}
	a.printFunc = pfunc
	client.SetPrintFunc(pfunc)

	// Set up the logger.
	if err := a.InitLogger(cfg.LogFile, cfg.Debug); err != nil {
		return nil, fmt.Errorf("error initializing logger: %w", err)
	}

	return a, nil

}

// InitLogger sets up a slog.Logger to log to the file (or Stderr) at Info or
// Debug level, then calls SetLogger with it.
func (a *Agent) InitLogger(file string, debug bool) error {

	var handler *slog.JSONHandler
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	if file == "" {
		// Log to standard error.
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// Log to file.
		fh, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}

		handler = slog.NewJSONHandler(fh, &slog.HandlerOptions{
			Level: level,
		})
	}

	a.SetLogger(slog.New(handler).With(
		"agent",
		a.Ident(),
	))
	return nil
}

// RunCompletion runs a completion request for the given prompt, returning its
// CompletionResult after handling any tool calls in the responses and sending
// them back for new completions.  The *final* completion in such a chain is
// returned, but its RawCompletions field includes all round-trips.
func (a *Agent) RunCompletion(ctx context.Context, prompt string) (*CompletionResponse, error) {

	raws := []*RawCompletion{}
	req := &CompletionRequest{Content: prompt}
	res, err := a.client.RunCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error running completion: %w", err)
	}
	raws = append(raws, res.RawCompletions...)
	for len(res.ToolCalls) > 0 {

		// Run all the tool calls, keeping their responses.
		//
		// FOR NOW we just bail on bad calls, but IRL maybe we should just
		// give the LLM back an error?  Nah, CONFIG THIS.
		// TODO: figure out best approach for this, probably config AllowToolError
		results := make([]*ToolResult, len(res.ToolCalls))
		for idx, call := range res.ToolCalls {
			// NB: we have no actual guarantee that the registered tools have
			// not changed since the last call; nor that the LLM is not trying
			// to call a disallowed tool. Thus we need to check that the tool
			// is both allowed, and currently registered.
			//
			// TODO (someday): support changing allowed tools at runtime.
			if !slices.Contains(a.config.Tools, call.Name) {
				return nil, fmt.Errorf("no such tool for agent: %s", call.Name)
			}
			tool, err := registry.Get(call.Name)
			if err != nil {
				return nil, err
			}
			// Only log the tool args if that's configured, which by default
			// it's not.
			if a.config.LogToolArgs {
				a.logger.Info("calling tool", "tool", call.Name, "args", call.Args)
			} else {
				a.logger.Info("calling tool", "tool", call.Name)
			}
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
		res, err = a.client.RunCompletion(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("error running tool-result completion: %w", err)
		}
		raws = append(raws, res.RawCompletions...)
		// TODO: limit loops on tools!
	}

	// TODO: limits
	// TODO: bail on refusal
	// TODO: save to DumpDir if desired
	// TODO: print output if not stream and not silent

	// Any tool calls have completed and we have a result plus a set of raw
	// completions that override the current one.
	return &CompletionResponse{
		Content:        res.Content,
		RawCompletions: raws,
	}, nil

}
