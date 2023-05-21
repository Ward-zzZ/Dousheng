package main

import (
	"context"
	RelationServer "tiktok-demo/shared/kitex_gen/RelationServer"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, req *RelationServer.DouyinRelationActionRequest) (resp *RelationServer.DouyinRelationActionResponse, err error) {
	// TODO: Your code here...
	return
}

// MGetRelationFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MGetRelationFollowList(ctx context.Context, req *RelationServer.DouyinRelationFollowListRequest) (resp *RelationServer.DouyinRelationFollowListResponse, err error) {
	// TODO: Your code here...
	return
}

// MGetUserRelationFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MGetUserRelationFollowerList(ctx context.Context, req *RelationServer.DouyinRelationFollowerListRequest) (resp *RelationServer.DouyinRelationFollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryRelation implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) QueryRelation(ctx context.Context, req *RelationServer.DouyinQueryRelationRequest) (resp *RelationServer.DouyinQueryRelationResponse, err error) {
	// TODO: Your code here...
	return
}
