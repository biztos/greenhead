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

// ApiKeyCmd represents the "api key" subcommand.
var ApiKeyCmd = &cobra.Command{
	Use:   "key <encoded_key>",
	Short: "Check the encoded key against configured keys.",
	Long: `Checks the encoded key and prints the key data if found.

If the key is not found, an error is printed.

Note that a full config is required in order for this to work, as an API must
be instantiated before its keys can be searched.

If --raw-keys is set then the key should not be encoded.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.CheckKey(Stdout, args[0])
	},
}

// ApiEncodeCmd represents the "api encode" subcommand.
var ApiEncodeCmd = &cobra.Command{
	Use:   "encode <raw_key>...",
	Short: "Encode raw keys to make client-facing hashes..",
	Long: `Prints the auth-key encoded form of each argument.

The raw key and the encoded key are printed as space-delimited pairs:

	first-key niS1U1bvKhIatf0aPuijyRkeg384U-OZo2npBE_xg_c
	second-key 85fyYKJ1zE1C55ZcVWFnvz8Giu1JHTWo9bOMjC25a7A

This is useful for manually configuring API Keys.  No special configuration is
required.

Note that the application may override the encoder; and that if --raw-keys is
set, encoding returns the input string unchanged.

Whitespace in keys (encoded or raw) is not explicitly supported, and will make
the output difficult to parse.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		r.EncodeKeys(Stdout, args)
		return nil
	},
}

// ApiCheckCmd represents the "api check" subcommand.
var ApiCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Instantiate and validate the API without serving.",
	Long: `Sets up an API instance.

This is useful for checking configs, particularly roles and keys.

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
address.

Note that the normal output controls for streaming and logging do not apply
and are ignored if set.

The --log-fiber option is recommended for local testing, as the default
request logs are meant to be machine-readable.

It is STRONGLY recommended that the --no-keys option, and its configuration
equivalent, be used for testing only.`,
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
	ApiCmd.PersistentFlags().StringVar(&Config.API.ListenAddress, "listen", ":3030",
		"Address at which to listen for requests.")
	ApiCmd.PersistentFlags().BoolVar(&Config.API.RawKeys, "raw-keys", false,
		"Do NOT encode API keys: clients must use the original raw key.")
	ApiCmd.PersistentFlags().BoolVar(&Config.API.NoKeys, "no-keys", false,
		"Do NOT require API keys.")
	ApiCmd.PersistentFlags().BoolVar(&Config.API.NoUI, "no-ui", false,
		"Do NOT expose the web UI.")
	ApiCmd.PersistentFlags().BoolVar(&Config.API.LogFiber, "log-fiber", false,
		"Use Fiber logging for API requests.")

	// Registration:
	ApiCmd.AddCommand(ApiKeyCmd)
	ApiCmd.AddCommand(ApiEncodeCmd)
	ApiCmd.AddCommand(ApiCheckCmd)
	ApiCmd.AddCommand(ApiServeCmd)
	RootCmd.AddCommand(ApiCmd)
}
