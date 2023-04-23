package apichatrespstream

type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Delta        map[string]interface{} `json:"delta"`
	Index        int                    `json:"index"`
	FinishReason string                 `json:"finish_reason"`
}
