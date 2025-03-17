package main

import (
	"context"
	"fmt"

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
	fmt.Println(my_agent.String())
	// IRL maybe not to check very often, or just have as a CLI option
	if err := my_agent.Client.Check(context.Background()); err != nil {
		return err
	}
	prompt := "Parse the following URLs: https://biztos.com/misc/?a=b, https://google.com/?q=foobar"
	_, err = my_agent.RunCompletion(context.Background(), prompt)
	if err != nil {
		return err
	}
	// should have printed from Stream...
	fmt.Println("* COMPLETED OK")

	return nil
}

func init() {
	runner.RunnerFunc = Runner
}
