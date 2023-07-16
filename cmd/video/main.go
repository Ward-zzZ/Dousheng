package main

import (
	"context"
	"net"
	"strconv"

	"tiktok-demo/cmd/video/config"
	"tiktok-demo/cmd/video/initialize"
	"tiktok-demo/cmd/video/pkg/minio"
	"tiktok-demo/cmd/video/pkg/mysql"
	"tiktok-demo/shared/consts"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	VideoServer "tiktok-demo/shared/kitex_gen/VideoServer/videosrv"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	IP, Port := initialize.InitFlag()
	r, info := initialize.InitRegistry(Port)
	db := initialize.InitDB()
	minioClient := initialize.InitMinio()

	initialize.InitRelation()
	initialize.InitComment()
	initialize.InitUser()
	initialize.InitFavorite()

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	srv := VideoServer.NewServer(&VideoSrvImpl{
		MysqlManager:    mysql.NewManager(db),
		MinioManager:    minio.NewManager(minioClient, config.GlobalServerConfig.MinioInfo.VideoBucket, config.GlobalServerConfig.MinioInfo.CoverBucket),
		RealtionManager: config.RelationClient,
		CommentManager:  config.CommentClient,
		UserManager:     config.UserClient,
		FavoriteManager: config.FavoriteClient,
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
