package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/runner"
)

func Runner(cfg *runner.Config) error {
	// ignore all the actual opts, just want to get a r/t going here.
	//
	// FIRST do this with the fake weather tool, before trying to call 4reals.
	prompt := "Parse the following URL: https://biztos.com/misc/?a=b"
	res, err := ExecuteChat(CompletionRequest(prompt), cfg.Stream)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

// Now do this with a real tool!
func GetTool(name string) openai.Tool {

	fmt.Println("USING TOOL:", name)
	tool := registry.Get(name)
	if tool == nil {
		panic("tool not found") // don't care, IRL we will have tools already.
	}
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        tool.Name(),
			Description: tool.Description(),
			Strict:      true,
			Parameters:  tool.InputSchema(),
		},
	}

}

func CompletionRequest(prompt string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		// MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Tools: []openai.Tool{GetTool("parse_url")},
	}
}

// ExecuteChat executes a chat completion request and returns a ChatCompletionResponse.
// If stream is true, it streams the response and prints content as it arrives.
// The returned ChatCompletionResponse is the same whether streaming is used or not.
func ExecuteChat(r openai.ChatCompletionRequest, streaming bool) (*openai.ChatCompletionResponse, error) {

	// Create OpenAI client
	token := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(token)
	ctx := context.Background()

	if !streaming {
		// Non-streaming path
		res, err := client.CreateChatCompletion(ctx, r)
		return &res, err
	}

	// Streaming path
	r.Stream = true
	r.StreamOptions = &openai.StreamOptions{IncludeUsage: true}

	stream, err := client.CreateChatCompletionStream(ctx, r)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	// Build up the response as we receive chunks
	var fullResponse = &openai.ChatCompletionResponse{}
	var contentBuilder strings.Builder
	var finishReason string
	var role string
	var toolCalls []openai.ToolCall
	var got_usage bool
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		// If this is the first chunk, initialize the response
		if fullResponse.ID == "" {
			fullResponse.ID = response.ID
			fullResponse.Object = response.Object // TODO: remove ".chunk"
			fullResponse.Created = response.Created
			fullResponse.Model = response.Model
			fullResponse.SystemFingerprint = response.SystemFingerprint
		}

		// Handle the Usage chunk, which by definition is the last one.
		if got_usage {
			return nil, fmt.Errorf("chunk after usage: %v", response)
		}
		if len(response.Choices) == 0 {
			// This is the usage chunk; there should be only one.
			if response.Usage == nil {
				return nil, fmt.Errorf("unexpected nil Usage: %v", response)
			}
			if fullResponse.Usage.TotalTokens > 0 {
				return nil, fmt.Errorf("dupe usage chunk: %v", response.Usage)
			}

			fullResponse.Usage = *response.Usage
			fmt.Println("")
			got_usage = true
			continue

		}

		// Extract role (typically "assistant") from first delta
		if role == "" && response.Choices[0].Delta.Role != "" {
			role = response.Choices[0].Delta.Role
		}

		// Append the content
		content := response.Choices[0].Delta.Content
		if content != "" {
			contentBuilder.WriteString(content)
			fmt.Print(content) // Print content as it arrives
		}

		// Handle tool calls (if present)
		if len(response.Choices[0].Delta.ToolCalls) > 0 {
			// Process each tool call delta
			for _, toolCallDelta := range response.Choices[0].Delta.ToolCalls {
				// For new tool calls
				if toolCallDelta.Index != nil {

					index := *toolCallDelta.Index

					// If this is a new tool call
					if index >= len(toolCalls) {
						toolCalls = append(toolCalls, openai.ToolCall{
							ID:   toolCallDelta.ID,
							Type: toolCallDelta.Type,
						})
					}

					// Update function details
					if toolCallDelta.Function.Name != "" {
						toolCalls[index].Function.Name = toolCallDelta.Function.Name
						fmt.Printf("Tool call: %s ", toolCallDelta.Function.Name)
					}

					// Append to function arguments as they stream in
					if toolCallDelta.Function.Arguments != "" {
						if toolCalls[index].Function.Arguments == "" {
							toolCalls[index].Function.Arguments = toolCallDelta.Function.Arguments
						} else {
							toolCalls[index].Function.Arguments += toolCallDelta.Function.Arguments
						}

						// Print tool call information as it arrives.
						// This *should* be safe for the index, but what if it's not?
						// TODO: wait until done with a tool maybe? or till index change?
						fmt.Printf("%s", toolCallDelta.Function.Arguments)
					}
				}
			}
		}

		// Save the finish reason if present
		if response.Choices[0].FinishReason != "" {
			finishReason = string(response.Choices[0].FinishReason)
		}
	}

	// Construct the final response to match the non-streaming format
	message := openai.ChatCompletionMessage{
		Role:    role,
		Content: contentBuilder.String(),
	}

	// Add tool calls if any were received
	if len(toolCalls) > 0 {
		message.ToolCalls = toolCalls
	}

	fullResponse.Choices = []openai.ChatCompletionChoice{
		{
			Message:      message,
			FinishReason: openai.FinishReason(finishReason),
			Index:        0,
		},
	}

	// Add a new line after streaming is complete
	fmt.Println()

	return fullResponse, nil
}

func init() {

	runner.RunnerFunc = Runner
}
