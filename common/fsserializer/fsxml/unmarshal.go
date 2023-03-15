/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: unmarshal to go object
@author: fanky
@version: 1.0
@date: 2022-07-08
**/

package fsxml

import (
	"fmt"
	"reflect"
)

func (*S_Node) unmarshalTo(rt reflect.Type, rv reflect.Value) error {
	return nil
}

func (this *S_Node) UnmarshalTo(obj any) error {
	tobj := reflect.TypeOf(obj)
	if tobj == nil {
		return fmt.Errorf("can't unmarshal to an inconclusive type nil value")
	}
	if tobj.Kind() != reflect.Ptr {
		return fmt.Errorf("the object unmarshal to must be a pinter")
	}
	return this.unmarshalTo(tobj, reflect.ValueOf(obj))
}
