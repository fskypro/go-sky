/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: encoder/decodeer
@author: fanky
@version: 1.0
@date: 2025-06-03
**/

package codec

// 数据编码解码器接口
type I_Codec[T any] interface {
	Encode(T) ([]byte, error)
	Decode([]byte, *T) error
}
