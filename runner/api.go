package runner

import (
	"fmt"
	"io"

	"github.com/biztos/greenhead/api"
)

// CheckAPI instantiates an API and writes a log to its logger, printing OK
// to w if successful.
func (r *Runner) CheckAPI(w io.Writer) error {

	api, err := api.NewAPI(r.Config.API, r.Agents)
	if err != nil {
		return err
	}
	// TODO: logs!
	fmt.Println(api)
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
