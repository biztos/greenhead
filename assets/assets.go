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
// ext.  If strip is true, the extension is removed.
func PrefixNamesExt(prefix string, ext string, strip bool) []string {
	matches := []string{}
	for _, name := range PrefixNames(prefix) {
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
// prefix removed.
//
// NOTE: this may not work in windows-origin builds.
// TODO: make sure binsanity is properly UNIX-ifying paths!
func PrefixNames(prefix string) []string {
	matches := []string{}
	for _, name := range AssetNames() {
		if strings.HasPrefix(name, prefix) {
			matches = append(matches, strings.TrimPrefix(name, prefix))
		}
	}
	return matches

}
