/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: heartbeat client
@author: fanky
@version: 1.0
@date: 2021-03-26
**/

package tcprpc

import (
	"context"
	"net"
	"time"

	"fsky.pro/fsrpcs"
)

// -------------------------------------------------------------------
// 心跳服务
// -------------------------------------------------------------------
const _heartbeatServiceName = "HeartbeatService_"

type s_HeartbeatService struct {
	Reply []byte
}

func newHeartbeatService() *s_HeartbeatService {
	return &s_HeartbeatService{[]byte("hello")}
}

func (this *s_HeartbeatService) Name() string {
	return _heartbeatServiceName
}

func (this *s_HeartbeatService) Hello(req []byte, reply *[]byte) error {
	*reply = this.Reply
	return nil
}

// -------------------------------------------------------------------
// 心跳客户端
// -------------------------------------------------------------------
// 心跳参数
type S_HeartbeatInfo struct {
	Interval int         // 发送心跳包时间间隔（秒）
	Request  []byte      // 心跳表内容
	Respond  chan []byte // 心跳回复内容
}

func NewHeartbeatInfo() *S_HeartbeatInfo {
	return &S_HeartbeatInfo{
		Interval: 10,
		Request:  []byte("hello!"),
		Respond:  nil,
	}
}

// ---------------------------------------------------------
type s_HeartbeatClient struct {
	*fsrpcs.S_Client
	info *S_HeartbeatInfo
}

func newHeartbeatClient(proxy fsrpcs.I_ServiceProxy) *s_HeartbeatClient {
	return &s_HeartbeatClient{
		S_Client: fsrpcs.NewClient(_heartbeatServiceName, proxy),
		info:     NewHeartbeatInfo(),
	}
}

func (this *s_HeartbeatClient) updateInfo(info *S_HeartbeatInfo) {
	if info.Interval < 1 {
		info.Interval = 1
	}
	this.info = info
}

func (this *s_HeartbeatClient) hello() error {
	var reply []byte
	err := this.Call("Hello", this.info.Request, &reply)
	if err == nil && this.info.Respond != nil {
		select {
		case <-time.After(time.Second * 2):
			break
		case this.info.Respond <- reply:
			break
		}
	}
	return err
}

func (this *s_HeartbeatClient) loopSend(ctx context.Context, conn *net.TCPConn) {
	interval := time.Second * time.Duration(this.info.Interval)
	//outtime := time.Now().Add(interval).Add(time.Second * 3)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			//conn.SetReadDeadline(outtime)
			if this.hello() != nil {
				return
			}
			time.Sleep(interval)
		}
	}
}
