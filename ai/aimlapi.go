package ai

type AIMLAPIRequest struct {
	Model               string                   `json:"model"`
	Messages            []AIMLAPIRequestMessages `json:"messages"`
	MaxCompletionTokens int                      `json:"max_completion_tokens"`
	MaxTokens           int                      `json:"max_tokens"`
	Stream              bool                     `json:"stream"`
	StreamOptions       `json:"stream_options"`
	Temperature         int `json:"temperature"`
	TopP                int `json:"top_p"`
	Seed                int `json:"seed"`
	MinP                int `json:"min_p"`
	TopK                int `json:"top_k"`
	RepetitionPenalty   int `json:"repetition_penalty"`
	TopA                int `json:"top_a"`
	FrequencyPenalty    int `json:"frequency_penalty"`
	Prediction          `json:"prediction"`
	PresencePenalty     int    `json:"presence_penalty"`
	Tools               []Tool `json:"tools"`
	ToolChoice          string `json:"tool_choice"`
	ParallelToolCalls   bool   `json:"parallel_tool_calls"`
	Stop                string `json:"stop"`
	Logprobs            bool   `json:"logprobs"`
	TopLogprobs         int    `json:"top_logprobs"`
	LogitBias           `json:"logit_bias"`
	ResponseFormat      `json:"response_format"`
}

type AIMLAPIRequestMessages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type Prediction struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Function struct {
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Parameters  interface{} `json:"parameters"`
	Strict      bool        `json:"strict"`
	Required    []string    `json:"required"`
}

type Tool struct {
	Type     string `json:"type"`
	Function `json:"function"`
}

type LogitBias struct {
	AdditionalProperty int `json:"ANY_ADDITIONAL_PROPERTY"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type AIMLAPIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Choices []Choice `json:"choices"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Choice struct {
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
	Logprobs     `json:"logprobs"`
	Message      AIMLAPIResponseMessage `json:"message"`
}

type Logprobs struct {
	Content []interface{} `json:"content"`
	Refusal []interface{} `json:"refusal"`
}

type AIMLAPIResponseMessage struct {
	Role    string      `json:"role"`
	Content string      `json:"content"`
	Refusal interface{} `json:"refusal"`
}
