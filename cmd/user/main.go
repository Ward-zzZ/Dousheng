package main

import (
	"context"
	"net"
	"strconv"

	"tiktok-demo/cmd/user/config"
	"tiktok-demo/cmd/user/initialize"
	"tiktok-demo/cmd/user/pkg/md5"
	"tiktok-demo/cmd/user/pkg/mq"
	"tiktok-demo/cmd/user/pkg/mysql"
	"tiktok-demo/cmd/user/pkg/redis"

	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/kitex_gen/UserServer/userservice"

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
	MysqlManager := mysql.NewUserManager(db, config.GlobalServerConfig.MysqlInfo.Salt)
	mq.InitMq()
	mq.MysqlManager = MysqlManager
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	relationClient := initialize.InitRelation()
	srv := userservice.NewServer(&UserServiceImpl{
		MysqlManager:    MysqlManager,
		RedisManager:    redis.NewManager(redisClient),
		RealtionManager: relationClient,
		EncryptManager:  &md5.EncryptManager{Salt: config.GlobalServerConfig.MysqlInfo.Salt},
	},
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}))
	err := srv.Run()
	if err != nil {
		klog.Fatalf("run server failed: %v", err)
	}
}
