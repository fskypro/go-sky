/**
@copyright: yxht 2015
@website: http://www.exht.com
@brief: json encoder/decoder
@author: fanky
@version: 1.0
@date: 2025-06-02
**/

package jsoncodec

import "encoding/json"

type S_JsonCodec[T any] struct{}

func (S_JsonCodec[T]) Encode(obj T) ([]byte, error) {
	return json.Marshal(obj)
}

func (S_JsonCodec[T]) Decode(data []byte, obj *T) error {
	return json.Unmarshal(data, obj)
}
