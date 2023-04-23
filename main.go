package main

import (
	"xyhelper-gateway/config"
	"xyhelper-gateway/v1/chat"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/ping", func(r *ghttp.Request) {
		r.Response.WriteJsonExit(g.Map{
			"message": "pong",
		})
	})
	v1group := s.Group("/v1")
	v1group.POST("/chat/completions", chat.ChatCompletions)
	s.SetPort(config.PORT)
	s.Run()
}
