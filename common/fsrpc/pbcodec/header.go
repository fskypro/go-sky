/**
@copyright: fantasysky 2016
@brief: 请求和回复头扩展方法
@author: fanky
@version: 1.0
@date: 2018-09-18
**/

package pbcodec

import "fsky.pro/fsrpc"

// 将 rpc 内核格式请求头转换为 pb 格式
func toPBReqHeader(req *fsrpc.S_ReqHeader) *S_ReqHeader {
	return &S_ReqHeader{
		ServiceName: req.ServiceName,
		MethodName:  req.MethodName,
		ReqID:       req.ReqID,
	}
}

// 将 rpc 内核格式回复头转换为 pb 格式
func toPBRspHeader(rsp *fsrpc.S_RspHeader) *S_RspHeader {
	return &S_RspHeader{
		ServiceName: rsp.ServiceName,
		MethodName:  rsp.MethodName,
		ReqID:       rsp.ReqID,
		Fail:        rsp.Fail,
		Error:       rsp.Error,
	}
}

// ------------------------------------------------------------------
// 将 pb 格式请求头转换为 rpc 内核格式
func (pbReq *S_ReqHeader) toReqHeader(req *fsrpc.S_ReqHeader) {
	req.ServiceName = pbReq.ServiceName
	req.MethodName = pbReq.MethodName
	req.ReqID = pbReq.ReqID
}

// 将 pb 格式回复头转换为 rpc 内核格式
func (pbRsp *S_RspHeader) toRspHeader(rsp *fsrpc.S_RspHeader) {
	rsp.ReqID = pbRsp.ReqID
	rsp.Fail = pbRsp.Fail
	rsp.Error = pbRsp.Error
}
