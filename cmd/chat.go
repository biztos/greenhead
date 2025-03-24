package cmd

import (
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"

	"github.com/biztos/greenhead/agent"
	"github.com/biztos/greenhead/runner"
)

const ExitCodeChatError = 3

// ChatCmd represents the "chat" command.
var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with an agent.",
	Long: `The chat command starts a chat session with an agent defined by the
provided config file.

TODO: that config, duh.`,
	Run: func(cmd *cobra.Command, args []string) {
		// OK, this is actually a great place to figure out how to do the
		// agent setup from config -- should be EASY to do this for a single
		// agent!
		// TODO obviously, just using the wtf hardcode for now.
		// something like agent.NewFromConfigs(agent-config,runner-config)
		// and get second config easy as can...
		agent_cfg := &agent.Config{
			Type:   "openai",
			Model:  openai.GPT4o,
			Name:   "WTF",
			Tools:  []string{"parse_url"},
			Stream: Config.Stream,
			Color:  "lightblue",
			Context: []agent.ContextItem{
				{
					Role:    "system",
					Content: "You are a helpful assistant, but you speak like a pirate.",
				},
			},
		}
		chat_agent, err := agent.NewAgent(agent_cfg)
		if err != nil {
			BailErr(ExitCodeChatError, err)
		}
		runner.RunChat(chat_agent)
	},
}

func init() {
	// TODO: config file
	RootCmd.AddCommand(ChatCmd)
}
