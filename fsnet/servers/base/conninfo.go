/**
@copyright: fantasysky 2016
@brief: 连接信息基础结构
@author: fanky
@version: 1.0
@date: 2019-03-19
**/

package base

import "fmt"

// 客户端基础信息
type S_ConnInfo struct {
	ConnID   T_ConnID    // 连接编号
	IP       string      // 客户端IP地址
	Port     uint16      // 客户端端口
	UserData interface{} // 自定义用户额外信息

	state T_ConnState // 是否是被服务器踢下线（只有在离线的时候使用）
}

func NewConnInfo(ip string, port uint16) *S_ConnInfo {
	return &S_ConnInfo{
		ConnID: GenConnID(),
		IP:     ip,
		Port:   port,
		state:  CONN_STATE_NEW,
	}
}

func (this *S_ConnInfo) String() string {
	return fmt.Sprintf(`{connid: %d, addr:"%s:%d", state="%v"}`, this.ConnID, this.IP, this.Port, this.state)
}

// IsState 判断是否是指定状态
// 调用该方法前，需要先调用 LockState 方法获得锁，并在调用 UnlockState 解锁
func (this *S_ConnInfo) IsState(state T_ConnState) bool {
	return this.state == state
}

// SetState 设置连接状态
// 调用该方法前，需要先调用 LockState 方法获得锁，并在调用 UnlockState 解锁
func (this *S_ConnInfo) SetState(state T_ConnState) {
	this.state = state
}
