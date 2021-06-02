/**
@copyright: fantasysky 2016
@brief: 客户端编码解码器接口
@author: fanky
@version: 1.0
@date: 2018-09-22
**/

package client

import "io"
import "fsky.pro/fsrpc"

type I_ClientCodec interface {
	// 初始化 codec
	Initialize(rwc io.ReadWriteCloser)

	// 将请求数据写入网络数据流
	// 第二个参数为请求参数，如果服务方法不需要参数，则服务器方法定义参数为 fsrpc.EArg 或 *fsrpc.EArg
	// 客户端调用时可以给 fsrpc.EArg 传入 nil 或者 fsrpc.EArg{}
	WriteRequest(*fsrpc.S_ReqHeader, interface{}) error

	// 读取回复数据头
	ReadResponseHeader(*fsrpc.S_RspHeader) error

	// 读取回复数据体
	// 如果客户端不需要关注回复内容，则可以对参数传入 nil 或者 fsrpc.EReply{}
	ReadResponseReply(interface{}) error

	// 关闭编码解码器
	Close() error
}
