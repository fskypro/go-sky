/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: string key set
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
// AnyKeySet
// -------------------------------------------------------------------
type S_StrKeySet[K ~string, V any] struct {
	*s_HSet[K, V, S_StrKeySet[K, V]]
}

func NewStrKeySet[K ~string, V any](client *redis.Client, hashKey string, codec codec.I_Codec[V]) *S_StrKeySet[K, V] {
	return &S_StrKeySet[K, V]{
		s_HSet: newHSet[K, V, S_StrKeySet[K, V]](context.Background(), client, hashKey, codec),
	}
}

func (S_StrKeySet[K, V]) keyToString(key K) string {
	return string(key)
}

func (S_StrKeySet[K, V]) stringToKey(str string) (K, error) {
	return K(str), nil
}
