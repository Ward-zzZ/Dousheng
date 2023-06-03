package rpc

import (
	"fmt"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	consul "github.com/kitex-contrib/registry-consul"

	"tiktok-demo/cmd/api/config"
	"tiktok-demo/shared/kitex_gen/CommentServer/commentservice"
	mw "tiktok-demo/shared/middleware"
)

func initComment() {
	// init resolver
	r, err := consul.NewConsulResolver(fmt.Sprintf("%s:%d",
		config.GlobalConsulConfig.Host,
		config.GlobalConsulConfig.Port))
	if err != nil {
		klog.Fatalf("new consul client failed: %s", err.Error())
	}
	// init opentelemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.CommentSrvInfo.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)

	// create a new client
	c, err := commentservice.NewClient(
		config.GlobalServerConfig.CommentSrvInfo.Name,
		client.WithResolver(r),                                     // service discovery
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
		client.WithMuxConnection(1),                                // multiplexing
		client.WithMiddleware(mw.CommonMiddleware),                 // TODO: rpc info tracing,different
		client.WithInstanceMW(mw.ClientMiddleware),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.CommentSrvInfo.Name}),
	)
	if err != nil {
		klog.Fatalf("ERROR: cannot init client: %v\n", err)
	}
	config.GlobalCommentClient = c
}
