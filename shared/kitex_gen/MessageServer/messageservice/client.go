// Code generated by Kitex v0.5.2. DO NOT EDIT.

package messageservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	MessageServer "tiktok-demo/shared/kitex_gen/MessageServer"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	MessageChat(ctx context.Context, Req *MessageServer.DouyinMessageChatRequest, callOptions ...callopt.Option) (r *MessageServer.DouyinMessageChatResponse, err error)
	MessageAction(ctx context.Context, Req *MessageServer.DouyinMessageActionRequest, callOptions ...callopt.Option) (r *MessageServer.DouyinMessageActionResponse, err error)
	MessageLatestMsg(ctx context.Context, Req *MessageServer.DouyinMessageLatestMsgRequest, callOptions ...callopt.Option) (r *MessageServer.DouyinMessageLatestMsgResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kMessageServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kMessageServiceClient struct {
	*kClient
}

func (p *kMessageServiceClient) MessageChat(ctx context.Context, Req *MessageServer.DouyinMessageChatRequest, callOptions ...callopt.Option) (r *MessageServer.DouyinMessageChatResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MessageChat(ctx, Req)
}

func (p *kMessageServiceClient) MessageAction(ctx context.Context, Req *MessageServer.DouyinMessageActionRequest, callOptions ...callopt.Option) (r *MessageServer.DouyinMessageActionResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MessageAction(ctx, Req)
}

func (p *kMessageServiceClient) MessageLatestMsg(ctx context.Context, Req *MessageServer.DouyinMessageLatestMsgRequest, callOptions ...callopt.Option) (r *MessageServer.DouyinMessageLatestMsgResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MessageLatestMsg(ctx, Req)
}