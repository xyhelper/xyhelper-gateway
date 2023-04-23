package config

import (
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	PORT      int
	Tokens    = garray.NewStrArrayFrom([]string{"xyhelper-gateway-default-token"})
	TimeOutMs = 300000
	BaseURI   = "https://freechat.xyhelper.cn"
)

func init() {
	ctx := gctx.GetInitCtx()
	port := g.Cfg().MustGetWithEnv(ctx, "PORT")
	if port.IsEmpty() {
		PORT = 8080
	} else {
		PORT = port.Int()
	}
	tokens := g.Cfg().MustGetWithEnv(ctx, "TOKENS")
	if !tokens.IsEmpty() {
		// 将tokens按,分割
		tokensArray := gstr.SplitAndTrim(tokens.String(), ",")
		Tokens = garray.NewStrArrayFrom(tokensArray)
	}
	timeOutMs := g.Cfg().MustGetWithEnv(ctx, "TIMEOUTMS")
	if !timeOutMs.IsEmpty() {
		TimeOutMs = timeOutMs.Int()
	}
	baseURI := g.Cfg().MustGetWithEnv(ctx, "BASEURI")
	if !baseURI.IsEmpty() {
		BaseURI = baseURI.String()
	}

}
