// agent/client.go

package agent

import (
	"context"
	"errors"
	"log/slog"
)

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

	Streaming       bool
	StreamToolCalls bool
	PrintFunc       func(a ...any)
	Logger          *slog.Logger
}

// SetLogger implements ApiClient.
func (c *BasicApiClient) SetLogger(logger *slog.Logger) {
	c.Logger = logger
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

// SetStreamToolCalls implements ApiClient.
func (c *BasicApiClient) SetStreamToolCalls(streaming bool) {
	c.StreamToolCalls = streaming
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
