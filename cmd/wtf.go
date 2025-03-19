package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	// "github.com/biztos/greenhead/runner"
)

const ExitCodeWTF = 999

// WtfCmd is for work in progress.
var WtfCmd = &cobra.Command{
	Use:   "wtf",
	Short: "WTF? WTF!",
	Long: `What.
The.
(Actual.)
F**K?!`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HHELO WTTTF")
	},
}

func init() {
	RootCmd.AddCommand(WtfCmd)
}
