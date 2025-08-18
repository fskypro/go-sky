/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 自定义类型键的集合
@author: fanky
@version: 1.0
@date: 2025-06-02
**/

package hset

import (
	"context"

	"fsky.pro/fsredis/codec"
	"github.com/redis/go-redis/v9"
)

// -------------------------------------------------------------------
// 自定义 Key 的约束接口
// -------------------------------------------------------------------
type I_KeyType[K any] interface {
	comparable
	KeyToString() string
	StringToKey(string) (K, error)
}

// -------------------------------------------------------------------
// AnyKeySet
// -------------------------------------------------------------------
type S_AnyKeySet[K I_KeyType[K], V any] struct {
	*s_HSet[K, V, S_AnyKeySet[K, V]]
}

func NewAnyKeySet[K I_KeyType[K], V any](client *redis.Client, hashKey string, codec codec.I_Codec[V]) *S_AnyKeySet[K, V] {
	return NewAnyKeySetContext[K](context.Background(), client, hashKey, codec)
}

func NewAnyKeySetContext[K I_KeyType[K], V any](ctx context.Context, client *redis.Client, hashKey string, codec codec.I_Codec[V]) *S_AnyKeySet[K, V] {
	return &S_AnyKeySet[K, V]{
		s_HSet: newHSet[K, V, S_AnyKeySet[K, V]](ctx, client, hashKey, codec),
	}
}

func (S_AnyKeySet[K, V]) keyToString(key K) string {
	// 这里不能使用方法接收者
	return key.KeyToString()
}

func (S_AnyKeySet[K, V]) stringToKey(str string) (K, error) {
	// 这里不能使用方法接收者
	var k K
	return k.StringToKey(str)
}
