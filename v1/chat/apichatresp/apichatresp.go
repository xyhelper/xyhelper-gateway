package apichatresp

type ChatCompletion struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int64           `json:"created"`
	Model   string          `json:"model"`
	Usage   Usage           `json:"usage"`
	Choices []ChoiceMessage `json:"choices"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChoiceMessage struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
