// Code generated by Kitex v0.5.2. DO NOT EDIT.

package commentservice

import (
	"context"
	"fmt"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
	CommentServer "tiktok-demo/shared/kitex_gen/CommentServer"
)

func serviceInfo() *kitex.ServiceInfo {
	return commentServiceServiceInfo
}

var commentServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "CommentService"
	handlerType := (*CommentServer.CommentService)(nil)
	methods := map[string]kitex.MethodInfo{
		"CommentAction": kitex.NewMethodInfo(commentActionHandler, newCommentActionArgs, newCommentActionResult, false),
		"CommentList":   kitex.NewMethodInfo(commentListHandler, newCommentListArgs, newCommentListResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "CommentServer",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.5.2",
		Extra:           extra,
	}
	return svcInfo
}

func commentActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(CommentServer.DouyinCommentActionRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(CommentServer.CommentService).CommentAction(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *CommentActionArgs:
		success, err := handler.(CommentServer.CommentService).CommentAction(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*CommentActionResult)
		realResult.Success = success
	}
	return nil
}
func newCommentActionArgs() interface{} {
	return &CommentActionArgs{}
}

func newCommentActionResult() interface{} {
	return &CommentActionResult{}
}

type CommentActionArgs struct {
	Req *CommentServer.DouyinCommentActionRequest
}

func (p *CommentActionArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(CommentServer.DouyinCommentActionRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *CommentActionArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *CommentActionArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *CommentActionArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, fmt.Errorf("No req in CommentActionArgs")
	}
	return proto.Marshal(p.Req)
}

func (p *CommentActionArgs) Unmarshal(in []byte) error {
	msg := new(CommentServer.DouyinCommentActionRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var CommentActionArgs_Req_DEFAULT *CommentServer.DouyinCommentActionRequest

func (p *CommentActionArgs) GetReq() *CommentServer.DouyinCommentActionRequest {
	if !p.IsSetReq() {
		return CommentActionArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *CommentActionArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *CommentActionArgs) GetFirstArgument() interface{} {
	return p.Req
}

type CommentActionResult struct {
	Success *CommentServer.DouyinCommentActionResponse
}

var CommentActionResult_Success_DEFAULT *CommentServer.DouyinCommentActionResponse

func (p *CommentActionResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(CommentServer.DouyinCommentActionResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *CommentActionResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *CommentActionResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *CommentActionResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, fmt.Errorf("No req in CommentActionResult")
	}
	return proto.Marshal(p.Success)
}

func (p *CommentActionResult) Unmarshal(in []byte) error {
	msg := new(CommentServer.DouyinCommentActionResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *CommentActionResult) GetSuccess() *CommentServer.DouyinCommentActionResponse {
	if !p.IsSetSuccess() {
		return CommentActionResult_Success_DEFAULT
	}
	return p.Success
}

func (p *CommentActionResult) SetSuccess(x interface{}) {
	p.Success = x.(*CommentServer.DouyinCommentActionResponse)
}

func (p *CommentActionResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *CommentActionResult) GetResult() interface{} {
	return p.Success
}

func commentListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(CommentServer.DouyinCommentListRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(CommentServer.CommentService).CommentList(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *CommentListArgs:
		success, err := handler.(CommentServer.CommentService).CommentList(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*CommentListResult)
		realResult.Success = success
	}
	return nil
}
func newCommentListArgs() interface{} {
	return &CommentListArgs{}
}

func newCommentListResult() interface{} {
	return &CommentListResult{}
}

type CommentListArgs struct {
	Req *CommentServer.DouyinCommentListRequest
}

func (p *CommentListArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(CommentServer.DouyinCommentListRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *CommentListArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *CommentListArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *CommentListArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, fmt.Errorf("No req in CommentListArgs")
	}
	return proto.Marshal(p.Req)
}

func (p *CommentListArgs) Unmarshal(in []byte) error {
	msg := new(CommentServer.DouyinCommentListRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var CommentListArgs_Req_DEFAULT *CommentServer.DouyinCommentListRequest

func (p *CommentListArgs) GetReq() *CommentServer.DouyinCommentListRequest {
	if !p.IsSetReq() {
		return CommentListArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *CommentListArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *CommentListArgs) GetFirstArgument() interface{} {
	return p.Req
}

type CommentListResult struct {
	Success *CommentServer.DouyinCommentListResponse
}

var CommentListResult_Success_DEFAULT *CommentServer.DouyinCommentListResponse

func (p *CommentListResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(CommentServer.DouyinCommentListResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *CommentListResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *CommentListResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *CommentListResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, fmt.Errorf("No req in CommentListResult")
	}
	return proto.Marshal(p.Success)
}

func (p *CommentListResult) Unmarshal(in []byte) error {
	msg := new(CommentServer.DouyinCommentListResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *CommentListResult) GetSuccess() *CommentServer.DouyinCommentListResponse {
	if !p.IsSetSuccess() {
		return CommentListResult_Success_DEFAULT
	}
	return p.Success
}

func (p *CommentListResult) SetSuccess(x interface{}) {
	p.Success = x.(*CommentServer.DouyinCommentListResponse)
}

func (p *CommentListResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *CommentListResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) CommentAction(ctx context.Context, Req *CommentServer.DouyinCommentActionRequest) (r *CommentServer.DouyinCommentActionResponse, err error) {
	var _args CommentActionArgs
	_args.Req = Req
	var _result CommentActionResult
	if err = p.c.Call(ctx, "CommentAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) CommentList(ctx context.Context, Req *CommentServer.DouyinCommentListRequest) (r *CommentServer.DouyinCommentListResponse, err error) {
	var _args CommentListArgs
	_args.Req = Req
	var _result CommentListResult
	if err = p.c.Call(ctx, "CommentList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
