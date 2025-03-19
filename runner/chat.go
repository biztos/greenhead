package runner

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
)

// RunChat runs an interactive chat session.
//
// TODO: take the agent as arg here.
func RunChat(agt *agent.Agent) error {

	// Hmm, gonna want a way to suppress logging or log to file here.
	tmp, err := os.CreateTemp("", "wtf-log-*.log")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmp.Close()

	// Get the absolute path to print it for demonstration
	absPath, _ := filepath.Abs(tmp.Name())

	// Create a JSON handler that writes to the temp file
	fileHandler := slog.NewJSONHandler(tmp, &slog.HandlerOptions{
		Level: slog.LevelInfo, // TODO: this stuff for setting debug on init
	})

	// Set the default logger to use our file handler
	agt.SetLogger(slog.New(fileHandler))

	fmt.Println("Chatting with:", agt.String())
	fmt.Println("Logs:", absPath)
	fmt.Println("Return twice to send prompt; empty prompt or Ctrl-D to quit.")
	fmt.Println("Note that context is NOT cleared!")

	rl, err := readline.New("> ")
	if err != nil {
		return fmt.Errorf("readline: %w", err)
	}
	defer rl.Close()

	prompt := ""
	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		prompt += line + "\n"
		if strings.HasSuffix(prompt, "\n\n") {
			prompt = strings.TrimSpace(prompt)
			if prompt == "" {
				break
			}
			_, err = agt.RunCompletion(context.Background(), prompt)
			if err != nil {
				return err
			}
			// Normally you'd want to stream, but in case not then... what?
			// Want to just dump the response message then.
			prompt = ""
		}
	}
	fmt.Println("* DONE")

	return nil
}

// TODO: obviously not do this shit here!  useful for demo though.
type ParseUrlInput struct {
	Url string `json:"url"`
}

func ParseUrl(ctx context.Context, in ParseUrlInput) (*url.URL, error) {
	return url.Parse(in.Url) // context ignored
}
func init() {

	parse_url := tools.NewTool[ParseUrlInput, *url.URL](
		"parse_url",
		"Parses an URL and returns its parts in a struct.",
		ParseUrl,
	)
	if err := registry.Register(parse_url); err != nil {
		panic(err)
	}

}
