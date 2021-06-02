/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: client
@author: fanky
@version: 1.0
@date: 2021-03-21
**/

package fsrpcs

import (
	"net/rpc"
)

// -------------------------------------------------------------------
// client base
// -------------------------------------------------------------------
type S_Client struct {
	ServiceName string
	proxy       I_ServiceProxy
}

func NewClient(name string, proxy I_ServiceProxy) *S_Client {
	return &S_Client{
		ServiceName: name,
		proxy:       proxy,
	}
}

func (this *S_Client) Call(method string, arg interface{}, reply interface{}) error {
	return this.proxy.Call(this.ServiceName+"."+method, arg, reply)
}

// 注意：cap(done) 一定要大于 1，否则会引起 panic
func (this *S_Client) Go(method string, arg interface{}, reply interface{}, done chan *rpc.Call) {
	this.proxy.Go(this.ServiceName+"."+method, arg, reply, done)
}

func (this *S_Client) Go2(method string, arg interface{}, reply interface{}) *rpc.Call {
	return this.proxy.Go(this.ServiceName+"."+method, arg, reply, nil)
}
