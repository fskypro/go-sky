/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: service proxy for client
@author: fanky
@version: 1.0
@date: 2021-03-21
**/

package fsrpcs

import "net/rpc"

type I_ServiceProxy interface {
	Call(string, interface{}, interface{}) error
	Go(string, interface{}, interface{}, chan *rpc.Call) *rpc.Call
}
