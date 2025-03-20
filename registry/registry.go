// Package registry holds the registered tools.
//
// All public functions are safe to use concurrently.  However, keep in mind
// that it is possible to replace a Tooler at runtime, in which case the next
// call to Get will return a different value.
package registry

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/biztos/greenhead/tools"
)

var mutex = sync.Mutex{}

var lockedForNew = false

var lockedForReplace = false

var registered = map[string]tools.Tooler{}

var ordered_names = []string{}

var ErrNewLocked = errors.New("registry is locked for new tools")
var ErrReplaceLocked = errors.New("registry is locked for replacement tools")

// Register adds a tool, replacing any same-named tool if allowed.  For any
// non-nil return value, the tool will not have been registered.
//
// Take note of the one-way Lock* functions for controlling the registry.
func Register(t tools.Tooler) error {
	mutex.Lock()
	defer mutex.Unlock()
	// Must have a non-blanco name.
	name := t.Name()
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("empty name for tool")
	}
	if registered[name] != nil {
		// We are replacing, if allowed.
		if lockedForReplace {
			return ErrReplaceLocked
		}
		idx := slices.Index(ordered_names, name)
		ordered_names = slices.Delete(ordered_names, idx, idx+1)
	} else if lockedForNew {
		return ErrNewLocked
	}

	// The names now have the new tool at the end in all cases.
	ordered_names = append(ordered_names, name)
	registered[name] = t
	return nil

}

// Get returns a Tooler by name, or nil if none is found.
func Get(name string) tools.Tooler {
	mutex.Lock()
	defer mutex.Unlock()
	return registered[name]
}

// Names returns all the registered Tooler names, in order of registration.
func Names() []string {
	mutex.Lock()
	defer mutex.Unlock()
	return slices.Clone(ordered_names)

}

// DisplayLines returns the name and description for all registered Toolers,
// with formatting, in order of registration.  Sorting the return strings
// alphabetically will sort by tool name.
//
// If the description is multi-line, the first line is used.
func DisplayLines() []string {
	mutex.Lock()
	defer mutex.Unlock()
	lines := make([]string, len(ordered_names))
	max_name := 0
	descs := make([]string, len(ordered_names))
	for i, name := range ordered_names {
		if len(name) > max_name {
			max_name = len(name)
		}
		desc_lines := strings.Split(registered[name].Description(), "\n")
		descs[i] = desc_lines[0]
	}
	fmt_str := fmt.Sprintf("%%-%ds - %%s", max_name)
	for i, name := range ordered_names {
		lines[i] = fmt.Sprintf(fmt_str, name, descs[i])
	}
	return lines
}

// Clear clears all registered extensions.  Note that Clear does *not* unlock
// the registry.
//
// Clear is intended for the use-case of having runtime tool registration only
// but starting with compiled-in tools you do not want.  Calling Clear in the
// init phase may not do what you expect.
//
// TODO: prove it works in init phase of an *external* package; it should!
func Clear() {
	mutex.Lock()
	defer mutex.Unlock()
	registered = map[string]tools.Tooler{}
	ordered_names = []string{}
}

// Lock locks the registry for both new and replacement tools.
//
// There is no corresponding Unlock function.
//
// Use Lock as a guard against accidentally creating tools at runtime.
//
// Keep in mind that a tool could call system functions; a tool could create
// and register new tools; and an AI could (theoretically) do a "gain of
// function" without your knowledge.
func Lock() {
	mutex.Lock()
	defer mutex.Unlock()
	lockedForNew = true
	lockedForReplace = true
}

// LockForNew applies a selective lock, preventing only new registrations.
//
// If LockForReplace or Lock has been called already, replacements will remain
// locked as well.
func LockForNew() {
	mutex.Lock()
	defer mutex.Unlock()
	lockedForNew = true
}

// LockForReplace applies a selective lock, preventing only replacements.
//
// If LockForNew or Lock has been called already, new registrations will
// remain locked as well.
func LockForReplace() {
	mutex.Lock()
	defer mutex.Unlock()
	lockedForReplace = true
}
