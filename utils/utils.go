// Package utils includes miscellaneous utilities to reduce boilerplate code.
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/titanous/json5"
)

// MustToml marshals v to a byte array or panics trying.
func MustToml(v any) []byte {
	b, err := toml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// MustTomlString marshals v to a string or panics trying.
func MustTomlString(v any) string {
	return string(MustToml(v))

}

// MustJson marshals v to a byte array or panics trying.
func MustJson(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// MustJsonString marshals v to a string or panics trying.
func MustJsonString(v any) string {
	return string(MustJson(v))
}

// MustJsonPretty marshals v to a byte array with indent or panics trying.
func MustJsonPretty(v any) []byte {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	return b
}

// MustJsonStringPretty marshals v to a string with indent or panics trying.
func MustJsonStringPretty(v any) string {
	return string(MustJsonPretty(v))
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

// MarshalJsonFile marshals v to JSON and writes the data to file.
func MarshalJsonFile(v any, file string) error {

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if err := os.WriteFile(file, b, 0666); err != nil {
		return err
	}
	return nil
}

// MustMarshalJsonFile calls MarshalJsonFile and panics on error.
func MustMarshalJsonFile(v any, file string) {

	if err := MarshalJsonFile(v, file); err != nil {
		panic(err)
	}
}

// MarshalTomlFile marshals v to TOML and writes the data to file.
func MarshalTomlFile(v any, file string) error {

	b, err := toml.Marshal(v)
	if err != nil {
		return err
	}
	if err := os.WriteFile(file, b, 0666); err != nil {
		return err
	}
	return nil
}

// Comma-separated list of lowercased file extensions that can be Marshaled
// or Unmarshaled in this package.
const UnmarshalFileExtensions = ".toml,.json,.json5"

// UnmarshalFile reads file and unmarshals it into v if possible.
//
// Intended for reading config files, the default format is TOML.  Any other
// types in UnmarshalFileExtensions will be converted to TOML first.
func UnmarshalFile(file string, v any) error {
	ext := strings.ToLower(filepath.Ext(file))
	if !slices.Contains(strings.Split(UnmarshalFileExtensions, ","), ext) {
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
		case ".json", ".json5":
			x := &map[string]any{}
			if err := json5.Unmarshal(b, x); err != nil {
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
