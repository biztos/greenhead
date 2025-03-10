// Package extensions defines the extension interface and has subpackages
// containing the built-in extensions.
package extensions

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// Extension presents a set of callable functions available to an LLM Agent.
type Extension interface {
	Name() string
	Description() string
	Functions() []*openai.FunctionDefinition
	Call(ctx context.Context, name string, a ...any) (string, error)
	Stop(ctx context.Context) error
}
