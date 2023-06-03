// Code generated by hertz generator.

package ApiServer

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
	ApiServer "tiktok-demo/cmd/api/biz/model/ApiServer"
	// "github.com/cloudwego/hertz/pkg/common/hlog"
	// mw "tiktok-demo/cmd/api/biz/middleware"
	// "tiktok-demo/cmd/api/config"
	// "tiktok-demo/cmd/api/pkg"
	// "tiktok-demo/shared/consts"
	// "tiktok-demo/shared/errno"
	// "tiktok-demo/shared/kitex_gen/CommentServer"
	// "tiktok-demo/shared/kitex_gen/FavoriteServer"
	// "tiktok-demo/shared/kitex_gen/RelationServer"
	// "tiktok-demo/shared/kitex_gen/UserServer"
	// "tiktok-demo/shared/kitex_gen/VideoServer"
	// "tiktok-demo/shared/tools"
)

// Feed .
// @router /douyin/feed/ [GET]
func Feed(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ApiServer.DouyinFeedRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(ApiServer.DouyinFeedResponse)

	c.JSON(consts.StatusOK, resp)
}

// PublishAction .
// @router /douyin/publish/action/ [POST]
func PublishAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ApiServer.DouyinPublishActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(ApiServer.DouyinPublishActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// PublishList .
// @router /douyin/publish/list/ [GET]
func PublishList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ApiServer.DouyinPublishListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(ApiServer.DouyinPublishListResponse)

	c.JSON(consts.StatusOK, resp)
}
