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
	"tiktok-demo/shared/kitex_gen/RelationServer/relationservice"
	mw "tiktok-demo/shared/middleware"
)

func initRelation() {
	// init resolver
	r, err := consul.NewConsulResolver(fmt.Sprintf("%s:%d",
		config.GlobalConsulConfig.Host,
		config.GlobalConsulConfig.Port))
	if err != nil {
		klog.Fatalf("new consul client failed: %s", err.Error())
	}
	// init opentelemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.RelationSrvInfo.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)

	// create a new client
	c, err := relationservice.NewClient(
		config.GlobalServerConfig.RelationSrvInfo.Name,
		client.WithResolver(r),                                     // service discovery
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
		client.WithMuxConnection(1),                                // multiplexing
		client.WithMiddleware(mw.CommonMiddleware),                 // TODO: rpc info tracing,different
		client.WithInstanceMW(mw.ClientMiddleware),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.RelationSrvInfo.Name}),
	)
	if err != nil {
		klog.Fatalf("ERROR: cannot init client: %v\n", err)
	}
	config.GlobalRelationClient = c
}
