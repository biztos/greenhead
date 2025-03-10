// Package demo defines a Demo extension that registers itself on init.
package demo

import (
	"context"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/sashabaranov/go-openai"

	"github.com/biztos/greenhead/registry"
)

// Demo defines the demonstration extension.
type Demo struct {
	id       ulid.ULID
	defs     []*openai.FunctionDefinition
	funcmaps map[string]func(ctx context.Context, name string, a ...any) (string, error)
	memory   []string
}

// Name implements extensions.Extension.
func (d *Demo) Name() string {
	return "Demo"
}

// Description implements extensions.Extension.
func (d *Demo) Description() string {
	return "Demonstration of callable functions in an extension."
}

// Functions implements extensions.Extension.
func (d *Demo) Functions() []*openai.FunctionDefinition {
	return d.defs
}

// Call implements extensions.Extension.
func (d *Demo) Call(ctx context.Context, name string, a ...any) (string, error) {

	/*

		plan... don't really want to overuse reflection but might include
		as an alt way to do it?

		map is cumbersome though

		need to do this:

		HAVE METHOD? NO -> ERR

		HAVE DEF? NO -> ERR

		CONFORM CALL ANYS TO DEFINED ARGS FOR THING

		(should we use anys or strings which can be converted?)

		FAIL IF CAN'T CONFORM

		MAKE THE CALL WITH THE CORRECT STUFF

		...kinda... shitty?




		ideally use reflection to get the

		get the def


	*/
	return "TBD", nil
}

// Stop implements extensions.Extension (and is a noop).
func (d *Demo) Stop(ctx context.Context) error {
	return nil
}

// Hello returns a greeting that includes identifying information.
func (d *Demo) Hello() string {
	return fmt.Sprintf("Hello, I'm %s!", d.id)
}

// Store stashes s in memory.
func (d *Demo) Store(s string) {
	d.memory = append(d.memory, s)
}

// Sum returns the sum of the input floats.
func (d *Demo) Sum(vals ...float64) float64 {
	v := 0.0
	for _, val := range vals {
		v += val
	}
	return v
}

// Retrieve returns all that was saved in Store, joined with newlines.
func (d *Demo) Retrieve() string {
	return strings.Join(d.memory, "\n")
}

// NewDemo returns a *Demo with a unique id and wrappers for the Hello,
// Store, Recall, and Sum functions.
func NewDemo() *Demo {
	d := &Demo{id: ulid.Make()}
	d.defs := []*openai.FunctionDefinition{
		{
			Name: "hello",
			// could include the func wrapper itself and json omit
		},
	}

	return d
}

func init() {
	if err := registry.Register(NewDemo()); err != nil {
		panic(err)
	}
}
