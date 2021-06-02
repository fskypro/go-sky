/**
@copyright: fantasysky 2016
@brief: 生成连接 ID
@author: fanky
@version: 1.0
@date: 2019-01-09
**/

package base

import "math"

// ---------------------------------------------------------
// 连接编号
// ---------------------------------------------------------
type T_ConnID uint64

var _currMax T_ConnID = 0

// 单机的话可用该函数产生唯一ID
func GenConnID() T_ConnID {
	if _currMax < math.MaxUint64 {
		_currMax += 1
	} else {
		_currMax = 1
	}
	return _currMax
}

// ---------------------------------------------------------
// 连接状态
// ---------------------------------------------------------
type T_ConnState uint8

// 连接状态
const (
	CONN_STATE_NEW     T_ConnState = iota // 新建状态
	CONN_STATE_ONLINE                     // 连接中
	CONN_STATE_KICKOUT                    // 被踢下线
	CONN_STATE_LOST                       // 客户端主动离线
)

var _strStates = map[T_ConnState]string{
	CONN_STATE_NEW:     `CONN_STATE_NEW`,
	CONN_STATE_ONLINE:  `CONN_STATE_ONLINE`,
	CONN_STATE_KICKOUT: `CONN_STATE_KICKOUT`,
	CONN_STATE_LOST:    `CONN_STATE_LOST`,
}

func (self T_ConnState) String() string {
	return _strStates[self]
}

// ---------------------------------------------------------
// 回调
// ---------------------------------------------------------
// 上线回调
type F_OnlineCb func(*S_ConnInfo)

// 离线回调
type F_OfflineCb func(*S_ConnInfo)

// 接收消息回调
type F_ReceiveCb func(*S_ConnInfo, []byte)

// 遍历连接信息回调
type F_IterConnsHandler func(*S_ConnInfo) bool
