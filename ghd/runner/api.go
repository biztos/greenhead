package runner

import (
	"fmt"
	"io"

	"github.com/biztos/greenhead/ghd/api"
	"github.com/biztos/greenhead/ghd/utils"
)

// EncodeKeys encodes raw_keys according to the API configuration, printing
// the pairs to w.
func (r *Runner) EncodeKeys(w io.Writer, raw_keys []string) {
	for _, raw := range raw_keys {
		var enc string
		if r.Config.API.RawKeys {
			enc = raw
		} else {
			enc = api.EncodeAuthKey(raw)
		}
		fmt.Fprintln(w, raw, enc)
	}
}

// CheckKey checks an encoded key against configured keys in an API.
// On success, prints the TOML encoding of the key to w.
func (r *Runner) CheckKey(w io.Writer, encoded_key string) error {
	api, err := api.NewAPI(r.Config.API, r.Agents)
	if err != nil {
		return err
	}
	key := api.GetKey(encoded_key)
	if key == nil {
		return fmt.Errorf("key not found: %q", encoded_key)
	}
	fmt.Fprintln(w, utils.MustTomlString(key))
	return nil
}

// CheckAPI instantiates an API and prints "OK" to w if successful.
func (r *Runner) CheckAPI(w io.Writer) error {

	_, err := api.NewAPI(r.Config.API, r.Agents)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, "OK")
	return nil

}

// ServeAPI instantiates an API and calls Listen.
func (r *Runner) ServeAPI(w io.Writer) error {

	api, err := api.NewAPI(r.Config.API, r.Agents)
	if err != nil {
		return err
	}
	return api.Listen()

}
