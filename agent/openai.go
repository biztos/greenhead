// agent/openai.go

package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/utils"
)

var OpenAiClientDefaultModel = openai.GPT4o

// OpenAiClient is an ApiClient that builds on BasicApiClient.
type OpenAiClient struct {
	BasicApiClient
	Client  *openai.Client
	History []openai.ChatCompletionMessage
}

// NewOpenAiClient returns a client initialized for the OpenAI API.
//
// The environment variable OPENAI_API_KEY must be set to a valid key.
func NewOpenAiClient() (ApiClient, error) {

	token := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(token)
	return &OpenAiClient{
		BasicApiClient{
			Client: client,
			Logger: slog.Default(),
		},
		client,
		nil,
	}, nil

}

// ClearContext implements ApiClient by clearing the initial context and also
// the message history.
func (c *OpenAiClient) ClearContext() {
	c.ContextItems = nil
	c.History = nil
}

// Check implements ApiClient by querying the model list.
func (c *OpenAiClient) Check(ctx context.Context) error {

	c.Logger.Info("checking API with ListModels")
	start_ts := time.Now()
	models, err := c.Client.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("error running check with ListModels: %w", err)
	} else if len(models.Models) == 0 {
		return fmt.Errorf("no models found")
	}
	c.Logger.Info("check successful", utils.DurLog(start_ts)...)

	return nil
}

// RunCompletion implements ApiClient by running the OpenAI chat completion.
func (c *OpenAiClient) RunCompletion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {

	// Create the context we will send, in the native format.
	msgs := make([]openai.ChatCompletionMessage, 0,
		len(c.ContextItems)+len(c.History)+len(req.ToolResults)+1)
	for _, item := range c.ContextItems {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    item.Role,
			Content: item.Content,
		})
	}
	msgs = append(msgs, c.History...)

	// Add tool results if applicable.
	if len(req.ToolResults) > 0 {
		for _, tr := range req.ToolResults {
			b, err := json.Marshal(tr.Output)
			if err != nil {
				// TODO: prove this %T uses the concrete type, otherwise ditch.
				return nil, fmt.Errorf("error marshaling JSON of %T: %w",
					tr.Output, err)
			}
			msgs = append(msgs, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				ToolCallID: tr.Id,
				Content:    string(b),
			})
		}
	} else {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Content,
		})
	}

	// Get the tools in openai format.
	tools := make([]openai.Tool, 0, len(c.Tools))
	for _, name := range c.Tools {
		t := registry.Get(name)
		if t == nil {
			return nil, fmt.Errorf("tool not registered: %s", name)
		}
		tools = append(tools, t.OpenAiTool())
	}

	// Create openai-specific request.
	oai_req := openai.ChatCompletionRequest{
		Model:               c.Model,
		Tools:               tools,
		Messages:            msgs,
		Stream:              c.Streaming,
		MaxCompletionTokens: c.MaxCompletionTokens,
		// TODO: temperature and so on -- want a custom config I guess.
		// ...and so on...

	}
	if c.PreFunc != nil {
		if err := c.PreFunc(c, &oai_req); err != nil {
			return nil, fmt.Errorf("error from preprocessor: %w", err)
		}
	}
	// Run that -- we want to keep an eye on durations.
	//
	// TODO: overall instrumentation plan, this ain't it obviously, but it's
	// a start and we could stick it in Grafana.
	start_ts := time.Now()
	c.Logger.Info("creating chat completion", "model", c.Model, "stream", c.Streaming)
	res, err := c.CreateChatCompletion(ctx, oai_req)
	// TODO: give it some thought, is this a good way to do the duration logging?
	// We want to say what took how long, not what the result was.  Need to be
	// able to search by the "what" and the fact of the duration.
	c.Logger.Info("creating chat completion", utils.DurLog(start_ts)...)
	if err != nil {
		return nil, fmt.Errorf("error creating chat completion: %w", err)
	}
	// Post-process the result before dealing with it here, if applicable.
	if c.PostFunc != nil {
		err := c.PostFunc(c, &res)
		if err != nil {
			return nil, fmt.Errorf("error from postprocessing function: %w", err)
		}
	}

	// Result needs to have one choice, getting multiples (or none!) is Bad.
	if len(res.Choices) != 1 {
		return nil, fmt.Errorf("wrong number of choices in response: %d",
			len(res.Choices))
	}
	msg := res.Choices[0].Message
	// TODO: consider what to usefully do with FinishReason.
	fin := string(res.Choices[0].FinishReason)
	// TODO: consider ContentFilterResults, could be sticky bastards.

	// TODO: do something else with msg.Refusal, maybe?  Need examples.
	// This is supposed to be "safety" related, per openai:
	// https://platform.openai.com/docs/guides/structured-outputs/refusals?api-mode=responses
	if msg.Refusal != "" {
		return nil, fmt.Errorf("endpoint refused to create completion: %s", msg.Refusal)
	}

	usage := &Usage{
		Input:  res.Usage.PromptTokens,
		Output: res.Usage.CompletionTokens,
		Total:  res.Usage.TotalTokens,
	}
	if details := res.Usage.PromptTokensDetails; details != nil {
		usage.CachedInput = details.CachedTokens
	}
	if details := res.Usage.CompletionTokensDetails; details != nil {
		usage.Reasoning = details.ReasoningTokens
	}

	// Update the context window now (do NOT add to context window before
	// running error-free, otherwise retry will be wrong).
	msgs = append(msgs, msg)
	c.History = append(c.History, msg)

	// Get our preferred tool-call format.
	tool_calls := make([]*ToolCall, len(msg.ToolCalls))
	for idx, tc := range msg.ToolCalls {
		// watch for anything not a function
		// TODO: figure out whether non-function tool calls get reported back
		// or not?
		if tc.Type != openai.ToolTypeFunction {
			return nil, fmt.Errorf("unexpected tool call type: %s", tc.Type)
		}
		tool_calls[idx] = &ToolCall{
			Id:   tc.ID,
			Name: tc.Function.Name,
			Args: tc.Function.Arguments,
		}
	}

	// Return in result format with raw objects.
	return &CompletionResponse{
		FinishReason: fin,
		Content:      msg.Content,
		ToolCalls:    tool_calls,
		Usage:        usage,
		RawCompletions: []*RawCompletion{
			{
				Request:  oai_req,
				Response: res,
			},
		},
	}, nil
}

// CreateChatCompletion executes a chat completion for both streaming and
// non-streaming cases.
func (c *OpenAiClient) CreateChatCompletion(ctx context.Context, r openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {

	if !c.Streaming {
		return c.Client.CreateChatCompletion(ctx, r)
	}

	// Streaming path is trickier; watch for Miss Steaks!
	var empty openai.ChatCompletionResponse
	r.Stream = true
	r.StreamOptions = &openai.StreamOptions{IncludeUsage: true}

	stream, err := c.Client.CreateChatCompletionStream(ctx, r)
	if err != nil {
		return empty, err
	}
	defer stream.Close()

	// Build up the response as we receive chunks
	var res = openai.ChatCompletionResponse{}
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
			return empty, err
		}

		// If this is the first chunk, initialize the response
		if res.ID == "" {
			res.ID = response.ID
			res.Object = response.Object // TODO: remove ".chunk"
			res.Created = response.Created
			res.Model = response.Model
			res.SystemFingerprint = response.SystemFingerprint
		}

		// Handle the Usage chunk, which by definition is the last one.
		if got_usage {
			return empty, fmt.Errorf("chunk after usage: %s", utils.MustJsonString(response))
		}
		if len(response.Choices) == 0 {
			// This is the usage chunk; there should be only one.
			if response.Usage == nil {
				return empty, fmt.Errorf("unexpected nil Usage: %v", response)
			}
			if res.Usage.TotalTokens > 0 {
				return empty, fmt.Errorf("dupe usage chunk: %v", response.Usage)
			}

			res.Usage = *response.Usage
			c.PrintFunc("\n")
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
			c.PrintFunc(content) // Print content as it arrives
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
						if c.StreamToolCalls {
							frag := fmt.Sprintf("\n* Tool call: %s ", toolCallDelta.Function.Name)
							c.PrintFunc(frag)
						}
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
						// NB: only do this if configured to, which by default, nope.
						if c.StreamToolCalls {
							c.PrintFunc(toolCallDelta.Function.Arguments)
						}
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

	res.Choices = []openai.ChatCompletionChoice{
		{
			Message:      message,
			FinishReason: openai.FinishReason(finishReason),
			Index:        0,
		},
	}

	// Add a new line after streaming is complete
	c.PrintFunc("\n")

	return res, nil
}

func init() {
	RegisterNewApiClientFunc("openai", NewOpenAiClient)
}
