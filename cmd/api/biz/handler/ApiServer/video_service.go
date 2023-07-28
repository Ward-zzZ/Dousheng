// Code generated by hertz generator.

package ApiServer

import (
	"bytes"
	"context"
	"io"
	"time"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	ApiServer "tiktok-demo/cmd/api/biz/model/ApiServer"
	"tiktok-demo/cmd/api/config"
	"tiktok-demo/cmd/api/pkg"
	Globalconsts "tiktok-demo/shared/consts"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/VideoServer"
	// "tiktok-demo/shared/tools"
)

// Feed .
// @router /douyin/feed/ [GET]
func Feed(ctx context.Context, c *app.RequestContext) {
	// var err error
	// var req ApiServer.DouyinFeedRequest
	// err = c.BindAndValidate(&req)
	// if err != nil {
	// 	c.String(consts.StatusBadRequest, err.Error())
	// 	return
	// }

		// token := c.Query("token")
	latestTime := c.Query("latest_time")
	var timestamp int64 = 0
	if latestTime != "" {
		timestamp, _ = strconv.ParseInt(latestTime, 10, 64)
	} else {
		timestamp = time.Now().UnixMilli()
	}

	// var laststTime, userId int64
	// // 获取最近的时间并判断处理
	// lastst_time := c.Query("latest_time")
	// if len(lastst_time) != 0 {
	// 	if latesttime, err := strconv.Atoi(lastst_time); err != nil {
	// 		hlog.Error("strconv.Atoi(lastst_time) error", err)
	// 		pkg.SendFeedResponse(c, errno.ConvertErr(err), nil)
	// 		return
	// 	} else {
	// 		laststTime = int64(latesttime)
	// 	}
	// }

	//获取token中传来的user id
	var userId int64
	user, _ := c.Get(Globalconsts.IdentityKey)
	if user == nil {
		// default userId when not login
		userId = 0
	} else {
		userId = user.(*ApiServer.User).Id
	}

	//调用feed rpc
	videosFeedResp, _ := config.GlobalVideoClient.Feed(context.Background(), &VideoServer.DouyinFeedRequest{
		LatestTime: timestamp,
		UserId:     userId,
	})
	if videosFeedResp.BaseResp.StatusCode != errno.SuccessCode {
		respErr := errno.NewErrNo(videosFeedResp.BaseResp.StatusCode, videosFeedResp.BaseResp.StatusMsg)
		hlog.Error("feed rpc error", "error", respErr)
		pkg.SendFeedResponse(c, respErr, nil)
		return
	}
	pkg.SendFeedResponse(c, nil, videosFeedResp.VideoList)
}

// PublishAction .
// @router /douyin/publish/action/ [POST]
func PublishAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ApiServer.DouyinPublishActionRequest

	req.Title = c.PostForm("title")
	req.Token = c.PostForm("token")
	//接收视频文件并处理。
	fileHeader, err := c.Request.FormFile("data")
	if err != nil {
		pkg.SendPublishActionResponse(c, errno.ConvertErr(err), nil)
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		hlog.Error("fileHeader.Open() error", err)
		pkg.SendPublishActionResponse(c, errno.ConvertErr(err), nil)
		return
	}
	defer file.Close()

	//copy为[]byte格式
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		hlog.Error("io.Copy(buf, file) error", err)
		pkg.SendPublishActionResponse(c, errno.ConvertErr(err), nil)
		return
	}

	//获取token中传来的user id
	var userId int64
	user, _ := c.Get(Globalconsts.IdentityKey)
	if user == nil {
		hlog.Error("user is not login")
		pkg.SendPublishActionResponse(c, errno.UserNotLoginErr, nil)
		return
	} else {
		userId = user.(*ApiServer.User).Id
	}

	PublishActionResp, _ := config.GlobalVideoClient.PublishAction(context.Background(), &VideoServer.DouyinPublishActionRequest{
		UserId: userId,
		Title:  req.Title,
		Token:  req.Token,
		Data:   buf.Bytes(),
	})

	if PublishActionResp.BaseResp.StatusCode != errno.SuccessCode {
		respErr := errno.NewErrNo(PublishActionResp.BaseResp.StatusCode, PublishActionResp.BaseResp.StatusMsg)
		hlog.Error("PublishActionResp error", "error", respErr)
		pkg.SendPublishActionResponse(c, respErr, nil)
		return

	}
	pkg.SendPublishActionResponse(c, errno.Success, nil)
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

	//获取token中传来的user id
	var userId int64
	user, _ := c.Get(Globalconsts.IdentityKey)
	if user == nil {
		// default userId when not login
		userId = 0
	} else {
		userId = user.(*ApiServer.User).Id
	}
	PublishListResp, _ := config.GlobalVideoClient.PublishList(context.Background(), &VideoServer.DouyinPublishListRequest{
		UserId: userId,
	})
	if PublishListResp.BaseResp.StatusCode != errno.SuccessCode {
		respErr := errno.NewErrNo(PublishListResp.BaseResp.StatusCode, PublishListResp.BaseResp.StatusMsg)
		hlog.Error("PublishList rpc error", "error", respErr)
		pkg.SendPublishListResponse(c, respErr, nil)
		return
	}
	pkg.SendPublishListResponse(c, errno.Success, PublishListResp.VideoList)
}
