/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: errors
@author: fanky
@version: 1.0
@date: 2021-03-20
**/

package client

import (
	"fmt"
)

// -----------------------------------------------------------------------------
// Error type
// -----------------------------------------------------------------------------
// 远程请求失败错误
type S_ServerError string

func (e S_ServerError) Error() string {
	return string(e)
}

// 调用服务函数返回错误
type T_ServiceError string

func (e T_ServiceError) Error() string {
	return string(e)
}

// ---------------------------------------------------------
// 发送请求错误
// ---------------------------------------------------------
type S_SendError struct {
	msg string
}

func NewSendError(reqID, msg string, args ...interface{}) *S_SendError {
	return &S_SendError{fmt.Sprintf("send request(%s) message error: ", reqID) + fmt.Sprintf(msg, args...)}
}

func (this *S_SendError) Error() string {
	return string(this.msg)
}
