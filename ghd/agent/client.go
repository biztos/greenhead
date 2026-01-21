// agent/client.go

package agent

import (
	"context"
	"errors"
	"log/slog"
)

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

	// SetDumpDir sets a directory into which the ApiClient *may* write any
	// debug information such as raw requests or responses, to supplement the
	// data dumped by the agent itself.
	SetDumpDir(string)

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

var ErrPlaceholder = errors.New("Placeholder function.")

// BasicApiClient satisfies the ApiClient interface, with placeholder
// implementations of Check and RunCompletion.
//
// This allows easy implementation of fairly simple clients, while still
// leaving the door open to highly customized clients with different types.
//
// See OpenApiClient for an idiomatic use of BasicApiClient.
type BasicApiClient struct {

	// Client should hold the core client used by CompletionFunc.
	Client any

	ContextItems        []ContextItem
	Tools               []string
	Model               string
	MaxCompletionTokens int

	PreFunc  func(ApiClient, any) error
	PostFunc func(ApiClient, any) error

	Streaming bool
	ShowCalls bool
	PrintFunc func(a ...any)
	Logger    *slog.Logger
	DumpDir   string
}

// SetLogger implements ApiClient.
func (c *BasicApiClient) SetLogger(logger *slog.Logger) {
	c.Logger = logger
}

// SetDumpDir implements ApiClient.
func (c *BasicApiClient) SetDumpDir(dir string) {
	c.DumpDir = dir
}

// SetPrintFunc implements ApiClient.
func (c *BasicApiClient) SetPrintFunc(f func(a ...any)) {
	c.PrintFunc = f
}

// SetPreFunc implements ApiClient.
func (c *BasicApiClient) SetPreFunc(f func(ApiClient, any) error) error {
	c.PreFunc = f
	return nil
}

// SetPostFunc implements ApiClient.
func (c *BasicApiClient) SetPostFunc(f func(ApiClient, any) error) error {
	c.PostFunc = f
	return nil
}

// SetModel implements ApiClient but does *not* validate the model.
func (c *BasicApiClient) SetModel(model string) error {
	c.Model = model
	return nil
}

// SetTools implements ApiClient but does *not* validate the tools.
func (c *BasicApiClient) SetTools(tools []string) error {
	c.Tools = tools
	return nil
}

// SetStreaming implements ApiClient.
func (c *BasicApiClient) SetStreaming(streaming bool) {
	c.Streaming = streaming
}

// SetShowCalls implements ApiClient.
func (c *BasicApiClient) SetShowCalls(show bool) {
	c.ShowCalls = show
}

// ClearContext implements ApiClient.
func (c *BasicApiClient) ClearContext() {
	c.ContextItems = nil
}

// AddContextItem implements ApiClient.
func (c *BasicApiClient) AddContextItem(item ContextItem) {
	c.ContextItems = append(c.ContextItems, item)
}

// SetMaxCompletionTokens implements ApiClient.
func (c *BasicApiClient) SetMaxCompletionTokens(limit int) {
	c.MaxCompletionTokens = limit
}

// RunCompletion implements ApiClient but returns ErrPlaceholder.
func (c *BasicApiClient) RunCompletion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	return nil, ErrPlaceholder
}

// Check implements ApiClient by calling CheckFunc.
func (c *BasicApiClient) Check(ctx context.Context) error {
	return ErrPlaceholder
}
