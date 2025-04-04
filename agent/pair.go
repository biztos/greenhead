package agent

import (
	"context"
	"fmt"
)

var ErrFirstCompletion = fmt.Errorf("error running first completion in pair")
var ErrSecondCompletion = fmt.Errorf("error running second completion in pair")

// CompletionPair represents a round-trip of completions for the two agents.
type CompletionPair struct {
	FirstRequest   *CompletionRequest
	FirstResponse  *CompletionResponse
	SecondRequest  *CompletionRequest
	SecondResponse *CompletionResponse
	Error          error
}

// Pair represents a pair of Agents in conversation.
type Pair struct {
	First          *Agent
	Second         *Agent
	MaxCompletions int

	completed int
}

// NewPair sets up a Pair for the given Agents.
func NewPair(first *Agent, second *Agent) *Pair {
	return nil
}

// Run calls RunCompletions successively until an error is found in the
// CompletionPair.  The conversation begins with prompt.
//
// NOTE: this does *not* output anything beyond what the agents themselves are
// configured to output.  Poor configuration can result in running without
// any limits, forever, until the singularity occurs or the Cloud Mafia blocks
// your account.
func (p *Pair) Run(ctx context.Context, prompt string) error {

	for {
		cp := p.RunCompletions(ctx, prompt)
		if cp.Error != nil {
			return cp.Error
		}
		prompt = cp.SecondResponse.Content
	}

}

// RunCompletions runs two completions for p.  Prompt is sent to p.First, then
// the respnose content from that completion is sent to p.Second.
//
// If either completion returns an error, or reaches p.MaxCompletions, the
// error is contained within the return value.
//
// The context is checked for errors prior to each agent's completion.
func (p *Pair) RunCompletions(ctx context.Context, prompt string) *CompletionPair {

	var err error
	cp := &CompletionPair{}

	// First.
	if err := ctx.Err(); err != nil {
		cp.Error = fmt.Errorf("%w: %w", ErrFirstCompletion, err)
		return cp
	}
	cp.FirstRequest = &CompletionRequest{Content: prompt}
	cp.FirstResponse, err = p.First.RunCompletion(ctx, cp.FirstRequest)
	if err != nil {
		cp.Error = fmt.Errorf("%w: %w", ErrFirstCompletion, err)
		return cp
	}
	if err := p.inc_completed(); err != nil {
		cp.Error = fmt.Errorf("%w: %w", ErrFirstCompletion, err)
		return cp
	}

	// Second.
	if err := ctx.Err(); err != nil {
		cp.Error = fmt.Errorf("%w: %w", ErrSecondCompletion, err)
		return cp
	}
	cp.SecondRequest = &CompletionRequest{Content: cp.FirstResponse.Content}
	cp.SecondResponse, err = p.Second.RunCompletion(ctx, cp.SecondRequest)
	if err != nil {
		cp.Error = fmt.Errorf("%w: %w", ErrSecondCompletion, err)
	}
	if err := p.inc_completed(); err != nil {
		cp.Error = fmt.Errorf("%w: %w", ErrSecondCompletion, err)
	}
	return cp
}

func (p *Pair) inc_completed() error {
	p.completed++
	if p.MaxCompletions > 0 && p.completed >= p.MaxCompletions {
		return fmt.Errorf("%w: %d", ErrMaxCompletions, p.completed)
	}
	return nil
}
