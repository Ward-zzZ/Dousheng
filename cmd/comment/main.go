package main

import (
	"context"
	"net"
	"strconv"

	"tiktok-demo/cmd/comment/config"
	"tiktok-demo/cmd/comment/initialize"
	// "tiktok-demo/cmd/comment/pkg/mq"
	"tiktok-demo/cmd/comment/pkg/mysql"
	"tiktok-demo/cmd/comment/pkg/redis"
	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/kitex_gen/CommentServer/commentservice"

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
	MysqlManager := mysql.NewManager(db)
	// mq.InitMq()
	// mq.MysqlManager = MysqlManager
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	userClient := initialize.InitUser()
	svr := commentservice.NewServer(&CommentServiceImpl{
		MysqlManager: MysqlManager,
		RedisManager: redis.NewManager(redisClient),
		UserManager:  userClient,
	}, server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "comment_srv",
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
