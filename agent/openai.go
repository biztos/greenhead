// openai.go
//
// TODO: hooks to preprocess/postprocess stuff.
// ...want option to prepocess the whole outgoing request so you can manage context
//
//	this would include tool results b/c agent adds them
//
// ...want option to postprocess the result before the agent gets it
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

type OpenAiClient struct {
	Client          *openai.Client
	Model           string
	Tools           []openai.Tool
	Stream          bool
	StreamToolCalls bool
	ContextMessages []openai.ChatCompletionMessage

	streamPrint func(a ...any)
	preFunc     func(any) error
	postFunc    func(any) error
	logger      *slog.Logger
}

// NewOpenAiClient returns a client initialized according to cfg.
func NewOpenAiClient(cfg *Config) (ApiClient, error) {

	token := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(token)
	model := cfg.Model
	if model == "" {
		model = OpenAiClientDefaultModel
	}
	tools := make([]openai.Tool, len(cfg.Tools))
	for idx, name := range cfg.Tools {
		// This is redundant to the Agent checks, but you aren't guaranteed
		// to be instantiating from an Agent, so here we go, it's cheap.
		t := registry.Get(name)
		if t == nil {
			return nil, fmt.Errorf("tool not registered: %s", name)
		}
		tools[idx] = t.OpenAiTool()
	}
	sp_func, err := ColorPrintFunc(cfg.Color, cfg.BgColor)
	if err != nil {
		return nil, fmt.Errorf("streaming colors: %w", err)
	}
	msgs := make([]openai.ChatCompletionMessage, 0, len(cfg.Context))
	for _, item := range cfg.Context {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    item.Role,
			Content: item.Content,
		})
	}
	return &OpenAiClient{
		Client:          client,
		Model:           model,
		Stream:          cfg.Stream,
		ContextMessages: msgs,
		Tools:           tools,

		streamPrint: sp_func,
		logger:      slog.Default(),
	}, nil
}

// SetLogger implements ApiClient.
func (c *OpenAiClient) SetLogger(logger *slog.Logger) {
	c.logger = logger
}

// SetPreFunc implements ApiClient.
func (c *OpenAiClient) SetPreFunc(f func(any) error) {
	c.preFunc = f
}

// SetPostFunc implements ApiClient.
func (c *OpenAiClient) SetPostFunc(f func(any) error) {
	c.postFunc = f
}

// AddContext implements ApiClient.
func (c *OpenAiClient) AddContext(item *ContextItem) {
	c.ContextMessages = append(c.ContextMessages, openai.ChatCompletionMessage{
		Role:    item.Role,
		Content: item.Content,
	})
}

// RunCompletion implements ApiClient.
func (c *OpenAiClient) RunCompletion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {

	// Get our message array which we will need for the request and context.
	// We will add one response message below.
	add_len := 1
	if len(req.ToolResults) > 0 {
		add_len = len(req.ToolResults)
	}
	msgs := make([]openai.ChatCompletionMessage, 0, len(c.ContextMessages)+add_len+1)
	msgs = append(msgs, c.ContextMessages...)
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

	// Create openai-specific request.
	oai_req := openai.ChatCompletionRequest{
		Model:    c.Model,
		Tools:    c.Tools,
		Messages: msgs,
		Stream:   c.Stream,
		// TODO: MaxCompletionTokens support
		// TODO: temperature and so on -- want a custom config I guess.
		// ...and so on...

	}

	if c.preFunc != nil {
		err := c.preFunc(&oai_req)
		if err != nil {
			return nil, fmt.Errorf("error from preprocessing function: %w", err)
		}
	}

	// Run that -- we want to keep an eye on durations.
	//
	// TODO: overall instrumentation plan, this ain't it obviously, but it's
	// a start and we could stick it in Grafana.
	start_ts := time.Now()
	c.logger.Info("creating chat completion", "model", c.Model, "stream", c.Stream)
	res, err := c.CreateChatCompletion(ctx, oai_req)
	c.logger.Info("created chat completion", utils.DurLog(start_ts)...)
	if err != nil {
		return nil, fmt.Errorf("error creating chat completion: %w", err)
	}
	// Post-process the result before dealing with it here, if applicable.
	if c.postFunc != nil {
		err := c.postFunc(&res)
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
	c.ContextMessages = msgs

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

// Check implements ApiClient by making a call to the models endpoint, which
// requires authentication.
func (c *OpenAiClient) Check(ctx context.Context) error {
	c.logger.Info("checking API with ListModels")
	start_ts := time.Now()
	models, err := c.Client.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("error running check with ListModels: %w", err)
	}
	c.logger.Info("check successful",
		utils.DurLog(start_ts)...)

	// Not all that useful if we don't treat it as an error, but...
	if len(models.Models) == 0 {
		c.logger.Warn("no models found in check!")
	}

	return nil
}

// ExecuteCompletion executes a chat completion for both streaming and
// non-streaming cases.
func (c *OpenAiClient) CreateChatCompletion(ctx context.Context, r openai.ChatCompletionRequest) (res openai.ChatCompletionResponse, err error) {

	if !c.Stream {
		return c.Client.CreateChatCompletion(ctx, r)
	}

	// Streaming path is trickier; watch for Miss Steaks!
	r.Stream = true
	r.StreamOptions = &openai.StreamOptions{IncludeUsage: true}

	stream, serr := c.Client.CreateChatCompletionStream(ctx, r)
	if serr != nil {
		err = serr
		return
	}
	defer stream.Close()

	// Build up the response as we receive chunks
	// var res = openai.ChatCompletionResponse{}
	var contentBuilder strings.Builder
	var finishReason string
	var role string
	var toolCalls []openai.ToolCall
	var got_usage bool
	for {
		response, serr := stream.Recv()
		if errors.Is(serr, io.EOF) {
			break
		}
		if err != nil {
			err = serr
			return
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
			err = fmt.Errorf("chunk after usage: %s", utils.MustJsonString(response))
			return
		}
		if len(response.Choices) == 0 {
			// This is the usage chunk; there should be only one.
			if response.Usage == nil {
				err = fmt.Errorf("unexpected nil Usage: %v", response)
				return
			}
			if res.Usage.TotalTokens > 0 {
				err = fmt.Errorf("dupe usage chunk: %v", response.Usage)
				return
			}

			res.Usage = *response.Usage
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
			c.streamPrint(content) // Print content as it arrives
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
							c.streamPrint(frag)
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
							c.streamPrint(toolCallDelta.Function.Arguments)
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
	fmt.Println()

	return res, nil
}
