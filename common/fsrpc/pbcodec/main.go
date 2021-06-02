/**
@copyright: fantasysky 2016
@brief: 实现 protobuf 编码解码
@author: fanky
@version: 1.0
@date: 2018-09-16
**/

// protobuf 数据打包和解包
// 使用该包，需要用到 protobuf 编译工具
// 工具下载地址：https://github.com/google/protobuf/releases
// 同时用到 go protobuf 库：https://github.com/golang/protobuf.git
// 将 proto 文件转换为 go 代码文件：protoc --go_out=./ *.proto
package pbcodec

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"

	"fsky.pro/fslog"
	"fsky.pro/fsrpc"
)

// -----------------------------------------------------------------------------
// innter functions
// -----------------------------------------------------------------------------
// 读取网络缓存数据
func _readBuffer(r io.Reader) (buff []byte, err error) {
	sizeBuff := make([]byte, 4, 4)
	n, err := r.Read(sizeBuff)
	if err != nil {
		return
	}
	if n < 4 {
		err = errors.New("not enough content to be read in buffer.")
		return
	}

	var size uint32
	err = binary.Read(bytes.NewBuffer(sizeBuff), binary.BigEndian, &size)
	if err != nil {
		return
	}
	buff = make([]byte, size, size)
	_, err = r.Read(buff)
	return
}

// 写入数据到网络缓存
func _writeBuffer(w io.Writer, buff []byte) error {
	size := uint32(len(buff))
	err := binary.Write(w, binary.BigEndian, size)
	if err != nil {
		return err
	}
	_, err = w.Write(buff)
	return err
}

// -----------------------------------------------------------------------------
// S_PBServerCodec
// -----------------------------------------------------------------------------
type S_PBServerCodec struct {
	rwc    io.ReadWriteCloser
	writer *bufio.Writer
	closed bool
}

func NewServerCodec(rwc io.ReadWriteCloser) *S_PBServerCodec {
	iowriter := bufio.NewWriter(rwc)
	return &S_PBServerCodec{
		rwc:    rwc,
		writer: iowriter,
		closed: false,
	}
}

// ------------------------------------------------------------------
// 读取请求消息头
func (c *S_PBServerCodec) ReadRequestHeader(r *fsrpc.S_ReqHeader) error {
	header := S_ReqHeader{}
	buf, err := _readBuffer(c.rwc)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(buf, &header)
	if err != nil {
		return err
	}
	header.toReqHeader(r)
	return nil
}

// 读取请求消息体
func (c *S_PBServerCodec) ReadRequestArg(header *fsrpc.S_ReqHeader, arg interface{}) error {
	buff, err := _readBuffer(c.rwc)
	if err != nil {
		return err
	}

	// 远程调用失败
	if arg == nil {
		return nil
	}

	// 空参数
	if fsrpc.IsEmptyArg(arg) {
		if len(buff) == 0 {
			return nil
		} else {
			fslog.Warnf("fsrpc: servervice '%s.%s' require an empty argument, "+
				"but client pass an unempty argument!", header.ServiceName, header.MethodName)
		}
		return nil
	}
	return proto.Unmarshal(buff, arg.(proto.Message))
}

// 回复客户端
func (c *S_PBServerCodec) WriteResponse(header *fsrpc.S_RspHeader, reply interface{}) error {
	// 写入回复头
	pbRsp := toPBRspHeader(header)
	buff, err := proto.Marshal(pbRsp)
	if err != nil {
		c.Close()
		return err
	}
	if err = _writeBuffer(c.writer, buff); err != nil {
		c.Close()
		return err
	}

	// 调用失败或者无返回值
	if reply == nil || fsrpc.IsEmptyReply(reply) {
		buff = []byte{}
	} else {
		// 写入回复结果数据
		buff, err = proto.Marshal(reply.(proto.Message))
		if err != nil {
			c.Close()
			return err
		}
	}
	if err = _writeBuffer(c.writer, buff); err != nil {
		c.Close()
		return err
	}
	c.writer.Flush()
	return nil
}

// 关闭IO
func (c *S_PBServerCodec) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

// -----------------------------------------------------------------------------
// S_PBClientCodec
// -----------------------------------------------------------------------------
type S_PBClientCodec struct {
	rwc    io.ReadWriteCloser
	writer *bufio.Writer
}

func NewClienCodec() *S_PBClientCodec {
	return &S_PBClientCodec{}
}

// ------------------------------------------------------------------
// 初始化
func (c *S_PBClientCodec) Initialize(rwc io.ReadWriteCloser) {
	c.rwc = rwc
	c.writer = bufio.NewWriter(rwc)
}

// 写入请求数据
func (c *S_PBClientCodec) WriteRequest(header *fsrpc.S_ReqHeader, arg interface{}) error {
	// 写入请求头
	pbHeader := toPBReqHeader(header)
	bheader, err := proto.Marshal(pbHeader)
	if err != nil {
		return errors.New("encode protobuf's request header fail: " + err.Error())
	}
	if err = _writeBuffer(c.writer, bheader); err != nil {
		return errors.New("write header fail: " + err.Error())
	}

	// 写入请求内容
	var body []byte
	if arg == nil || fsrpc.IsEmptyArg(arg) {
		body = []byte{}
	} else {
		body, err = proto.Marshal(arg.(proto.Message))
	}
	if err != nil {
		return errors.New("encode protobuf's request argument fail: " + err.Error())
	}
	if err = _writeBuffer(c.writer, body); err != nil {
		c.writer.Reset(c.rwc)
		return errors.New("write argument fail: " + err.Error())
	}
	c.writer.Flush()
	return nil
}

// 读取回复数据头
func (c *S_PBClientCodec) ReadResponseHeader(header *fsrpc.S_RspHeader) error {
	pbHeader := new(S_RspHeader)
	buff, err := _readBuffer(c.rwc)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(buff, pbHeader)
	if err != nil {
		return err
	}
	pbHeader.toRspHeader(header)
	return nil
}

// 读取回复数据体
func (c *S_PBClientCodec) ReadResponseReply(reply interface{}) error {
	buff, err := _readBuffer(c.rwc)
	if err != nil {
		return err
	}

	// 不关注返回值
	if reply == nil || fsrpc.IsEmptyReply(reply) {
		return nil
	}

	return proto.Unmarshal(buff, reply.(proto.Message))
}

// 关闭链接
func (c *S_PBClientCodec) Close() error {
	return c.rwc.Close()
}
