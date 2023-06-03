package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/limiter"
	"tiktok-demo/shared/middleware"
)

// CommonMW is a middleware that is used for all routers.
func CommonMW() []app.HandlerFunc {
	return []app.HandlerFunc{
		// use recovery mw
		middleware.Recovery(),
		// use gzip mw
		gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".jpg", ".mp4", ".png"})),
		// use limiter mw
		limiter.AdaptiveLimit(limiter.WithCPUThreshold(900)),
	}
}
