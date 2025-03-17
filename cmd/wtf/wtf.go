package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/runner"
)

func Runner(cfg *runner.Config) error {
	agent_cfg := &agent.Config{
		Type:   "openai",
		Name:   "WTF",
		Tools:  []string{"parse_url"},
		Stream: cfg.Stream,
		Color:  "lightblue",
		Context: []*agent.ContextItem{
			{
				Role:    "system",
				Content: "You are a helpful assistant, but you speak like a pirate.",
			},
		},
	}
	my_agent, err := agent.NewAgent(agent_cfg)
	if err != nil {
		return err
	}

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
	my_agent.SetLogger(slog.New(fileHandler))

	fmt.Println("Chatting with:", my_agent.String())
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
			_, err = my_agent.RunCompletion(context.Background(), prompt)
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

func init() {
	runner.RunnerFunc = Runner
}
