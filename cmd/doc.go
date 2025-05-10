package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/runner"
)

var docConfig = &runner.DocConfig{}

// DocCmd represents the "doc" command set.
//
// Doc files are found under assets/src/doc and should be Markdown files;
// however care should be taken to not make them too Markdown-y as they will
// usually be printed to the screen.
//
// TODO: run all the doc stuff in CI to make sure it doesn't error on bad
// Markdown.
var DocCmd = &cobra.Command{
	Use:   "doc [topic]",
	Short: "Show detailed documentation.",
	Long: `The doc command prints detailed documentation for the given topic.

For example, "doc config" prints information about configuration, including
sample config TOML.

Note that documentation here is more comprehensive than the command help,
with the latter focused on running the commands.

Run "doc topics" for a list of topics.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Be liberal with the args, otherwise it gets confusing.
		if len(args) > 1 {
			return errors.New("One topic at a time please.")
		}
		var topic string
		if len(args) == 0 {
			topic = "topics"
		} else {
			topic = args[0]
		}

		return runner.PrintDocs(Stdout, topic, docConfig)
	},
}

func init() {
	DocCmd.Flags().BoolVar(&docConfig.Markdown, "md", false,
		"Print documentation as Markdown source.")
	DocCmd.Flags().BoolVar(&docConfig.Html, "html", false,
		"Print documentation as HTML.")
	DocCmd.Flags().BoolVar(&docConfig.Ascii, "ascii", false,
		"Print documentation as ASCII text without ANSI styling.")

	RootCmd.AddCommand(DocCmd)
}
