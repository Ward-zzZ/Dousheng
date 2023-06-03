package errno

import (
	"errors"
	"fmt"
)

type ErrNo struct {
	ErrCode int32
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

// create a new ErrNo
func NewErrNo(code int32, msg string) ErrNo {
	return ErrNo{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

// modify error message
func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}
	//s := ServiceErr
	//s.ErrMsg = err.Error()
	return Err
}

var (
	Success    = NewErrNo(SuccessCode, "Success")
	NoRoute    = NewErrNo(NoRouteCode, "No route")
	NoMethod   = NewErrNo(NoMethodCode, "No method")
	BadRequest = NewErrNo(BadRequestCode, "Bad request")
	ServiceErr = NewErrNo(ServiceErrCode, "Service is unable to start successfully")
	ParamErr   = NewErrNo(ParamErrCode, "Wrong Parameter has been given")
	FuncErr    = NewErrNo(FuncErrCode, "Error!")
	MysqlErr   = NewErrNo(MysqlErrCode, "Mysql error!")
	// user
	UserAlreadyExistErr      = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	UserNotExistErr          = NewErrNo(UserNotExistErrCode, "User not exists")
	AuthorizationFailedErr   = NewErrNo(AuthorizationFailedErrCode, "Authorization failed")
	StructConvertFailedErr   = NewErrNo(StructConvertFailedErrCode, "Struct Convert failed")
	ChangeUserFollowCountErr = NewErrNo(ChangeUserFollowCountErrCode, "Failed to modify the follow count")
	RelationRPCErr           = NewErrNo(RelationRPCErrCode, "Failed to use relation RPC")
	FindUserErr              = NewErrNo(FindUserErrCode, "Failed to use relation RPC")
	//follow
	FollowActionErr    = NewErrNo(FollowActionErrCode, "Follow action failed")
	ActionTypeErr      = NewErrNo(ActionTypeErrCode, "Wrong action-type has been given")
	QueryFollowErr     = NewErrNo(QueryFollowErrCode, "Query relation failed")
	UserRPCErr         = NewErrNo(UserRPCErrCode, "Failed to use user RPC")
	GetFollowListErr   = NewErrNo(GetFollowListErrCode, "Failed to get follow list")
	GetFollowerListErr = NewErrNo(GetFollowerListErrCode, "Failed to get follower list")
	GetFollowSetErr    = NewErrNo(GetFollowSetErrCode, "Failed to get follow set")
	//video
	PublishActionErr         = NewErrNo(PublishActionErrCode, "Publish Action failed")
	PublishListErr           = NewErrNo(PublishListErrCode, "Publish List failed")
	FeedErr                  = NewErrNo(FeedErrCode, "Feed videos failed")
	VideoRpcUserErr          = NewErrNo(VideoRpcUserErrCode, "Video rpc User failed")
	VideoRpcRelationErr      = NewErrNo(VideoRpcRelationErrCode, "Video rpc relation failed")
	VideoListNotFound        = NewErrNo(VideoListNotFoundErrCode, "Video List is empty")
	GetVideoListByVideoIdErr = NewErrNo(GetVideoListByVideoIdErrCode, "Get Video List By Video Id Err")
	//favorite
	FavoriteVideoListNotExistErr = NewErrNo(FavoriteVideoListNotExistErrCode, "Favorite not exist")
	FavoriteActionErr            = NewErrNo(FavoriteActionErrCode, "FavoriteAction failed")
	FavoriteActionTypeErr        = NewErrNo(FavoriteActionTypeErrCode, "FavoriteActionType is wrong")
	FavoriteVideoListErr         = NewErrNo(FavoriteVideoListErrCode, "FavoriteVideoListErrCode rpc List err")
	QueryUserLikeVideoErr        = NewErrNo(FavoriteQueryUserLikeVideoErrCode, "FavoriteQueryUserLikeVideoErr rpc err")
	//comment
	CommentActionErr  = NewErrNo(CommentActionErrCode, "Comment action failed")
	GetCommentListErr = NewErrNo(GetCommentListErrCode, "Failed to get comment list")
)