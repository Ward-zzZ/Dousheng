package config

import (
	"tiktok-demo/shared/kitex_gen/CommentServer/commentservice"
	"tiktok-demo/shared/kitex_gen/FavoriteServer/favoriteservice"
	"tiktok-demo/shared/kitex_gen/RelationServer/relationservice"
	"tiktok-demo/shared/kitex_gen/UserServer/userservice"
	"tiktok-demo/shared/kitex_gen/VideoServer/videosrv"
)

var (
	GlobalServerConfig ServerConfig
	GlobalConsulConfig ConsulConfig

	GlobalCommentClient  commentservice.Client
	GlobalFavoriteClient favoriteservice.Client
	GlobalRelationClient relationservice.Client
	GlobalUserClient     userservice.Client
	GlobalVideoClient    videosrv.Client
)
