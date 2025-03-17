// Package utils includes miscellaneous utilities to reduce boilerplate code.
package utils

import (
	"encoding/json"
	"time"
)

// MustJsonString marshals v to a string or panics trying.
func MustJsonString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

// Dur returns the string representation of the duration since t.
func Dur(t time.Time) string {
	return time.Since(t).String()
}

// DurLog returns a tuple that can be used as an attribute in slog.
//
// NB: this is only useful if you're logging *just* the duration and no other
// attributes.  More complex logging wrappers are TODO; for now we shall stick
// with having durations be their own log entries.
func DurLog(t time.Time) []any {
	return []any{"duration", Dur(t)}
}
