// Package utils includes miscellaneous utilities to reduce boilerplate code.
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/titanous/json5"
	"golang.org/x/term"
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

// MustUnToml unmarshals b to v or panics trying.
func MustUnToml(b []byte, v any) {
	err := toml.Unmarshal(b, v)
	if err != nil {
		panic(err)
	}
}

// MustUnTomlString unmarshals s to v or panics trying.
func MustUnTomlString(s string, v any) {
	MustUnToml([]byte(s), v)
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

// MustUnJson unmarshals b to v or panics trying.
func MustUnJson(b []byte, v any) {
	err := json.Unmarshal(b, v)
	if err != nil {
		panic(err)
	}
}

// MustUnJsonString unmarshals s to v or panics trying.
func MustUnJsonString(s string, v any) {
	MustUnJson([]byte(s), v)
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

// JsonFile marshals v to JSON and writes the data to file.
func JsonFile(v any, file string) error {

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if err := os.WriteFile(file, b, 0666); err != nil {
		return err
	}
	return nil
}

// JsonFilePretty is as JsonFile but with indent.
func JsonFilePretty(v any, file string) error {

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(file, b, 0666); err != nil {
		return err
	}
	return nil

}

// MustJsonFile calls JsonFile and panics on error.
func MustJsonFile(v any, file string) {

	if err := JsonFile(v, file); err != nil {
		panic(err)
	}
}

// MustJsonFilePretty calls JsonFilePretty and panics on error.
func MustJsonFilePretty(v any, file string) {

	if err := JsonFilePretty(v, file); err != nil {
		panic(err)
	}
}

// TomlFile marshals v to TOML and writes the data to file.
func TomlFile(v any, file string) error {

	b, err := toml.Marshal(v)
	if err != nil {
		return err
	}
	if err := os.WriteFile(file, b, 0666); err != nil {
		return err
	}
	return nil
}

// MustTomlFile calls TomlFile and panics on error.
func MustTomlFile(v any, file string) {

	err := TomlFile(v, file)
	if err != nil {
		panic(err)
	}
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

// IsExecutable checks whether a file is (i.e., appears to be) executable.
func IsExecutable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	mode := info.Mode()
	if !mode.IsRegular() {
		return false, nil
	}

	if runtime.GOOS == "windows" {
		ext := strings.ToLower(filepath.Ext(path))
		return ext == ".exe" || ext == ".bat" || ext == ".cmd" || ext == ".com", nil
	}

	// Check if any execution bit is set
	return mode&0111 != 0, nil
}

var DefaultTerminalWidth = 80

// GetTerminalWidth attempts to obtain the current terminal width for
// formatting output, trying Stdout first and falling back to /dev/tty,
// finally defaulting to DefaultTerminalWidth.
func GetTerminalWidth() int {
	// Try stdout first
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			return w
		}
	}

	// Fall back to /dev/tty
	tty, err := os.Open("/dev/tty")
	if err == nil {
		defer tty.Close()
		if w, _, err := term.GetSize(int(tty.Fd())); err == nil {
			return w
		}
	}

	// Give up and use default
	return 80
}
