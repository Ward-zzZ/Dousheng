package config

import (
	"tiktok-demo/shared/kitex_gen/CommentServer/commentservice"
	"tiktok-demo/shared/kitex_gen/FavoriteServer/favoriteservice"
	"tiktok-demo/shared/kitex_gen/RelationServer/relationservice"
	"tiktok-demo/shared/kitex_gen/UserServer/userservice"
)

var (
	GlobalServerConfig ServerConfig
	GlobalConsulConfig ConsulConfig

	RelationClient relationservice.Client
	UserClient     userservice.Client
	CommentClient  commentservice.Client
	FavoriteClient favoriteservice.Client
)
