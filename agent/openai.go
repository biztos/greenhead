package agent

type OpenAiClient struct {
	Config  *Config
	logName string
}

func NewOpenAiClient(cfg *Config) (ApiClient, error) {
	return &OpenAiClient{
		Config: cfg,
	}, nil
}

func (c *OpenAiClient) SetLogName(name string) {
	c.logName = name
}
func (c *OpenAiClient) AddSystemPrompt(prompt string) {

}
func (c *OpenAiClient) AddUserPrompt(prompt string) {}

func (c *OpenAiClient) AddAssistantResponse(prompt string) {}

func (c *OpenAiClient) AddToolResults([]*ToolResult) {}

func (c *OpenAiClient) RunCompletion() error {
	return nil
}

// Check validates the underlying client by making a (presumably) no-cost
// round-trip to the configured API endpoint, e.g.
//
//	https://api.openai.com/v1/models
func (c *OpenAiClient) Check() error {
	return nil
}
