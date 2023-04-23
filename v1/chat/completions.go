package chat

import (
	"strings"
	"sync"
	"time"
	"xyhelper-gateway/config"
	"xyhelper-gateway/v1/chat/apichatresp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xyhelper/chatgpt-go"
)

var (
	TokenLockMap = make(map[string]*sync.Mutex)
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionsRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}
type Author struct {
	Role string `json:"role"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

func ChatCompletions(r *ghttp.Request) {
	ctx := r.Context()
	// ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	// defer cancel()
	g.Log().Debug(ctx, "Conversation start....................")
	authkey := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

	// 如果 authkey 为空, 则返回 401
	if authkey == "" {
		r.Response.WriteStatusExit(401)
	}

	req := &ChatCompletionsRequest{}
	err := r.GetStruct(req)
	if err != nil {
		r.Response.Status = 400
		r.Response.WriteJsonExit(err.Error())
	}
	g.Log().Debug(ctx, "req", req)
	var newMessage string

	for _, message := range req.Messages {
		newMessage += message.Role + ":" + message.Content + "\n"
	}

	AccessToken, ok := config.Tokens.Rand()
	if !ok {
		r.Response.WriteStatusExit(500)
	}
	// 按token加锁
	if _, ok := TokenLockMap[AccessToken]; !ok {
		TokenLockMap[AccessToken] = &sync.Mutex{}
	}
	TokenLockMap[AccessToken].Lock()
	defer TokenLockMap[AccessToken].Unlock()
	cli := chatgpt.NewClient(
		chatgpt.WithAccessToken(AccessToken),
		chatgpt.WithTimeout(time.Duration(config.TimeOutMs*1000*1000)),
		chatgpt.WithBaseURI(config.BaseURI),
	)

	if req.Stream {
		g.Log().Debug(ctx, "Stream")
	} else {
		g.Log().Debug(ctx, "Not Stream")
		data, err := cli.GetChatText(newMessage)
		if err != nil {
			g.Log().Error(ctx, err)
			r.Response.WriteStatusExit(500)
		}
		// var resData *apichatresp.ChatCompletion
		resData := &apichatresp.ChatCompletion{
			ID:      data.MessageID,
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "gpt-3.5-turbo-0301",
			Usage: apichatresp.Usage{
				PromptTokens:     0,
				CompletionTokens: 0,
				TotalTokens:      0,
			},
			Choices: []apichatresp.ChoiceMessage{
				{
					FinishReason: "stop",
					Index:        0,
					Message: apichatresp.Message{
						Role:    "assistant",
						Content: data.Content,
					},
				},
			},
		}
		r.Response.WriteJsonExit(resData)
	}
}
