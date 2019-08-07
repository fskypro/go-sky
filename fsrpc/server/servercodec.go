/**
@copyright: fantasysky 2016
@brief: 服务器端编码解码器接口
@author: fanky
@version: 1.0
@date: 2018-09-22
**/

package server

import "io"
import . "fsky.pro/fsrpc"

// 所有解码结构都必须实现该接口
type S_ServerCodec interface {
	// 读取请求头
	ReadRequestHeader(*S_ReqHeader) error

	// 读取请求数据体
	// 参数如果是 nil，则表示解释请求头失败，不需要解释参数，直接将读出的数据丢弃即可
	ReadRequestArg(*S_ReqHeader, interface{}) error

	// 将回复客户端内容写入网络数据流
	// 第二个参数如果是 nil，则表示远程调用失败；如果是 *fsrpc.EReply 则表示该远程调用不关注返回参数
	WriteResponse(*S_RspHeader, interface{}) error

	// 数据流通道
	Close() error
}

// 编码解码器生成函数
type F_CodecCreator func(io.ReadWriteCloser) S_ServerCodec
