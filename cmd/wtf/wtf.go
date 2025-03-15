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
	// GPT is reliably calling the tool, that's good.  Now maybe do another
	// r/t with the results?
	prompt := "Parse the following URLs: https://biztos.com/misc/?a=b, https://google.com/?q=foobar"
	req := CompletionRequest(prompt)
	res, err := ExecuteChat(req, cfg.Stream)
	if err != nil {
		return fmt.Errorf("error executing chat completion request: %w", err)
	}
	if len(res.Choices) != 1 {
		return fmt.Errorf("unexpected number of choices in response: %d",
			len(res.Choices))
	}
	choice := res.Choices[0]
	tool_calls := choice.Message.ToolCalls
	finish_reason := choice.FinishReason
	if len(tool_calls) > 0 && finish_reason != "tool_calls" {
		return fmt.Errorf("have tool calls but finish_reason is %s",
			finish_reason)
	}

	// TODO: gather async tools and run in a waitgroup, include sync tools
	// IN ORDER as one task in the waitgroup.  mutex for outputs.
	// use a type that holds the mutex? not copied... referred... works?
	outputs := make([]openai.ToolOutput, len(tool_calls))
	for idx, tool_call := range tool_calls {
		if tool_call.Type != openai.ToolTypeFunction {
			return fmt.Errorf("tool call %d has non-function type: %s",
				idx, tool_call.Type)
		}
		name := tool_call.Function.Name
		args := tool_call.Function.Arguments
		// TODO: get from agent's own list of tools b/c could be subset.
		// TBD: what to do if tool not found?  case is the AI hallucinates the
		// tool, so presumably send an error response.  Anyway this should be
		// a method on the agent.
		tool := registry.Get(name)
		if tool == nil {
			panic("tool not found:" + name)
		}
		// we know the type of the tool input?  maybe not.  can call how?
		output, err := tool.Exec(context.Background(), args)
		if err != nil {
			output = map[string]string{"error": err.Error()}
		}
		outputs[idx] = openai.ToolOutput{
			ToolCallID: tool_call.ID,
			Output:     output,
		}
	}
	// TODO: at this point if we do NOT have tool calls then do what? Then
	// we stop the agent and wait for further input which I guess will go
	// under... Continue?
	if len(tool_calls) != 0 {
		// Construct a new call and keep going I guess... could loop on this?
		// Could indeed.  So have to figure out how to handle that.  Probably
		// the agent knows in its config how many tool calls in a row it can
		// get. Easy enough to test with fake responses.  Default 1?  Well.
		// Think about tic tac toe.  Probably depends on the type of call.
		// Maybe build in?
		// OK anyway let's do this as-is for now, just the one time, and dump
		// the response.
		req.Messages = append(req.Messages, choice.Message)
		for _, output := range outputs {
			b, err := json.Marshal(output.Output)
			if err != nil {
				// TODO: prove this uses the concrete type, otherwise ditch.
				return fmt.Errorf("error marshaling JSON of %T: %w",
					output.Output, err)
			}
			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
				Role:       "tool",
				ToolCallID: output.ToolCallID,
				Content:    string(b),
			})
		}

		// Now run it again and see what we get!
		fmt.Printf("* Sending tool response: %d\n", len(outputs))
		res, err = ExecuteChat(req, cfg.Stream)
		if err != nil {
			return fmt.Errorf("error executing chat completion request: %w", err)
		}
	}

	// b, err := json.MarshalIndent(res, "", "  ")
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(b))

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
						fmt.Printf("\n* Tool call: %s ", toolCallDelta.Function.Name)
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
