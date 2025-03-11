// Package tools defines the types for tools (functions) available to the
// LLMs.
package tools

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Tooler defines the interface to which Tools conform.
//
// Tools are managed as Toolers; for complex use-cases you may wish to skip
// Tool[T any] altogether and define your own type.
//
// For simple use-cases, just use NewTool.
type Tooler interface {
	Name() string
	Description() string
	Exec(context.Context, any) (any, error)

	OpenAiTool() (*openai.Tool, error)
}

// Tool is a tool which can be called by LLMs once registered.
//
// T is the input type for the function, and should be JSON serializable.
type Tool[T any] struct {
	name string
	desc string
	f    func(context.Context, T) (any, error)
}

// NewTool returns a Tool for type T.
func NewTool[T any](name, desc string, f func(context.Context, T) (any, error)) *Tool[T] {
	return &Tool[T]{
		name: name,
		desc: desc,
		f:    f,
	}
}

// Name implements Tooler.
func (t *Tool[T]) Name() string {
	return t.name
}

// Description implements Tooler.
func (t *Tool[T]) Description() string {
	return t.desc
}

// Exec implements Tooler by calling Func with input type-checked.
func (t *Tool[T]) Exec(ctx context.Context, input any) (any, error) {
	inp_t, ok := input.(T)
	if !ok {
		// TODO: panic here, probably -- once we have the rigging.
		return nil, fmt.Errorf("Wrong input type for %s: %T", t.name, input)
	}
	return t.f(ctx, inp_t)
}

// OpenAiTool implements Tooler and returns an OpenAI-compatibile definition
// of t.
//
// "Strict" is assumed and the ToolType is always "Function."
func (t *Tool[T]) OpenAiTool() (*openai.Tool, error) {
	oai_tool := &openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        t.name,
			Description: t.desc,
			Strict:      true,
			// TODO: figure out the right way to do parameters here...
			// Presumably want to set up a JSONSchema struct based on T.
			Parameters: nil,
		},
	}
	return oai_tool, nil
}
