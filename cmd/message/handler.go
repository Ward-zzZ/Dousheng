package main

import (
	"context"

	"tiktok-demo/cmd/message/pkg/crypt"
	"tiktok-demo/cmd/message/pkg/mysql"
	"tiktok-demo/cmd/message/pkg/pack"
	"tiktok-demo/shared/errno"
	MessageServer "tiktok-demo/shared/kitex_gen/MessageServer"

	"github.com/cloudwego/kitex/pkg/klog"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct {
	MysqlManager
	// RedisManager
}

type MysqlManager interface {
	AddMessage(fromUserId int64, toUserId int64, content string) (*mysql.Message, error)
	GetMessageByUserId(userId int64, toUserId int64, timeStamp int64) ([]*mysql.Message, error)
	GetMessageById(id int64) (*mysql.Message, error)
	GetLatestMessage(userId int64, toUserId int64) (*mysql.Message, error)
}

type RedisManager interface {
	SetLastMsgTime(userId int64, toUserId int64, timeStamp int64) error
	GetLastMsgTime(userId int64, toUserId int64) (int64, error)
}

// MessageChat implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageChat(ctx context.Context, req *MessageServer.DouyinMessageChatRequest) (resp *MessageServer.DouyinMessageChatResponse, err error) {

	// // redis 读取上次获取消息的最后一条消息的时间戳
	// lastMsgTime, err := s.RedisManager.GetLastMsgTime(req.UserId, req.ToUserId)
	if err != nil {
		klog.Errorf("Redis GetLastMsgTime err:%v", err)
		return pack.BuildMessageChatResp(errno.InternalServerErr, nil), nil
	}

	var messages []*mysql.Message
	lastMsgTime := req.PreMsgTime

	messages, err = s.MysqlManager.GetMessageByUserId(req.UserId, req.ToUserId, lastMsgTime)
	if err != nil {
		klog.Errorf("Mysql GetMessageByUserId err:%v", err)
		return pack.BuildMessageChatResp(errno.InternalServerErr, nil), nil
	}

	// // 更新redis中的最后一条消息的时间戳
	// if len(messages) > 0 {
	// 	lastMsgTime = messages[len(messages)-1].CreatedAt.Unix()
	// 	if err := s.RedisManager.SetLastMsgTime(req.UserId, req.ToUserId, lastMsgTime); err != nil {
	// 		klog.Errorf("Redis SetLastMsgTime err:%v", err)
	// 		return pack.BuildMessageChatResp(errno.InternalServerErr, nil), nil
	// 	}
	// }

	msg, err := pack.MsgInfoConvert(messages)
	if err != nil {
		klog.Errorf("MsgInfoConvert err:%v", err)
		return pack.BuildMessageChatResp(errno.InternalServerErr, nil), nil
	}
	return pack.BuildMessageChatResp(nil, msg), nil

}

// MessageAction implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageAction(ctx context.Context, req *MessageServer.DouyinMessageActionRequest) (resp *MessageServer.DouyinMessageActionResponse, err error) {
	if (req.ActionType != 1 && req.ActionType != 2) || req.UserId == req.ToUserId {
		return pack.BuildMessageActionResp(errno.ParamErr), nil
	}
	//todo: relation查询判断是否好友
	// aes+base64加密消息内容
	encryptContent, err := crypt.EncryptByAes([]byte(req.Content))
	if err != nil {
		klog.Errorf("EncryptByAes err:%v", err)
		return pack.BuildMessageActionResp(errno.InternalServerErr), nil
	}
	// add message to mysql
	if req.ActionType == 1 {
		_, err := s.MysqlManager.AddMessage(req.UserId, req.ToUserId, string(encryptContent))
		if err != nil {
			klog.Errorf("Mysql AddMessage err:%v", err)
			return pack.BuildMessageActionResp(errno.InternalServerErr), nil
		}
	}

	return pack.BuildMessageActionResp(nil), nil
}

// MessageLatestMsg implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageLatestMsg(ctx context.Context, req *MessageServer.DouyinMessageLatestMsgRequest) (resp *MessageServer.DouyinMessageLatestMsgResponse, err error) {

	// check params
	if req.UserId == 0 || len(req.ToUserIdList) == 0 {
		return pack.BuildMessageLatestMsgResp(errno.ParamErr, nil, nil), nil
	}
	typeList := make([]int32, len(req.ToUserIdList))
	contentList := make([]string, len(req.ToUserIdList))
	for i, toUserId := range req.ToUserIdList {
		latestMsg, err := s.MysqlManager.GetLatestMessage(req.UserId, toUserId)
		if err != nil {
			return pack.BuildMessageLatestMsgResp(errno.InternalServerErr, nil, nil), nil
		}
		if latestMsg == nil {
			continue
		}
		// aes+base64解密消息内容
		decryptContent, err := crypt.DecryptByAes(latestMsg.Content)
		if err != nil {
			klog.Errorf("DecryptByAes err:%v", err)
			return pack.BuildMessageLatestMsgResp(errno.InternalServerErr, nil, nil), nil
		}
		if latestMsg.FromUserId == req.UserId {
			typeList[i] = 1 // 当前用户发送的消息
		} else {
			typeList[i] = 0 // 当前用户接收的消息
		}
		contentList[i] = string(decryptContent)
	}
	return pack.BuildMessageLatestMsgResp(nil, typeList, contentList), nil
}
