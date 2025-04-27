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
	"help/config.md",
}

// only decode once per asset.
var binsanity_cache = map[string][]byte{}

// assets are gzipped and base64 encoded
var binsanity_data = []string{
	"H4sIAAAAAAAA/2xUwW7bOhC88ysG1uU9IJYdtKcEPgQ5tLcCSW5BGjHSSiJMcQlyZcd/X5CU2zjtyfRwODPcXarCPbveDHPQYtgpdfEXJkIL5MiwdCAbbyAjIczOUcB/wr7g/0O7Lm/pgZzUuBMVRQeZ/VWGfeAh6AmBdBfR8jRp162tcQT2ySniOJp2hOh9YlNLHbmWwAcKqs2ZzszbpOig3ekcZNn3gQ+mow5HI2O2bdbrstcsh69UsZn0CUYi2T4dFm1cSb5ILdePt/levXHa2lN2/MhSvbEUET21pjeffDPxbFsr9fAxavznPS/EY43Ummg6CpDRREQzeUuK3nX6BffQaAq5Fp5sc6NU0zRpqSo8zp7CPQePb4HIjaS7pdXK8vCaomOH1eagw8bysBnOrNrysFKTfn9tORmV7uxwvd1mVJhtO6aKJUw9P+fU8eVFAU5PWfVpJDyYrrMUVgroKLbB5Eqk3btzsWXUAmv2FCEM0XYP4xDywQgjGEgi+sBTahaSc6yToJx89mFPTpuETNyRTdDgZf2VMynRscPzavNzIHktupvVCyoY19q5I7Czp9yxgw6G57iYL1ZJ968bbxf0sjpfUumV+k6B/rwDHI21eMsPpgxH8+nk7nq7bfKQNRdOu+tto96o1XMsesJ+nZ/aedhHHSEjxxTdzhRrpaoKy5SVRscEVfgxi58lYRLYFuxR9ySnsn5itngkS235AlQV7nL63yq/AgAA//9jiUl/JQQAAA==",
}
