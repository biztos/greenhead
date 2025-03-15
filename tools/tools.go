// Package tools defines the types for tools (functions) available to the
// LLMs.
//
// NOTE: this currently uses go-openai's jsonschema which is minimalistic and
// *might* break for more complex input types.  If so we can fall back to
// github.com/invopop/jsonschema, however that package gives such complete
// schemas that the LLMs might choke on them.
//
// NOTE: as of this writing, ChatGPT does *not* support JSON schemas for tool
// output.  It expects a string, which could be JSON -- and if it is, the
// LLM will try to interpret it based on the obviousness of the properties.
// Not ideal, but it's what we're stuck with.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// Tooler defines the interface to which Tools conform.
//
// Tools are managed as Toolers; for complex use-cases you may wish to skip
// Tool[T, R] altogether and define your own type.
//
// For simple use-cases, just use NewTool.
type Tooler interface {

	// Name returns the immutable name of the tool, which will be used both
	// to tell the LLM what the tool is called, and to look it up when the
	// LLM returns a tool call.
	Name() string

	// Description is the description of the tool as passed to the LLM.
	// It should be tuned to that purpose.
	Description() string

	// InputSchema returns a simplified JSON schema describing the underlying
	// concrete type of the input (args).
	InputSchema() *jsonschema.Definition

	// Exec executes the tool.  The input string must be valid JSON conforming
	// to the schema returned by InputSchema.
	Exec(ctx context.Context, input string) (any, error)

	// Help returns information on the tool for a human user of the CLI.
	Help() string

	// OpenAiTool returns a valid openai.Tool.
	//
	// TODO: move this!  Since it's derived from the above things, no sense in
	// making someone implement it.  Put in apis maybe?  But we may have our
	// own API at some point...
	OpenAiTool() openai.Tool
}

// Tool is a tool which can be called by LLMs once registered.
//
// T is the input type for the function; R is the return type for the
// non-error value.  Both must be JSON serializable/deserializable or
// runtime errors will occur.
type Tool[T any, R any] struct {
	name    string
	desc    string
	f       func(context.Context, T) (R, error)
	zeroT   T // arguably only need the schemas but keep around for now.
	zeroR   R // ...because perhaps useful for error messages etc.
	schemaT *jsonschema.Definition
}

// NewTool returns a Tool for input type T and output type R.
func NewTool[T any, R any](name, desc string, f func(context.Context, T) (R, error)) *Tool[T, R] {
	var zeroT T
	var zeroR R
	schemaT, err := jsonschema.GenerateSchemaForType(zeroT)
	if err != nil {
		panic(fmt.Sprintf("Input Schema for %s %T: %s", name, zeroT, err))
	}
	return &Tool[T, R]{
		name:  name,
		desc:  desc,
		f:     f,
		zeroT: zeroT,
		zeroR: zeroR,
		// Reflect once, we will be handing these out like candy later.
		// (Although that also involves reflection, so... bench it someday.)
		schemaT: schemaT,
	}
}

// Name implements Tooler.
func (t *Tool[T, R]) Name() string {
	return t.name
}

// Description implements Tooler.
func (t *Tool[T, R]) Description() string {
	return t.desc
}

// InputSchema implements Tooler.
func (t *Tool[T, R]) InputSchema() *jsonschema.Definition {
	return t.schemaT
}

// Exec implements Tooler by calling its function with args as a JSON string.
func (t *Tool[T, R]) Exec(ctx context.Context, args string) (any, error) {
	var input T
	err := json.Unmarshal([]byte(args), &input)
	if err != nil {
		// This could be programmer/prompter error or a hallucination; at
		// least openAI docs *say* the JSON schema should be respected.
		return nil, fmt.Errorf("error parsing json for %T: %w", input, err)
	}
	return t.f(ctx, input)
}

// Help implements Tooler by returning a (hopefully) useful summary of the
// Tool.
func (t *Tool[T, R]) Help() string {
	s := fmt.Sprintf("%s\n\n%s\n\n", t.name, t.desc)
	b, err := json.MarshalIndent(t.schemaT, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	s += fmt.Sprintf("Input Schema:\n\n%s\n\n", string(b))
	s += fmt.Sprintf("Return Type: %T, error\n\n", t.zeroR)

	return s
}

// OpenAiTool implements Tooler and returns an OpenAI-compatible definition
// of t.
//
// "Strict" is assumed and the ToolType is always "Function."
func (t *Tool[T, R]) OpenAiTool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        t.name,
			Description: t.desc,
			Strict:      true, // TODO: what does this mean?
			// TODO: prove this works, it *should* be good to go.
			Parameters: t.schemaT,
		},
	}
}
