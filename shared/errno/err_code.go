package errno

// Basic error code
const (
	SuccessCode int32 = iota + 0
	NoRouteCode
	NoMethodCode
	BadRequestCode
	ServiceErrCode
	ParamErrCode
	FuncErrCode
	MysqlErrCode
	InternalServerErrCode
	UserNotLoginErrCode
)

// user module error code
const (
	UserAlreadyExistErrCode int32 = iota + 10001
	UserNotExistErrCode
	AuthorizationFailedErrCode
	StructConvertFailedErrCode
	ChangeUserFollowCountErrCode
	RelationRPCErrCode
	FindUserErrCode
)

// follow module error code
const (
	FollowActionErrCode int32 = iota + 10101
	ActionTypeErrCode
	QueryFollowErrCode
	UserRPCErrCode
	GetFollowListErrCode
	GetFollowerListErrCode
	GetFollowSetErrCode
)

// comment module error code
const (
	CommentActionErrCode int32 = iota + 10201
	CommentActionTypeErrCode
	CommentAddErrCode
	CommentDelErrCode
	GetCommentListErrCode
)

// video module error code
const (
	PublishActionErrCode int32 = iota + 10301
	VideoNotExistErrCode
	VideoTypeErrCode
	PublishListErrCode
	FeedErrCode
	VideoRpcUserErrCode
	VideoRpcRelationErrCode
	VideoListNotFoundErrCode
	GetVideoListByVideoIdErrCode
	VideoContentTypeErrCode
)

// favorite module error code
const (
	FavoriteActionTypeErrCode int32 = iota + 10401
	FavoriteVideoListNotExistErrCode
	FavoriteActionErrCode
	FavoriteVideoListErrCode
	FavoriteQueryUserLikeVideoErrCode
)
