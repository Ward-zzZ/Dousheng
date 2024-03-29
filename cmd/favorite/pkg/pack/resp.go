package pack

import (
	"errors"

	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/FavoriteServer"
	"tiktok-demo/shared/kitex_gen/VideoServer"
)

func favoriteActionResp(err errno.ErrNo) *FavoriteServer.DouyinFavoriteActionResponse {
	resp := new(FavoriteServer.DouyinFavoriteActionResponse)
	resp.BaseResp = &FavoriteServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	return resp
}

func BuildfavoriteActionResp(err error) *FavoriteServer.DouyinFavoriteActionResponse {
	if err == nil {
		return favoriteActionResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return favoriteActionResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return favoriteActionResp(s)
}

// 获取点赞视频列表的pack
func getFavoriteListResp(err errno.ErrNo, videos []*FavoriteServer.Video) *FavoriteServer.DouyinFavoriteListResponse {
	resp := new(FavoriteServer.DouyinFavoriteListResponse)
	resp.BaseResp = &FavoriteServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.VideoList = videos

	return resp
}

func BuildgetFavoriteListResp(err error, videos []*FavoriteServer.Video) *FavoriteServer.DouyinFavoriteListResponse {
	if err == nil {
		return getFavoriteListResp(errno.Success, videos)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return getFavoriteListResp(e, nil)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return getFavoriteListResp(s, nil)
}

func queryUserFavoriteVideoResp(err errno.ErrNo, isFav bool) *FavoriteServer.DouyinQueryFavoriteResponse {
	resp := new(FavoriteServer.DouyinQueryFavoriteResponse)
	resp.BaseResp = &FavoriteServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.Favorite = isFav
	return resp
}

func BuildQueryUserFavoriteVideoResp(err error, isFav bool) *FavoriteServer.DouyinQueryFavoriteResponse {
	if err == nil {
		return queryUserFavoriteVideoResp(errno.Success, isFav)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return queryUserFavoriteVideoResp(e, isFav)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return queryUserFavoriteVideoResp(s, isFav)
}

func favoriteVideoQueryResp(err errno.ErrNo, count int64) *FavoriteServer.DouyinVideoBeFavoriteResponse {
	resp := new(FavoriteServer.DouyinVideoBeFavoriteResponse)
	resp.BaseResp = &FavoriteServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.FavoriteCount = count
	return resp
}

func BuildfavoriteVideoQueryResp(err error, count int64) *FavoriteServer.DouyinVideoBeFavoriteResponse {
	if err == nil {
		return favoriteVideoQueryResp(errno.Success, count)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return favoriteVideoQueryResp(e, count)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return favoriteVideoQueryResp(s, count)
}

func ConvertVideos(videos []*VideoServer.Video) []*FavoriteServer.Video {
	videosList := make([]*FavoriteServer.Video, 0)
	for _, v := range videos {
		videosList = append(videosList, &FavoriteServer.Video{
			Id: v.Id,
			Author: &FavoriteServer.User{
				Id:            v.Author.Id,
				Name:          v.Author.Name,
				FollowCount:   v.Author.FollowCount,
				FollowerCount: v.Author.FollowerCount,
				IsFollow:      v.Author.IsFollow,
				//TotalBeFavorite: ,todo
				//FavoriteCount: ,  todo
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}
	return videosList
}
