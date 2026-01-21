// Package assets manages compiled-in files.
//
// TODO: bake the following code into binsanity itself.
package assets

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrNotFound = errors.New("asset not found") // TODO: put in binsanity!

// Header returns up to length lines from the top of the named asset.
func Header(name string, length int) (string, error) {
	s, err := AssetString(name)
	if err != nil {
		return "", err
	}
	lines := strings.Split(s, "\n")
	if len(lines) <= length {
		return s, nil
	}
	return strings.Join(lines[:length], "\n"), nil

}

// AsssetString returns the asset as a string, or ErrNotFound.
func AssetString(name string) (string, error) {
	b, err := Asset(name)
	// TODO: proper err check into Asset itself, not this.
	if err != nil {
		return "", ErrNotFound
	}
	return string(b), nil

}

// PrefixNamesExt calls PrefixNames and filters out assets without extension
// ext.  If strip is true, the prefix and extension are removed.
func PrefixNamesExt(prefix string, ext string, strip bool) []string {
	matches := []string{}
	for _, name := range PrefixNames(prefix, strip) {
		if filepath.Ext(name) == ext {
			if strip {
				name = strings.TrimSuffix(name, ext)
			}
			matches = append(matches, name)
		}
	}
	return matches
}

// PrefixNames returns the names of assets with a given prefix, with the
// prefix removed.  If strip is true, the prefix is removed.
//
// NOTE: this may not work in windows-origin builds.
// TODO: make sure binsanity is properly UNIX-ifying paths!
func PrefixNames(prefix string, strip bool) []string {
	matches := []string{}
	for _, name := range AssetNames() {
		if strings.HasPrefix(name, prefix) {
			if strip {
				name = strings.TrimPrefix(name, prefix)
			}
			matches = append(matches, name)
		}
	}
	return matches

}
