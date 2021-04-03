/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: request info
@author: fanky
@version: 1.0
@date: 2021-03-20
**/

package client

import (
	"time"
)

// -----------------------------------------------------------------------------
// request info
// -----------------------------------------------------------------------------
// 请求信息
type S_ReqInfo struct {
	ID          string          // 请求编号
	ServiceName string          // 请求的服务对象名称
	MethodName  string          // 请求的方法
	Arg         interface{}     // 请求参数
	Reply       interface{}     // 请求回复参数
	Error       error           // 请求失败或调用服务器服务函数返回错误
	ReqCh       chan *S_ReqInfo // 请求返回后数据放入的通道
}

// 结束远程调用，将结构写入用户通道
func (this *S_ReqInfo) send() {
	select {
	case this.ReqCh <- this:
	default:
		// 如果 req.ReqCh 阻塞（已满），则程序会跑这里来，
		// 这时，等待 100 毫秒再尝试往通道发送
		go func() {
			time.Sleep(time.Millisecond * 100)
			select {
			case this.ReqCh <- this:
			default:
				// 如果 100 毫秒后，还是无法发送，则丢弃该请求
				// 这样做的目的是防止程序被无限期阻塞，因此 req.ReqCh 必须设置足够的缓冲数量
				this.Error = NewSendError(this.ID, "a reqiest(%s.%s) has been discarded!", this.ServiceName, this.MethodName)
			}
		}()
	}
}
