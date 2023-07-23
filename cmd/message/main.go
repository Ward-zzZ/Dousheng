package main

import (
	"context"
	"net"
	"strconv"

	"tiktok-demo/cmd/message/config"
	"tiktok-demo/cmd/message/initialize"
	"tiktok-demo/cmd/message/pkg/mysql"
	// "tiktok-demo/cmd/message/pkg/redis"
	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/kitex_gen/MessageServer/messageservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func main() {

	initialize.InitLogger()
	initialize.InitConfig()
	IP, Port := initialize.InitFlag()
	r, info := initialize.InitRegistry(Port)
	db := initialize.InitDB()
	// redisClient := initialize.InitRedis()
	MysqlManager := mysql.NewManager(db)

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	svr := messageservice.NewServer(&MessageServiceImpl{
		MysqlManager: MysqlManager,
		// RedisManager: redis.NewManager(redisClient),
	}, server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "message_srv",
	}),
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}))
	err := svr.Run()

	if err != nil {
		klog.Fatalf("run server failed: %v\n", err)
	}

}
