// Code generated by hertz generator.

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"tiktok-demo/cmd/api/config"
	"tiktok-demo/cmd/api/initialize"
	"tiktok-demo/cmd/api/initialize/rpc"
	"tiktok-demo/shared/consts"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	// "github.com/hertz-contrib/cors"
	cfg "github.com/hertz-contrib/http2/config"
	"github.com/hertz-contrib/http2/factory"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"
	"github.com/hertz-contrib/pprof"
)

func main() {
	// initialize
	initialize.InitLogger()
	initialize.InitConfig()
	r, info := initialize.InitRegistry()
	initialize.InitSentinel()
	initialize.InitJwt()
	tracer, trcCfg := hertztracing.NewServerTracer()
	// tlsCfg := initialize.InitTLS()
	// corsCfg := initialize.InitCors()
	rpc.Init()
	// create a new server
	h := server.New(
		tracer,
		// server.WithALPN(true),
		// server.WithTLS(tlsCfg),
		server.WithHostPorts(fmt.Sprintf(":%d", config.GlobalServerConfig.Port)),
		server.WithRegistry(r, info),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(1024*1024*1024), // 1G上传限制
	)
	// add http2 protocol
	h.AddProtocol("h2", factory.NewServerFactory(
		cfg.WithReadTimeout(time.Minute),
		cfg.WithDisableKeepAlive(false)))
	// tlsCfg.NextProtos = append(tlsCfg.NextProtos, "h2")
	// use pprof&&cors&&tracing&&sentinel
	pprof.Register(h)
	// h.Use(cors.New(corsCfg))
	h.Use(hertztracing.ServerMiddleware(trcCfg))
	h.Use(hertzSentinel.SentinelServerMiddleware(
		hertzSentinel.WithServerResourceExtractor(func(c context.Context, ctx *app.RequestContext) string {
			return consts.Tiktok
		}),
		// abort with status 429 by default
		hertzSentinel.WithServerBlockFallback(func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(http.StatusTooManyRequests, nil)
			ctx.Abort()
		}),
	))
	register(h)
	h.Spin()
}
