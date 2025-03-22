// Package utils includes miscellaneous utilities to reduce boilerplate code.
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// MustJsonString marshals v to a string or panics trying.
func MustJsonString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

// MustJsonStringPretty marshals v to a string with indent or panics trying.
func MustJsonStringPretty(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
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

// Comma-separated list of lowercased file extensions that can be Marshaled
// or Unmarshaled in this package.
const MarshalFileExtensions = ".toml,.json"

var canMarshalExt map[string]bool // set in init, easier to check against

// UnmarshalFile reads file and unmarshals it into v if possible.
//
// TOML is the canonical format, and should be used for any struct tags.
//
// Any other format will be converted to TOML.
//
// For supported formats see ToTomlSupports.
func UnmarshalFile(file string, v any) error {
	ext := strings.ToLower(filepath.Ext(file))
	if !canMarshalExt[ext] {
		return fmt.Errorf("unsupported extension: %q", ext)
	}
	b, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	if ext != ".toml" {
		switch ext {
		case ".toml":
			// already have.
		case ".json":
			x := &map[string]any{}
			if err := json.Unmarshal(b, x); err != nil {
				return fmt.Errorf("error parsing JSON: %w", err)
			}
			b, err = toml.Marshal(x)
			if err != nil {
				return fmt.Errorf("error marshaling TOML: %w", err)
			}
		}

	}
	if err := toml.Unmarshal(b, v); err != nil {
		return fmt.Errorf("error unmarshaling TOML: %w", err)
	}
	return nil
}

// MarshalFile marshals v in the format required and writes it to file.
//
// TOML is the canonical format, and should be used for any struct tags.
//
// Any other format will be converted from the TOML.
//
// For supported formats see ToTomlSupports.
func MarshalFile(v any, file string) error {
	ext := strings.ToLower(filepath.Ext(file))
	if !canMarshalExt[ext] {
		return fmt.Errorf("unsupported extension %s", ext)
	}
	return fmt.Errorf("TODO")
}

func init() {
	canMarshalExt = map[string]bool{}
	for _, s := range strings.Split(MarshalFileExtensions, ",") {
		canMarshalExt[s] = true
	}
}
