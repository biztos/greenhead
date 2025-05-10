package runner

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/glamour"

	"github.com/biztos/greenhead/assets"
	"github.com/biztos/greenhead/utils"
)

// DocConfig configures the documentation rendering.
//
// Earlier fields take precedence if true.  The default is the Glamour
// AutoStyle with wrap-width corrected based on screen width.
//
// This is used instead of a pseudo-enum in order to facilitate flags in the
// Cobra command.
type DocConfig struct {
	Markdown bool
	Html     bool
	Ascii    bool
}

// PrintDocs prints the documentation for the provided topic in the preferred
// format.
//
// If topic is "topics", prints a list of topics similar to help text, also in
// the preferred format.
func PrintDocs(w io.Writer, topic string, cfg *DocConfig) error {

	var src string
	var err error
	if topic == "topics" {
		src = "# Documentation Topics\n\n"
		src += strings.Join(DocTopicList(), "\n")
		src += "\n"
	} else {
		name := "doc/" + topic + ".md"
		src, err = assets.AssetString(name)
		if errors.Is(err, assets.ErrNotFound) {
			return fmt.Errorf("topic not found: %s", topic)
		}
		if err != nil {
			return err
		}
	}

	// Special case for source format:
	if cfg.Markdown {
		fmt.Fprintln(w, src)
		return nil
	}

	// HTML doesn't use Glamour for rendering:
	if cfg.Html {
		return fmt.Errorf("HTML output is TODO!")
	}

	// The rest use Glamour.
	wrap := glamour.WithWordWrap(utils.GetTerminalWidth())
	var style glamour.TermRendererOption
	if cfg.Ascii {
		style = glamour.WithStylePath("ascii")
	} else {
		style = glamour.WithAutoStyle()
	}
	renderer, err := glamour.NewTermRenderer(style, wrap)
	if err != nil {
		// TODO: don't capture error if it's impossible to trigger! See if we
		// can break it with the wrap function though.
		return fmt.Errorf("error creating renderer: %w", err)
	}
	out, err := renderer.Render(src)
	if err != nil {
		return fmt.Errorf("error rendering topic: %w", err)
	}
	fmt.Fprintln(w, out)

	return nil
}

// DocTopicList returns a set of Markdown-formatted unordered list items of
// available documentation topics.
func DocTopicList() []string {

	// Get the topics and titles.
	topics := assets.PrefixNamesExt("doc/", ".md", true)
	titles := make([]string, len(topics))
	for i, topic := range topics {
		// assume well-behaved source here, we do control it after all.
		h1, _ := assets.Header("doc/"+topic+".md", 1)
		titles[i] = strings.TrimPrefix(h1, "# ")
	}

	// Get length for formatting.
	max_topic := 0
	for _, topic := range topics {
		if len(topic) > max_topic {
			max_topic = len(topic)
		}
	}

	f := fmt.Sprintf("* %%-%ds - %%s", max_topic)
	list := make([]string, len(topics))
	for i, topic := range assets.PrefixNamesExt("doc/", ".md", true) {
		list[i] = fmt.Sprintf(f, topic, titles[i])

	}
	return list

}
