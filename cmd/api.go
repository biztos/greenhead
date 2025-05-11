package cmd

import (
	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/runner"
)

// ApiCmd represents the "api" command set.
var ApiCmd = &cobra.Command{
	Use:   "api [run|check]",
	Short: "Manage the HTTP API.",
}

// ApiCheckCmd represents the "api check" subcommand.
var ApiCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Instantiate and validate the API without serving..",
	Long: `Sets up an API instance and writes a line to its logger.

Does NOT begin listening for requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.CheckAPI(Stdout)
	},
}

// ApiServeCmd represents the "api serve" subcommand.
var ApiServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the API server.",
	Long: `Sets up an API instance and listens to requests at the configured
address.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		// TODO: make sure no error message from this return style if we
		// cleanly kill the server.  Or we want an error?
		return r.ServeAPI(Stdout)
	},
}

func init() {
	// Flags:
	ApiCmd.PersistentFlags().StringVar(&Config.ApiListenAddress, "address", ":3000",
		"Address at which to listen for requests.")

	// Registration:
	ApiCmd.AddCommand(ApiCheckCmd)
	ApiCmd.AddCommand(ApiServeCmd)
	RootCmd.AddCommand(ApiCmd)
}
