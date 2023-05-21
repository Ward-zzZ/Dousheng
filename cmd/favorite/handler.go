package main

import (
	"context"
	FavoriteServer "tiktok-demo/shared/kitex_gen/FavoriteServer"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *FavoriteServer.DouyinFavoriteActionRequest) (resp *FavoriteServer.DouyinFavoriteActionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteList(ctx context.Context, req *FavoriteServer.DouyinFavoriteListRequest) (resp *FavoriteServer.DouyinFavoriteListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFavoriteUser implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteUser(ctx context.Context, req *FavoriteServer.DouyinUserBeFavoriteRequest) (resp *FavoriteServer.DouyinUserBeFavoriteResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFavoriteVideo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteVideo(ctx context.Context, req *FavoriteServer.DouyinVideoBeFavoriteRequest) (resp *FavoriteServer.DouyinVideoBeFavoriteResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryUserLikeVideo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) QueryUserLikeVideo(ctx context.Context, req *FavoriteServer.DouyinQueryFavoriteRequest) (resp *FavoriteServer.DouyinQueryFavoriteResponse, err error) {
	// TODO: Your code here...
	return
}
