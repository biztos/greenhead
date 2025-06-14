package cmd

import (
	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/runner"
)

// ChatCmd represents the "chat" command.
var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with an agent.",
	Long: `The chat command starts a chat session with an agent defined by the
provided config file(s).

Exactly one agent must be configured.

Logs will be written to "chat.log" by default.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if !Config.NoLog && Config.LogFile == "" {
			Config.LogFile = "chat.log"
		}
		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		return r.RunChat()

	},
}

func init() {
	// TODO: config file
	RootCmd.AddCommand(ChatCmd)
}
