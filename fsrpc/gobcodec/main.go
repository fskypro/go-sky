/**
* @copyright: fantasysky 2016
* @brief: 实现 gob 数据打包和解包
* @author: fanky
* @version: 1.0
* @date: 2018-08-30
**/

// gob 数据打包和解包
package gobcodec

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"

	"fsky.pro/fsrpc"
)

// -----------------------------------------------------------------------------
// S_CobServerCodecc
// -----------------------------------------------------------------------------
type S_CobServerCodecc struct {
	rwc     io.ReadWriteCloser
	dec     *gob.Decoder
	enc     *gob.Encoder
	encBuff *bufio.Writer
	closed  bool
}

func NewServerCodec(rwc io.ReadWriteCloser) *S_CobServerCodecc {
	iowriter := bufio.NewWriter(rwc)
	return &S_CobServerCodecc{
		rwc:     rwc,
		dec:     gob.NewDecoder(rwc),      // 用于读取并解码的解码器
		enc:     gob.NewEncoder(iowriter), // 用于编码并写入缓冲的编码器
		encBuff: iowriter,                 // 写入缓冲
		closed:  false,
	}
}

// ------------------------------------------------------------------
// 读取请求消息头
func (c *S_CobServerCodecc) ReadRequestHeader(r *fsrpc.S_ReqHeader) error {
	return c.dec.Decode(r)
}

// 读取请求消息体
func (c *S_CobServerCodecc) ReadRequestArg(header *fsrpc.S_ReqHeader, arg interface{}) error {
	return c.dec.Decode(arg)
}

// 回复客户端
func (c *S_CobServerCodecc) WriteResponse(header *fsrpc.S_RspHeader, reply interface{}) (err error) {
	if err = c.enc.Encode(header); err != nil {
		// 理论上不会跑这里来，如果确实跑这里来了，需要刷新缓冲，关闭链接
		if c.encBuff.Flush() == nil {
			err = errors.New("encode gob's reponse header fail: " + err.Error())
			c.Close()
		}
		return
	}
	if err = c.enc.Encode(reply); err != nil {
		// 理论上不会跑这里来，如果确实跑这里来了，需要刷新缓冲，关闭链接
		if c.encBuff.Flush() == nil {
			err = errors.New("encode gob's response reply fail: " + err.Error())
			c.Close()
		}
		return
	}
	c.encBuff.Flush()
	return
}

// 关闭IO
func (c *S_CobServerCodecc) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

// -----------------------------------------------------------------------------
// S_GobClientCodec
// -----------------------------------------------------------------------------
type S_GobClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
}

func NewClienCodec() *S_GobClientCodec {
	return &S_GobClientCodec{}
}

// ------------------------------------------------------------------
// 初始化
func (c *S_GobClientCodec) Initialize(rwc io.ReadWriteCloser) {
	encBuf := bufio.NewWriter(rwc)
	c.rwc = rwc
	c.dec = gob.NewDecoder(rwc)
	c.enc = gob.NewEncoder(encBuf)
	c.encBuf = encBuf
}

// 写入请求数据
func (c *S_GobClientCodec) WriteRequest(header *fsrpc.S_ReqHeader, arg interface{}) (err error) {
	if err = c.enc.Encode(header); err != nil {
		err = errors.New("encode gob's request header fail: " + err.Error())
		return
	}
	if err = c.enc.Encode(arg); err != nil {
		err = errors.New("fsrpc: encode gob's request argument fail: " + err.Error())
		c.encBuf.Reset(c.rwc)
		return
	}
	c.encBuf.Flush()
	return
}

// 读取回复数据头
func (c *S_GobClientCodec) ReadResponseHeader(header *fsrpc.S_RspHeader) error {
	err := c.dec.Decode(header)
	if err != nil {
		err = errors.New("fsrpc: decode gob's response header fail: " + err.Error())
	}
	return err
}

// 读取回复数据体
func (c *S_GobClientCodec) ReadResponseReply(reply interface{}) error {
	err := c.dec.Decode(reply)
	if err != nil {
		err = errors.New("fsrpc: decode gob's response body fail: " + err.Error())
	}
	return err
}

// 关闭链接
func (c *S_GobClientCodec) Close() error {
	return c.rwc.Close()
}
