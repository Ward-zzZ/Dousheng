package pack

import (
	"errors"

	"tiktok-demo/cmd/message/pkg/crypt"
	"tiktok-demo/cmd/message/pkg/mysql"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/MessageServer"
)

func MessageActionResp(err errno.ErrNo) *MessageServer.DouyinMessageActionResponse {
	resp := new(MessageServer.DouyinMessageActionResponse)
	resp.BaseResp = &MessageServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	return resp
}

func MessageChatResp(err errno.ErrNo, messages []*MessageServer.ChatMessage) *MessageServer.DouyinMessageChatResponse {
	resp := new(MessageServer.DouyinMessageChatResponse)
	resp.BaseResp = &MessageServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.ChatList = messages
	return resp
}

func MessageLatestMsgResp(err errno.ErrNo, typeList []int32, contentList []string) *MessageServer.DouyinMessageLatestMsgResponse {
	resp := new(MessageServer.DouyinMessageLatestMsgResponse)
	resp.BaseResp = &MessageServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.TypeList = typeList
	resp.ContentList = contentList
	return resp
}

func BuildMessageActionResp(err error) *MessageServer.DouyinMessageActionResponse {
	if err == nil {
		return MessageActionResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return MessageActionResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return MessageActionResp(s)
}

func BuildMessageChatResp(err error, messages []*MessageServer.ChatMessage) *MessageServer.DouyinMessageChatResponse {
	if err == nil {
		return MessageChatResp(errno.Success, messages)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return MessageChatResp(e, nil)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return MessageChatResp(s, nil)
}

func BuildMessageLatestMsgResp(err error, typeList []int32, contentList []string) *MessageServer.DouyinMessageLatestMsgResponse {
	if err == nil {
		return MessageLatestMsgResp(errno.Success, typeList, contentList)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return MessageLatestMsgResp(e, nil, nil)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return MessageLatestMsgResp(s, nil, nil)
}

// mysql convert to rpc
func MsgInfoConvert(m []*mysql.Message) ([]*MessageServer.ChatMessage, error) {
	var messages []*MessageServer.ChatMessage
	for _, v := range m {
		DescryptContent, err := crypt.DecryptByAes(v.Content)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &MessageServer.ChatMessage{
			Id:         int64(v.ID),
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			Content:    string(DescryptContent),
			CreateTime: v.CreatedAt.Unix(),
		})
	}
	return messages, nil
}
