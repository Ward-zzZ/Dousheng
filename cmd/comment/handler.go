package main

import (
	"context"
	CommentServer "tiktok-demo/shared/kitex_gen/CommentServer"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *CommentServer.DouyinCommentActionRequest) (resp *CommentServer.DouyinCommentActionResponse, err error) {
	// TODO: Your code here...
	return
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *CommentServer.DouyinCommentListRequest) (resp *CommentServer.DouyinCommentListResponse, err error) {
	// TODO: Your code here...
	return
}
