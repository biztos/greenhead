package cmd

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/assets"
)

var docAsHtml bool
var docAsMarkdown bool
var docAsAscii bool

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

Run "doc" with no topic for a list of topics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("only one topic at a time please")
		}
		if len(args) == 0 {
			fmt.Fprintln(Stdout, "Documentation topics:")
			fmt.Fprintln(Stdout, "")
			for i, topic := range assets.PrefixNamesExt("doc/", ".md", true) {
				fmt.Fprintf(Stdout, "%d.  %s\n", i+1, topic)
			}
			fmt.Fprintln(Stdout, "")
			fmt.Fprintln(Stdout, "Use `doc <topic>` for the doc page.")
		} else {
			name := "doc/" + args[0] + ".md"
			md, err := assets.AssetString(name)
			if errors.Is(err, assets.ErrNotFound) {
				return fmt.Errorf("topic not found: %s", args[0])
			}
			if err != nil {
				return err
			}

			// Source-dump trumps look bump.
			if docAsMarkdown {
				fmt.Fprintln(Stdout, md)
				return nil
			}

			// HTML is TODO but presumably use goldmark
			if docAsHtml {
				return fmt.Errorf("HTML TODO")
			}

			// ASCII is just another Glamour style.
			if docAsAscii {
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
	DocCmd.Flags().BoolVar(&docAsMarkdown, "md", false,
		"Print documentation as Markdown source.")
	DocCmd.Flags().BoolVar(&docAsHtml, "html", false,
		"Print documentation as HTML.")
	DocCmd.Flags().BoolVar(&docAsAscii, "ascii", false,
		"Print documentation as ASCII text without ANSI styling.")

	RootCmd.AddCommand(DocCmd)
}
