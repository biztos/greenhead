package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/runner"
)

// WtfCmd is for work in progress.
var WtfCmd = &cobra.Command{
	Hidden: true,
	Use:    "wtf",
	Short:  "WTF? WTF!",
	Long: `What.
The.
(Actual.)
F**K?!

Work-in-progress random experiments.  Use at your own risk!

TODO: equivalent of this, but in an external package.  Hard? Easy?

Maybe just hook into pre/postrun funcs?`,
	RunE: func(cmd *cobra.Command, args []string) error {

		r, err := runner.NewRunner(Config)
		if err != nil {
			return err
		}
		r.Logger.Info("WTF info")
		r.Logger.Debug("bugggggin'")
		r.Logger.Warn("Open your eyes!", "vision", 1.245)
		r.Logger.Error("you done it now", "bomb", "BOOM")

		// Default and with?
		d := slog.Default()
		d.Info("i am default", "this", "that")
		sub := d.With("with-some", "ting")
		sub.Info("my message", "this", "that")

		if len(r.Agents) > 0 {

			a := r.Agents[0]
			a.Logger().Info("I am agentic!", "this", "that",
				"pretty-girl", "mantra")
		}

		fmt.Println("EHLO WTF")

		return nil
	},
}

func init() {

	RootCmd.AddCommand(WtfCmd)

	// well at least it's clear how we test this...
	runner.NamedAgentConfigs["chatty"].Color = "red" // e.g. "break me"

}
