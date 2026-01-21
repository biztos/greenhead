// Package demo defines a Demo extension that registers itself on init.
//
// This is intended as an idiomatic way of creating new tools as subpackages.
//
// TODO: break out into more realistic examples.
//
// Say one that is stateless, one stateful.
//
// Stateful can be something real, like say Employee Registry.
//
// Stateless... maybe uptime? github.com/shirou/gopsutil
package demo

import (
	"context"

	"github.com/biztos/greenhead/ghd/registry"
	"github.com/biztos/greenhead/ghd/tools"
)

// Demo defines the demonstration tool.
type Demo struct {
	memory []string
}

// TODO: generic input types for all the standard single-value things, e.g.
// StringArrayInput, IntArrayInput, AnyArrayInput...

// NullInput defines an input with no payload.
// TODO: see if this actually works in terms of the JSON setup, first for GPT.
type NullInput struct{}

// StoreInput defines the input to Store.
type StoreInput struct {
	Value string `json:"value"`
}

// SumInput defines the input to Sum.
type SumInput struct {
	Values []float64 `json:"values"`
}

// Store stores a string for recall later.
func (d *Demo) Store(val string) {
	d.memory = append(d.memory, val)
}

// Recall returns all the stored strings.
func (d *Demo) Recall() []string {
	return d.memory
}

// Sum returns the sum of all its input values.
func (d *Demo) Sum(vals []float64) float64 {
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum

}

func init() {
	// Wrap the functions.
	//
	// Note that there is only one global object here.
	//
	// Instantiating them on the fly is TBD but should be doable;
	// cf. PROBLEMS.md
	d := &Demo{}
	store := tools.NewTool[StoreInput, string](
		"demo_store",
		"Stores a string value for later Recall.",
		func(ctx context.Context, in StoreInput) (string, error) {
			d.Store(in.Value) // ctx is ignored for now
			return "stored", nil
		},
	)
	recall := tools.NewTool[NullInput, []string](
		"demo_recall",
		"Returns the stored values.",
		func(ctx context.Context, in NullInput) ([]string, error) {
			return d.Recall(), nil // ctx is ignored for now
		},
	)
	sum := tools.NewTool[SumInput, float64](
		"demo_sum",
		"Sum input values.",
		func(ctx context.Context, in SumInput) (float64, error) {
			return d.Sum(in.Values), nil // ctx is ignored for now
		},
	)

	if err := registry.Register(store); err != nil {
		panic(err)
	}
	if err := registry.Register(recall); err != nil {
		panic(err)
	}
	if err := registry.Register(sum); err != nil {
		panic(err)
	}
}
