package runner

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
)

// RunChat runs an interactive chat session.
//
// TODO: make fun! cf. TODO.md.
func (r *Runner) RunChat() error {

	// Require exactly one agent, at least for now.
	if len(r.Agents) != 1 {
		return fmt.Errorf("chat requires one agent configured, not %d",
			len(r.Agents))
	}
	agent := r.Agents[0]

	// If log file is not specified, then log to temp file, because running
	// chat and logging to the console at the same time is unusable.
	log_file := r.Config.LogFile
	if log_file == "" {
		log_file = filepath.Join(os.TempDir(), fmt.Sprintf("%s.log", agent.ULID))
		agent.InitLogger(log_file, r.Config.Debug)
	}

	fmt.Println("Chatting with:", agent.String())
	fmt.Println("Logs:", log_file)
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
			_, err = agent.RunCompletionPrompt(prompt)
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
