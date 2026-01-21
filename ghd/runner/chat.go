package runner

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)

// TODO: (long-term) - by default log internally to the chat and allow the
// examination of logs from within the chat session.  (Really long term!)
var ErrChatRequiresLogFile = errors.New("chat requires a log file")

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

	// Require that we log to file, because logging to output makes the chat
	// unusable.
	if !r.Config.NoLog && r.Config.LogFile == "" {
		return ErrChatRequiresLogFile
	}

	fmt.Println("Chatting with:", agent.String())
	fmt.Println("Logs:", r.Config.LogFile)
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
