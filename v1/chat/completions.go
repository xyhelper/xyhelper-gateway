package chat

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"xyhelper-gateway/config"
	"xyhelper-gateway/v1/chat/apichatresp"
	"xyhelper-gateway/v1/chat/apichatrespstream"

	"github.com/gogf/gf/v2/encoding/gjson"
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
	newMessage := "下面是我们的对话记录,请继续回答问题,请只给出答案内容\n"

	for _, message := range req.Messages {
		newMessage += message.Role + ":" + message.Content + "\n"
	}
	g.Log().Debug(ctx, "newMessage", newMessage)

	// AccessToken, ok := config.Tokens.Rand()

	// if !ok {
	// 	r.Response.WriteStatusExit(500)
	// }

	AccessToken := authkey
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
		var message string
		stream, err := cli.GetChatStream(newMessage)
		for err != nil {
			if err.Error() == "send message failed: 202 Accepted" {
				g.Log().Debug(ctx, "共享池新会话分配到未登录账号，重新获取会话", req)
				stream, err = cli.GetChatStream(newMessage)
				continue
			} else {
				g.Log().Error(ctx, err)
				r.Response.WriteStatusExit(500)
			}
		}

		rw := r.Response.RawWriter()
		flusher, ok := rw.(http.Flusher)
		if !ok {
			g.Log().Error(ctx, "rw.(http.Flusher) error")
			r.Response.WriteStatusExit(500)
			return
		}
		r.Response.Header().Set("Content-Type", "text/event-stream")
		r.Response.Header().Set("Cache-Control", "no-cache")
		r.Response.Header().Set("Connection", "keep-alive")
		var resData *apichatrespstream.ChatCompletion

		for text := range stream.Stream {
			if text.Role != "assistant" {
				continue
			}
			// g.Log().Debug(ctx, "message", message)

			answer := strings.Replace(text.Content, message, "", 1)
			message = text.Content

			choice := &apichatrespstream.Choice{
				Delta: map[string]interface{}{
					"content": answer,
				},
				Index: 0,
			}

			resData = &apichatrespstream.ChatCompletion{
				ID:      text.MessageID,
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "gpt-3.5-turbo-0301",
				Choices: []apichatrespstream.Choice{
					*choice,
				},
			}
			// g.Log().Debug(ctx, "resData", resData)
			resJson := gjson.New(resData)
			// g.Log().Debug(ctx, "resJson", resJson.MustToJsonString())

			_, err = fmt.Fprintf(rw, "data: %s\n\n", resJson.MustToJsonString())

			if err != nil {
				g.Log().Error(ctx, "fmt.Fprintf error", err)
				break
			}
			flusher.Flush()

		}
		_, err = fmt.Fprintf(rw, "data: %s\n\n", "[DONE]")
		if err != nil {
			g.Log().Error(ctx, "fmt.Fprintf error", err)
		}
		flusher.Flush()
	} else {
		g.Log().Debug(ctx, "Not Stream")
		data, err := cli.GetChatText(newMessage)
		for err != nil {
			if err.Error() == "send message failed: 202 Accepted" {
				g.Log().Debug(ctx, "共享池新会话分配到未登录账号，重新获取会话", req)
				data, err = cli.GetChatText(newMessage)
				continue
			} else {
				g.Log().Error(ctx, err)
				r.Response.WriteStatusExit(500)
			}
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
		g.Log().Debug(ctx, "resData", resData)
		r.Response.WriteJsonExit(resData)
	}
}
