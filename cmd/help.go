package cmd

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/assets"
)

var helpAsHtml bool
var helpAsMarkdown bool
var helpAsAscii bool

// HelpCmd represents the "help" command set.
//
// Help files are found under assets/src/help and should be Markdown files;
// however care should be taken to not make them too Markdown-y as they will
// usually be printed to the screen.
//
// TODO: run all the help stuff in CI to make sure it doesn't error on bad
// Markdown.
var HelpCmd = &cobra.Command{
	Use:   "help [topic]",
	Short: "Show detailed help text.",
	Long: `The help command prints detailed help text for provided topic.

For example, "help config" prints information about configuration, including
sample config TOML.

Run "help" with no topic for a list of topics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("only one topic at a time please")
		}
		if len(args) == 0 {
			fmt.Fprintln(Stdout, "Help topics:")
			fmt.Fprintln(Stdout, "")
			for i, topic := range assets.PrefixNamesExt("help/", ".md", true) {
				fmt.Fprintf(Stdout, "%d.  %s\n", i+1, topic)
			}
		} else {
			name := "help/" + args[0] + ".md"
			md, err := assets.AssetString(name)
			if errors.Is(err, assets.ErrNotFound) {
				return fmt.Errorf("topic not found: %s", args[0])
			}
			if err != nil {
				return err
			}

			// Source-dump trumps look bump.
			if helpAsMarkdown {
				fmt.Fprintln(Stdout, md)
				return nil
			}

			// HTML is TODO but presumably use goldmark
			if helpAsHtml {
				return fmt.Errorf("HTML TODO")
			}

			// ASCII is just another Glamour style.
			if helpAsAscii {
				out, err := glamour.Render(md, "ascii")
				if err != nil {
					return fmt.Errorf("error rendering topic: %w", err)
				}
				fmt.Fprintln(Stdout, out)
				return nil
			}

			// By default we want to auto-style.
			r, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
			)
			if err != nil {
				return fmt.Errorf("error creating renderer: %w", err)
			}
			out, err := r.Render(md)
			if err != nil {
				return fmt.Errorf("error rendering topic: %w", err)
			}
			fmt.Fprintln(Stdout, out)

		}
		return nil
	},
}

func init() {
	HelpCmd.Flags().BoolVar(&helpAsMarkdown, "md", false,
		"Print help text as Markdown source.")
	HelpCmd.Flags().BoolVar(&helpAsHtml, "html", false,
		"Print help text as HTML.")
	HelpCmd.Flags().BoolVar(&helpAsAscii, "ascii", false,
		"Print help text as ASCII text without ANSI styling.")

	RootCmd.AddCommand(HelpCmd)
}
