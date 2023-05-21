package main

import (
	"context"
	VideoServer "tiktok-demo/shared/kitex_gen/VideoServer"
)

// VideoSrvImpl implements the last service interface defined in the IDL.
type VideoSrvImpl struct{}

// Feed implements the VideoSrvImpl interface.
func (s *VideoSrvImpl) Feed(ctx context.Context, req *VideoServer.DouyinFeedRequest) (resp *VideoServer.DouyinFeedResponse, err error) {
	// TODO: Your code here...
	return
}

// PublishAction implements the VideoSrvImpl interface.
func (s *VideoSrvImpl) PublishAction(ctx context.Context, req *VideoServer.DouyinPublishActionRequest) (resp *VideoServer.DouyinPublishActionResponse, err error) {
	// TODO: Your code here...
	return
}

// PublishList implements the VideoSrvImpl interface.
func (s *VideoSrvImpl) PublishList(ctx context.Context, req *VideoServer.DouyinPublishListRequest) (resp *VideoServer.DouyinPublishListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetVideoListByVideoId implements the VideoSrvImpl interface.
func (s *VideoSrvImpl) GetVideoListByVideoId(ctx context.Context, req *VideoServer.DouyinVideoListByVideoId) (resp *VideoServer.DouyinPublishListResponse, err error) {
	// TODO: Your code here...
	return
}
