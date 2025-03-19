// Package registry holds the registered tools.
//
// Register and Clear are *not* safe for concurrency: tools should be
// registered at load.
package registry

import (
	"fmt"
	"strings"

	"github.com/biztos/greenhead/tools"
)

var locked = false

var registered []tools.Tooler

// Get returns a Tooler by name, or nil if none is found.
func Get(name string) tools.Tooler {
	for _, t := range registered {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

// Names returns all the registered Tooler names, in order of registration.
func Names() []string {
	names := make([]string, len(registered))
	for i, t := range registered {
		names[i] = t.Name()
	}
	return names

}

// Display returns the name and description for all registered Toolers, with
// formatting, or "<no tools>" if none are registered.
//
// If the description is multi-line, the first line is used.
func Display() string {
	if len(registered) == 0 {
		return "<no tools>"
	}
	max_name := 0
	names := make([]string, len(registered))
	descs := make([]string, len(registered))
	for i, t := range registered {
		names[i] = t.Name()
		if len(names[i]) > max_name {
			max_name = len(names[i])
		}
		desc_lines := strings.Split(t.Description(), "\n")
		descs[i] = desc_lines[0]

	}
	fmt_str := fmt.Sprintf("%%-%ds - %%s\n", max_name)
	disp := ""
	for i, n := range names {
		disp += fmt.Sprintf(fmt_str, n, descs[i])
	}
	return disp

}

// Register adds a Tool, with simple checks.  For any non-nil return value, t
// will *not* have been registered.
func Register(t tools.Tooler) error {
	if locked {
		return fmt.Errorf("registry is locked")
	}
	// Must have a non-blanco name.
	if strings.TrimSpace(t.Name()) == "" {
		return fmt.Errorf("empty name for tool")
	}
	// Good enough for now!
	// TODO: check the input for JSON-ability here, better than later.
	registered = append(registered, t)
	return nil

}

// Clear clears all registered extensions.  Calling Clear after doing anything
// with the extensions is *not* supported and is likely to break things.
//
// Use Clear when building a custom binary that should only have access to its
// own extensions.
func Clear() {
	registered = []tools.Tooler{}
}

// Lock marks the registered extensions as locked, so that no more can be
// registered.  It is *important* that the runtime call this function if you
// do not want tools to be able to register new tools.
//
// There is no Unlock function.
//
// Keep in mind that a tool could call system functions; a tool could create
// and register new tools; and an AI could (theoretically) do a "gain of
// function" without your knowledge.
func Lock() {
	locked = true
}

func init() {
	Clear()
}
