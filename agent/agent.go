// Package agent contains code for creating and operating individual Agents:
// LLM clients with tool-execution abilities.
package agent

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/oklog/ulid/v2"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/utils"
)

var DefaultPrintFunc = func(a ...any) { fmt.Print(a...) }

var ErrMaxCompletions = fmt.Errorf("max completions reached")

// Config describes the configuration of an Agent, and is usually supplied in
// a file.
//
// Note that in normal operation, runner configs will take precedence over
// agent configs.
type Config struct {
	Name        string   `toml:"name"`        // Name of the agent.
	Description string   `toml:"description"` // Description of the agent.
	Type        string   `toml:"type"`        // Type, e.g. AgentTypeOpenAi.
	Model       string   `toml:"model"`       // Model for the LLM, if applicable.
	Endpoint    string   `toml:"endpoint"`    // Endpoint if not default.
	Tools       []string `toml:"tools"`       // Allowed tools by name.

	Context []ContextItem `toml:"context"` // Context window for client.

	// Safety and limits:  (Zero generally means "no limit.")
	MaxCompletionTokens int  `toml:"max_completion_tokens"` // Max completion tokens *per completion* (may truncate responses).
	MaxCompletions      int  `toml:"max_completions"`       // Max number of completions to run.
	MaxTokens           int  `toml:"max_tokens"`            // Max number of total tokens for all operations.
	MaxToolChain        int  `toml:"max_toolchain"`         // Max number of tool call responses allowed in a row.
	AbortOnRefusal      bool `toml:"abort_on_refusal"`      // Abort if a completion is refused by an LLM.

	// Output control:
	Color     string `toml:"color"`           // Color for console output.
	BgColor   string `toml:"bg_color"`        // Background color for console output.
	Stream    bool   `toml:"stream"`          // Stream responses; if streaming not possible, print them.
	ShowCalls bool   `toml:"show_tool_calls"` // Show tool calls in output (experimental; can leak data).
	Silent    bool   `toml:"silent"`          // Suppress responses unless streamed.

	// Logging and debugging:
	Debug       bool   `toml:"debug"`         // Log at DEBUG level instead of INFO.
	LogFile     string `toml:"log_file"`      // Log to a file.
	NoLog       bool   `toml:"no_log"`        // Do not log at all.
	DumpDir     string `toml:"dump_dir"`      // Dump completions to this dir (can leak data).
	LogToolArgs bool   `toml:"log_tool_args"` // Log tool args (can leak data).

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
	Id   string `json:"id"`
	Name string `json:"name"`
	Args string `json:"args"`
}

// ToolResult holds the result of a tool call.
type ToolResult struct {
	Id     string `json:"id"`
	Output any    `json:"output"`
}

// RawCompletion represents a single round-trip completion request and
// response in its original types, which can be postprocessed with type
// assertions in the Agent.  It *must* be JSON serializable but it need not
// be deserializable from the Agent's point of view.
type RawCompletion struct {
	Request  any `json:"request"`
	Response any `json:"response"`
}

// CompletionRequest is a high-level representation of a message to the LLM
// from the "user."  If ToolResults are included, it is normal for Content to
// be empty.
type CompletionRequest struct {
	Content     string        `json:"content"`
	ToolResults []*ToolResult `json:"tool_results"`
}

// CompletionResponse is a high-level representation of a single-choice
// completion response.
//
// As of now it supports raw messages and tool calls.
// TODO: support files, images, whatever else we can create.
// (However this TODO is not super high priority -- we are mostly concerned
// with calling our own functions, not with non-text generation.)
type CompletionResponse struct {
	FinishReason   string           `json:"finish_reason"` // TODO: consider not including this...
	Content        string           `json:"content"`
	ToolCalls      []*ToolCall      `json:"tool_calls"`
	Usage          *Usage           `json:"usage"`
	RawCompletions []*RawCompletion `json:"raw_completions"`
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
	Input       int `json:"input"`
	CachedInput int `json:"cached_input"`
	Output      int `json:"output"`
	Reasoning   int `json:"reasoning"`
	Total       int `json:"total"` // nb: Total is just whatever was reported as total.
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
	ULID        ulid.ULID `json:"ulid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Model       string    `json:"model"`

	client    ApiClient
	toolnames []string

	completed int
	config    *Config
	printFunc func(a ...any)
	logger    *slog.Logger
	dumpdir   string
}

// Id returns the Agent's ULID as a string.
func (a *Agent) Id() string {
	return a.ULID.String()
}

// Ident returns an informative, uniquely identifying string.
func (a *Agent) Ident() string {
	s := fmt.Sprintf("%s:%s", a.ULID, a.Type)
	if a.Name != "" {
		s = fmt.Sprintf("%s:%s", s, a.Name)
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

// AddContextItem calls the ApiClient's AddContextItem.
func (a *Agent) AddContextItem(item ContextItem) {
	a.client.AddContextItem(item)
}

// NewAgent returns an agent initialized for use based on cfg.  If any of the
// configured Tools are not registered, an error is returned.
//
// TODO: consider the possibility of runtime tool registrations, in which case
// what do we do to keep the agent up to date?
func NewAgent(cfg *Config) (*Agent, error) {

	// Start with basics:
	a := &Agent{
		ULID:        ulid.Make(),
		Name:        cfg.Name,
		Description: cfg.Description,
		Type:        cfg.Type,
		Model:       cfg.Model,
		config:      cfg.Copy(),
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
	a.SetClient(client)
	client.SetTools(cfg.Tools)
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

	// Set up tools, checking for validity.
	if err := a.SetTools(cfg.Tools); err != nil {
		return nil, err
	}

	// Set up the streaming and color printing:
	pfunc, err := ColorPrintFunc(cfg.Color, cfg.BgColor)
	if err != nil {
		return nil, fmt.Errorf("error with print colors: %w", err)
	}
	a.printFunc = pfunc
	client.SetPrintFunc(pfunc)

	// Set up the logger.
	if cfg.NoLog {
		nologger := slog.New(slog.NewTextHandler(io.Discard, nil))
		a.SetLogger(nologger)
	} else {
		if err := a.InitLogger(cfg.LogFile, cfg.Debug); err != nil {
			return nil, fmt.Errorf("error initializing logger: %w", err)
		}
	}

	// Make sure the DumpDir exists if set, and is writable -- writing the
	// config there should do the trick!  Note that to keep things sane we
	// will make a subdir for the actual dump directory.
	if cfg.DumpDir != "" {
		a.dumpdir = filepath.Join(cfg.DumpDir, a.ULID.String())
		if err := os.MkdirAll(a.dumpdir, 0755); err != nil {
			return nil, fmt.Errorf("error creating dump directory: %w", err)
		}
		fn := fmt.Sprintf("%s-config.toml", a.ULID)
		cfg_file := filepath.Join(a.dumpdir, fn)
		if err := utils.MarshalTomlFile(cfg, cfg_file); err != nil {
			return nil, fmt.Errorf("error dumping config: %w", err)
		}
		client.SetDumpDir(a.dumpdir)
	}

	return a, nil

}

// SetClient sets the internal ApiClient to c, overriding anything set on
// initialization.
//
// This allows the use of arbitrary ApiClients that are not registered in
// this package.
func (a *Agent) SetClient(c ApiClient) {
	a.client = c
}

// SetTools sets the interal tools list for the agent and its ApiClient,
// handling regexp selection and checking for validity.
//
// Any input string deliminated with slashes, e.g. `/foo/`, is treated as a
// regular expression, and all registered names that match are included.
func (a *Agent) SetTools(names []string) error {
	valid_names, err := registry.MatchingNames(names)
	if err != nil {
		return err
	}
	a.toolnames = valid_names
	a.client.SetTools(valid_names)
	return nil
}

// Tools returns the list of tools available to the agent.
func (a *Agent) Tools() []string {
	return slices.Clone(a.toolnames)
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

// RunCompletionPrompt calls RunCompletion with background context and the
// provided prompt, returning the content of the response.
func (a *Agent) RunCompletionPrompt(prompt string) (string, error) {

	req := &CompletionRequest{Content: prompt}
	res, err := a.RunCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.Content, nil
}

// RunCompletion runs a completion request for the given prompt, returning its
// CompletionResult after handling any tool calls in the responses and sending
// them back for new completions.  The *final* completion in such a chain is
// returned, but its RawCompletions field includes all round-trips.
func (a *Agent) RunCompletion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {

	if a.config.MaxCompletions > 0 && a.completed >= a.config.MaxCompletions {
		return nil, fmt.Errorf("%w: %d", ErrMaxCompletions, a.completed)
	}

	raws := []*RawCompletion{}
	res, err := a.client.RunCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error running completion: %w", err)
	}
	raws = append(raws, res.RawCompletions...)
	tool_call_responses := 0
	for len(res.ToolCalls) > 0 {

		// We can in theory get multiple tool calls in succession, in which
		// case we watch for the tool chain.
		tool_call_responses++
		if a.config.MaxToolChain > 0 && tool_call_responses > a.config.MaxToolChain {
			// TODO: should this really be an error?  How best to handle these
			// abort cases?  Perhaps we should send a refusal to run tools?
			// A quota-exceeded tool output error?
			return nil, fmt.Errorf("max tool chain exceeded")
		}

		// Run all the tool calls, keeping their responses.
		//
		// FOR NOW we just bail on bad calls, but IRL maybe we should just
		// give the LLM back an error?  Nah, CONFIG THIS.
		// TODO: figure out best approach for this, probably config AllowToolError
		// TODO: concurrency if the tools allow it.
		results := make([]*ToolResult, len(res.ToolCalls))
		for idx, call := range res.ToolCalls {

			// Print tools as they arrive, if requested.
			// (Printing at the end will be confusing if tools take longer to run.)
			if !a.config.Stream && !a.config.Silent && a.config.ShowCalls {
				line := fmt.Sprintf("* tool_call: %s %s\n", call.Name, call.Args)
				a.printFunc(line)
			}

			// NB: we have no actual guarantee that the registered tools have
			// not changed since the last call; nor that the LLM is not trying
			// to call a disallowed tool. Thus we need to check that the tool
			// is both allowed, and currently registered.
			if !slices.Contains(a.toolnames, call.Name) {
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
	// TODO: print output if not stream and not silent

	// Print output, if desired.
	if !a.config.Stream && !a.config.Silent {
		a.printFunc(res.Content)
		a.printFunc("\n")

	}

	final_res := &CompletionResponse{
		Content:        res.Content,
		RawCompletions: raws,
	}

	// Dump the full round-trip if desired.
	a.completed++
	if a.config.DumpDir != "" {
		a.DumpCompletion(req, final_res)
	}

	// Any tool calls have completed and we have a result plus a set of raw
	// completions that override the current one.
	return final_res, nil

}

// DumpCompletion writes a JSON file for the full completion into the
// configured DumpDir, or the local directory if not set.
func (a *Agent) DumpCompletion(req *CompletionRequest, res *CompletionResponse) error {

	name := fmt.Sprintf("%s-%04d.json", a.ULID, a.completed)
	path := filepath.Join(a.dumpdir, name)
	v := map[string]any{"request": req, "response": res}
	if err := utils.MarshalJsonFile(v, path); err != nil {
		return fmt.Errorf("error dumping completion: %w", err)
	}
	return nil

}
