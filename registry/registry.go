// Package registry holds the registered tools.
//
// All public functions are safe to use concurrently.  However, keep in mind
// that it is possible to replace a Tooler at runtime, in which case the next
// call to Get will return a different value.
package registry

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/biztos/greenhead/tools"
)

var mutex = sync.Mutex{}

var lockedForNew = false

var lockedForReplace = false

var lockedForRemove = false

var registered = map[string]tools.Tooler{}

var ordered_names = []string{}

var ErrNotRegistered = errors.New("tool is not registered")
var ErrNewLocked = errors.New("registry is locked for new tools")
var ErrReplaceLocked = errors.New("registry is locked for replacement tools")
var ErrRemoveLocked = errors.New("registry is locked for removal")

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

// Remove removes a tool registration if it is registered.
func Remove(name string) error {
	mutex.Lock()
	defer mutex.Unlock()
	if registered[name] == nil {
		return fmt.Errorf("%w: %q", ErrNotRegistered, name)
	}
	if lockedForRemove {
		return ErrRemoveLocked
	}
	delete(registered, name)
	idx := slices.Index(ordered_names, name)
	ordered_names = slices.Delete(ordered_names, idx, idx+1)
	return nil
}

// Get returns a Tooler by name, or an ErrNotRegistered error if none is found.
func Get(name string) (tools.Tooler, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if registered[name] == nil {
		return nil, fmt.Errorf("%w: %q", ErrNotRegistered, name)
	}
	return registered[name], nil
}

// Names returns all the registered Tooler names, in order of registration.
func Names() []string {
	mutex.Lock()
	defer mutex.Unlock()
	return slices.Clone(ordered_names)

}

// MatchingNames checks names for validity and returns a deduplicated and
// (in the case of regexp names) expanded set of valid, registered tool names.
//
// If any item in names has no corresponding registered tool, an error is
// returned.
func MatchingNames(names []string) ([]string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	matched_names := make([]string, 0, len(ordered_names))
	have := map[string]bool{}
	for _, n := range names {
		if len(n) >= 2 && n[0] == '/' && n[len(n)-1] == '/' {
			re, err := regexp.Compile(n[1 : len(n)-1])
			if err != nil {
				return nil, fmt.Errorf("invalid tool regexp %q: %w", n, err)
			}
			matched := false
			for _, rn := range ordered_names {
				if re.MatchString(rn) {
					matched = true
					if !have[rn] {
						have[rn] = true
						matched_names = append(matched_names, rn)
					}
				}
			}
			if !matched {
				return nil, fmt.Errorf("no match for tool %q", n)
			}
		} else {
			if !have[n] {
				if registered[n] == nil {
					return nil, fmt.Errorf("%w: %q", ErrNotRegistered, n)
				}
				have[n] = true
				matched_names = append(matched_names, n)
			}

		}

	}
	return matched_names, nil
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

// Lock locks the registry for both new and replacement tools and also for
// removals.
//
// There is no corresponding Unlock function.
//
// Use Lock as a guard against accidentally changing the toolset at runtime.
//
// Keep in mind that a tool could call system functions; a tool could create
// and register new tools; and an AI could (theoretically) do a "gain of
// function" without your knowledge.
func Lock() {
	mutex.Lock()
	defer mutex.Unlock()
	lockedForNew = true
	lockedForReplace = true
	lockedForRemove = true
}

// LockForNew applies a selective lock, preventing only new registrations.
//
// Previously-set locks are unaffected.
func LockForNew() {
	mutex.Lock()
	defer mutex.Unlock()
	lockedForNew = true
}

// LockForReplace applies a selective lock, preventing only replacements.
//
// Note that while this is useful for allowing new registrations while
// blocking replacements, doing the opposite is dangerous.  When in doubt,
// just use Lock.
//
// Previously-set locks are unaffected.
func LockForReplace() {
	mutex.Lock()
	defer mutex.Unlock()
	lockedForReplace = true
}

// LockForRemove applies a selective lock, preventing only removals.
//
// Previously-set locks are unaffected.
func LockForRemove() {
	mutex.Lock()
	defer mutex.Unlock()
	lockedForRemove = true
}
