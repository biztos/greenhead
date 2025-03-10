// Package registry holds the registered extensions.
//
// Register and Clear are *not* safe for concurrency: extensions should be
// registered at load.
package registry

import (
	"fmt"
	"strings"

	"github.com/biztos/greenhead/extensions"
)

var registered = []extensions.Extension{}

// Get returns an extension by name, or nil if none is found.
func Get(name string) extensions.Extension {
	for _, e := range registered {
		if e.Name() == name {
			return e
		}
	}
	return nil
}

// Names returns all the registered extension names, in order of registration.
func Names() []string {
	names := make([]string, len(registered))
	for i, e := range registered {
		names[i] = e.Name()
	}
	return names

}

// Register adds an extension, with simple checks.  For any non-nil return
// value, the extension will *not* have been registered.
//
// Note that nothing but your own sanity is preventing you from making the
// values return different things on each call; that kind of chicanery is
// unsupported.
func Register(ext extensions.Extension) error {
	// Must have a non-blanco name.
	if strings.TrimSpace(ext.Name()) == "" {
		return fmt.Errorf("empty name for extension")
	}
	// Must have some functions.
	functions := ext.Functions()
	if len(functions) == 0 {
		return fmt.Errorf("no functions for extension %q", ext.Name())
	}
	// Each of its functions must have a unique name within the extension.
	have := map[string]bool{}
	for _, f := range functions {
		if strings.TrimSpace(f.Name) == "" {
			return fmt.Errorf("empty name for callable in %q - %q",
				ext.Name(), f.Name)
		}
		if have[f.Name] {
			return fmt.Errorf("duplicate name for callable in %q - %q",
				ext.Name(), f.Name)
		}
		have[f.Name] = true
	}
	// Good enough for now!
	registered = append(registered, ext)
	return nil

}

// Clear clears all registered extensions.  Calling Clear after doing anything
// with the extensions is *not* supported and is likely to break things.
//
// Use Clear when building a custom binary that should only have access to its
// own extensions.
func Clear() {
	registered = []extensions.Extension{}
}
