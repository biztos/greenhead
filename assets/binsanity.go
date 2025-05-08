/* binsanity.go - auto-generated; edit at your own peril!

More info: https://github.com/biztos/binsanity

*/

package assets

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"sort"
)

// Asset returns the byte content of the asset for the given name, or an error
// if no such asset is available.
func Asset(name string) ([]byte, error) {

	_, found := binsanity_cache[name]
	if !found {
		i := sort.SearchStrings(binsanity_names, name)
		if i == len(binsanity_names) || binsanity_names[i] != name {
			return nil, errors.New("Asset not found.")
		}

		// We ignore errors because we controlled the data from the begining.
		// It's not perfect but it seems better than having additional funcs
		// hanging around that might confuse the user: tried that already, not
		// nicer.
		decoded, _ := base64.StdEncoding.DecodeString(binsanity_data[i])
		buf := bytes.NewReader(decoded)
		gzr, _ := gzip.NewReader(buf)
		defer gzr.Close()
		data, _ := io.ReadAll(gzr)

		// Not cached, so decode and cache it.
		binsanity_cache[name] = data

	}
	return binsanity_cache[name], nil

}

// MustAsset returns the byte content of the asset for the given name, or
// panics if no such asset is available.
func MustAsset(name string) []byte {
	b, err := Asset(name)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// MustAssetString returns the string content of the asset for the given name,
// or panics if no such asset is available.  This is a convenience function
// for string(MustAsset(name)).
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetNames returns the sorted names of the assets.
func AssetNames() []string {
	return binsanity_names
}

// this must remain sorted or everything breaks!
var binsanity_names = []string{
	"doc/config.md",
	"doc/license.md",
	"doc/people.md",
	"doc/readme.md",
}

// only decode once per asset.
var binsanity_cache = map[string][]byte{}

// assets are gzipped and base64 encoded
var binsanity_data = []string{
	"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",
}
