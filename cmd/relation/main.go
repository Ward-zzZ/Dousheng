package main

import (
	"context"
	"net"
	"strconv"

	"tiktok-demo/cmd/relation/config"
	"tiktok-demo/cmd/relation/initialize"
	"tiktok-demo/cmd/relation/pkg/mq"
	"tiktok-demo/cmd/relation/pkg/mysql"
	"tiktok-demo/cmd/relation/pkg/redis"
	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/kitex_gen/RelationServer/relationservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func main() {
	//initialization
	initialize.InitLogger()
	initialize.InitConfig()
	IP, Port := initialize.InitFlag()
	r, info := initialize.InitRegistry(Port)
	db := initialize.InitDB()
	redisClient := initialize.InitRedis()
	MysqlManager := mysql.NewManager(db, config.GlobalServerConfig.MysqlInfo.Salt)
	mq.InitMq()
	mq.MysqlManager = MysqlManager
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	userClient := initialize.InitUser()
	svr := relationservice.NewServer(&RelationServiceImpl{
		MysqlManager: MysqlManager,
		RedisManager: redis.NewManager(redisClient),
		UserManager:  userClient,
	}, server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "relation_srv",
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
