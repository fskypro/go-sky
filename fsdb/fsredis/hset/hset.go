/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: redis hset
@author: fanky
@version: 1.0
@date: 2025-06-01
**/

package hset

import (
	"context"
	"fmt"
	"time"

	"fsky.pro/fsredis/codec"

	"fsky.pro/fsredis/codec/jsoncodec"

	"github.com/redis/go-redis/v9"
)

// -------------------------------------------------------------------
// 子类约束
// -------------------------------------------------------------------
type i_Sub[K any] interface {
	keyToString(K) string
	stringToKey(string) (K, error)
}

// -------------------------------------------------------------------
// 返回值封装
// -------------------------------------------------------------------
// 返回 key 列表
type S_RetKeys[K any] struct {
	InvalidKeys []string // 数据库中存在，但是无法转换成自定义 key 的字符串形式 key
	Keys        []K      // 获取成功的 key 列表
}

// 返回单个值
type S_RetValue[V any] struct {
	Exists    bool  // 键是否存在
	DecodeErr error // 解码值错误

	Value V      // 解释成功的值
	Data  []byte // 数据库中的原始存储值
}

// 返回键值对集合
type S_RetItems[K comparable, V any] struct {
	UnexistKeys []K               // 数据库中不存在的 key 列表
	InvalidKeys []string          // 数据库中存在，但是无法解释的 key 列表
	Items       map[K]V           // 获取成功的键值对集合
	DecodeErrs  map[K]error       // 对值解码失败的键和解码错误映射
	Datas       map[string][]byte // 数据库中的原始数据
}

// 迭代键
type S_IterKey[K comparable] struct {
	IsValid bool   // 是否解释键成功
	DBKey   string // 数据库中的键
	Key     K
}

// 迭代键值对
type S_IterItem[K comparable, V any] struct {
	IsKeyValid bool   // 还原 key 是否成功
	DecodeErr  error  // 解码错误
	DBKey      string // 数据库中原始 Key
	Key        K      // 自定义键
	Value      V      // 键对应的值
}

// -------------------------------------------------------------------
// HSet
// -------------------------------------------------------------------
type s_HSet[K comparable, V any, Sub i_Sub[K]] struct {
	ctx     context.Context
	client  *redis.Client
	hashKey string
	codec   codec.I_Codec[V]

	sub Sub
}

func newHSet[K comparable, V any, Sub i_Sub[K]](ctx context.Context, client *redis.Client, hashKey string, codec codec.I_Codec[V]) *s_HSet[K, V, Sub] {
	if codec == nil {
		codec = jsoncodec.S_JsonCodec[V]{}
	}
	var sub Sub
	return &s_HSet[K, V, Sub]{
		ctx:     ctx,
		client:  client,
		hashKey: hashKey,
		codec:   codec,
		sub:     sub,
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *s_HSet[K, V, Sub]) Clent() *redis.Client {
	return this.client
}

// 获取 set 中键值对的总数
func (this *s_HSet[K, V, Sub]) CountCtx(ctx context.Context) (int64, error) {
	return this.client.HLen(ctx, this.hashKey).Result()
}

// 获取 set 中键值对的总数
func (this *s_HSet[K, V, Sub]) Count() (int64, error) {
	return this.client.HLen(this.ctx, this.hashKey).Result()
}

// 查找是否存在指定 key
func (this *s_HSet[K, V, Sub]) HasKeyCtx(ctx context.Context, key K) (bool, error) {
	return this.client.HExists(ctx, this.hashKey, this.sub.keyToString(key)).Result()
}

// 查找是否存在指定 key
func (this *s_HSet[K, V, Sub]) HasKey(key K) (bool, error) {
	return this.HasKeyCtx(this.ctx, key)
}

// 获取所有键
func (this *s_HSet[K, V, Sub]) GetKeysCtx(ctx context.Context) (S_RetKeys[K], error) {
	fields, err := this.client.HKeys(ctx, this.hashKey).Result()
	if err != nil {
		return S_RetKeys[K]{}, err
	}
	loses := []string{}
	keys := []K{}
	for _, strKey := range fields {
		key, err := this.sub.stringToKey(strKey)
		if err != nil {
			loses = append(loses, strKey)
		} else {
			keys = append(keys, key)
		}
	}
	return S_RetKeys[K]{InvalidKeys: loses, Keys: keys}, nil
}

// 获取所有键
func (this *s_HSet[K, V, Sub]) GetKeys() (S_RetKeys[K], error) {
	return this.GetKeysCtx(this.ctx)
}

// 获取映射值
// 注意：
//
//	返回错误不等于 nil，只是表示访问数据库失败
//	还需要判断 S_RetValue[V].DecodeErr 是否等于 nil，不等于 nil，则表示反序列化值错误
func (this *s_HSet[K, V, Sub]) GetCtx(ctx context.Context, key K) (S_RetValue[V], error) {
	ret := S_RetValue[V]{}
	strKey := this.sub.keyToString(key)
	data, err := this.client.HGet(ctx, this.hashKey, strKey).Bytes()
	if err == redis.Nil {
		return ret, nil
	}
	if err != nil {
		return ret, err
	}

	ret.Exists = true
	ret.Data = data
	if err := this.codec.Decode(data, &ret.Value); err != nil {
		ret.DecodeErr = fmt.Errorf("decode value of str-key %q fail, %v", strKey, err)
	}
	return ret, nil
}

// 获取映射值
// 注意：
//
//	返回错误不等于 nil，只是表示访问数据库失败
//	还需要判断 S_RetValue[V].DecodeErr 是否等于 nil，不等于 nil，则表示反序列化值错误
func (this *s_HSet[K, V, Sub]) Get(key K) (S_RetValue[V], error) {
	return this.GetCtx(this.ctx, key)
}

// 获取指定一组 key 的值
func (this *s_HSet[K, V, Sub]) GetsCtx(ctx context.Context, keys ...K) (S_RetItems[K, V], error) {
	ret := S_RetItems[K, V]{
		UnexistKeys: []K{},
		InvalidKeys: make([]string, 0),
		Items:       map[K]V{},
		DecodeErrs:  map[K]error{},
		Datas:       map[string][]byte{},
	}
	strKeys := []string{}
	for _, key := range keys {
		strKeys = append(strKeys, this.sub.keyToString(key))
	}
	datas, err := this.client.HMGet(ctx, this.hashKey, strKeys...).Result()
	if err != nil {
		return ret, err
	}
	for index, key := range keys {
		var value V
		data := datas[index]
		if data == nil {
			continue
		}
		if _, ok := data.(string); !ok {
			continue
		}
		bdata := []byte(data.(string))
		ret.Datas[strKeys[index]] = bdata
		err := this.codec.Decode(bdata, &value)
		if err != nil {
			ret.DecodeErrs[key] = err
		} else {
			ret.Items[key] = value
		}
	}
	return ret, nil
}

// 获取指定一组 key 的值
func (this *s_HSet[K, V, Sub]) Gets(keys ...K) (S_RetItems[K, V], error) {
	return this.GetsCtx(this.ctx, keys...)
}

// 获取并删除指定键值
func (this *s_HSet[K, V, Sub]) GetDelsCtx(ctx context.Context, keys ...K) (S_RetItems[K, V], error) {
	ret := S_RetItems[K, V]{
		UnexistKeys: []K{},
		InvalidKeys: make([]string, 0),
		Items:       map[K]V{},
		DecodeErrs:  map[K]error{},
		Datas:       map[string][]byte{},
	}
	if len(keys) == 0 {
		return ret, nil
	}
	strKeys := []string{}
	for _, key := range keys {
		strKeys = append(strKeys, this.sub.keyToString(key))
	}
	datas, err := this.client.HGetDel(ctx, this.hashKey, strKeys...).Result()
	if err != nil {
		return ret, err
	}
	for index, data := range datas {
		var key = keys[index]
		var strKey = strKeys[index]
		ret.Datas[strKey] = []byte(data)
		if data == "" {
			// key 不存在是，data 为空字符串
			ret.UnexistKeys = append(ret.UnexistKeys, key)
			continue
		}
		var value V
		if err := this.codec.Decode([]byte(data), &value); err != nil {
			ret.DecodeErrs[key] = fmt.Errorf("decode value of str-key %q fail, %v", strKey, err)
		} else {
			ret.Items[key] = value
		}
	}
	return ret, nil
}

// 获取并删除指定键值
func (this *s_HSet[K, V, Sub]) GetDels(keys ...K) (S_RetItems[K, V], error) {
	return this.GetDelsCtx(this.ctx, keys...)
}

// ---------------------------------------------------------
// 获取全部映射值
func (this *s_HSet[K, V, Sub]) GetAllCtx(ctx context.Context) (S_RetItems[K, V], error) {
	ret := S_RetItems[K, V]{
		UnexistKeys: []K{},
		InvalidKeys: make([]string, 0),
		Items:       map[K]V{},
		DecodeErrs:  map[K]error{},
		Datas:       map[string][]byte{},
	}

	datas, err := this.client.HGetAll(ctx, this.hashKey).Result()
	if err != nil {
		return ret, err
	}

	for strKey, data := range datas {
		ret.Datas[strKey] = []byte(data)

		key, err := this.sub.stringToKey(strKey)
		if err != nil {
			// 解释 key 失败，说明数据已经无法还原
			ret.InvalidKeys = append(ret.InvalidKeys, strKey)
			continue
		}

		var value V
		if err := this.codec.Decode([]byte(data), &value); err != nil {
			ret.DecodeErrs[key] = fmt.Errorf("decode value of str-key %q fail, %v", strKey, err)
		} else {
			ret.Items[key] = value
		}
	}
	return ret, nil
}

// 获取全部映射值
func (this *s_HSet[K, V, Sub]) GetAll() (S_RetItems[K, V], error) {
	return this.GetAllCtx(this.ctx)
}

// 扫描分批获取键值对
// 参数 match 表示 key 的匹配器：
//
//	*：匹配任意长度、任意字符（包括空字符串）
//	?：匹配任意一个字符
//	[abc]：匹配 a、b、c 中的任意一个字符
//	[a-z]：匹配 a 到 z 的任意一个字符
//	\：转义下一个字符
//
// 如 "user:*" 只返回 field 以 user: 开头的项
func (this *s_HSet[K, V, Sub]) FetchCtx(ctx context.Context, match string, fun func(S_IterItem[K, V]) bool) error {
	if match == "" {
		match = "*"
	}
	// 游标，上一次遍历的位置
	var cursor uint64
	for {
		fields, newCursor, err := this.client.HScan(ctx, this.hashKey, cursor, match, 100).Result()
		if err != nil {
			return fmt.Errorf("scan fail, %v", err)
		}
		for i := 0; i < len(fields); i += 2 {
			strKey := fields[i]
			item := S_IterItem[K, V]{
				IsKeyValid: false,
				DecodeErr:  nil,
				DBKey:      strKey,
			}

			// 解释键
			key, err := this.sub.stringToKey(strKey)
			if err != nil {
				if !fun(item) {
					return nil
				}
				continue
			}
			item.IsKeyValid = true
			item.Key = key

			// 解码值
			data := []byte(fields[i+1])
			err = this.codec.Decode(data, &(item.Value))
			if err != nil {
				item.DecodeErr = fmt.Errorf("decode value of str-key %q fail, %v", strKey, err)
				if !fun(item) {
					return nil
				}
				continue
			}
			if !fun(item) {
				return nil
			}

		}
		if newCursor == 0 {
			return nil
		}

		// 更新上一次遍历的位置
		cursor = newCursor
	}
}

// 扫描分批获取键值对
// 参数 match 表示 key 的匹配器：
//
//	*：匹配任意长度、任意字符（包括空字符串）
//	?：匹配任意一个字符
//	[abc]：匹配 a、b、c 中的任意一个字符
//	[a-z]：匹配 a 到 z 的任意一个字符
//	\：转义下一个字符
//
// 如 "user:*" 只返回 field 以 user: 开头的项
func (this *s_HSet[K, V, Sub]) Fetch(match string, fun func(S_IterItem[K, V]) bool) error {
	return this.FetchCtx(this.ctx, match, fun)
}

// 遍历指定匹配规则的键
func (this *s_HSet[K, V, Sub]) FetchKeysCtx(ctx context.Context, match string, fun func(S_IterKey[K]) bool) error {
	if match == "" {
		match = "*"
	}
	var cursor uint64
	for {
		strKeys, newCursor, err := this.client.HScanNoValues(ctx, this.hashKey, cursor, match, 100).Result()
		if err != nil {
			return fmt.Errorf("scan fail, %v", err)
		}
		for _, strKey := range strKeys {
			iterKey := S_IterKey[K]{DBKey: strKey}
			key, err := this.sub.stringToKey(strKey)
			iterKey.IsValid = err == nil
			iterKey.Key = key
			if !fun(iterKey) {
				return nil
			}
		}
		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}
	return nil
}

// 遍历指定匹配规则的键
func (this *s_HSet[K, V, Sub]) FetchKeys(match string, fun func(S_IterKey[K]) bool) error {
	return this.FetchKeysCtx(this.ctx, match, fun)
}

// ---------------------------------------------------------
// 更新一组映射
func (this *s_HSet[K, V, Sub]) UpdateCtx(ctx context.Context, values map[K]V) error {
	if len(values) == 0 {
		return nil
	}
	datas := map[string]any{}
	for k, v := range values {
		strKey := this.sub.keyToString(k)
		value, err := this.codec.Encode(v)
		if err != nil {
			return fmt.Errorf("encode value of key %q fail, %v", strKey, err)
		}
		// 这里存值时，必须用字符串，否则 go-redis 将会报错
		// 本身 redis 是支持 []byte 存储的，但是 go-redis 却错误地将 []byte 判断为参数，从而造成无法存入 []byte
		// 这可以说是 go-redis 的一个 bug
		datas[strKey] = string(value)
	}
	_, err := this.client.HSet(ctx, this.hashKey, datas).Result()
	return err
}

// 更新一组映射
func (this *s_HSet[K, V, Sub]) Update(values map[K]V) error {
	return this.UpdateCtx(this.ctx, values)
}

// 设置映射值
func (this *s_HSet[K, V, Sub]) SetCtx(ctx context.Context, key K, value V) error {
	return this.UpdateCtx(ctx, map[K]V{key: value})
}

// 设置映射值
func (this *s_HSet[K, V, Sub]) Set(key K, value V) error {
	return this.UpdateCtx(this.ctx, map[K]V{key: value})
}

// 删除指定的键值
func (this *s_HSet[K, V, Sub]) DelsCtx(ctx context.Context, keys ...K) error {
	if len(keys) == 0 {
		return nil
	}
	strKeys := []string{}
	for _, key := range keys {
		strKeys = append(strKeys, this.sub.keyToString(key))
	}
	return this.client.HDel(ctx, this.hashKey, strKeys...).Err()
}

// 删除指定的键值
func (this *s_HSet[K, V, Sub]) Dels(keys ...K) error {
	return this.DelsCtx(this.ctx, keys...)
}

// 清空所有集合数据
// save 表示是否立刻存盘
func (this *s_HSet[K, V, Sub]) ClearCtx(ctx context.Context) error {
	return this.client.Del(ctx, this.hashKey).Err()
}

// 清空所有集合数据
// save 表示是否立刻存盘
func (this *s_HSet[K, V, Sub]) Clear() error {
	return this.client.Del(this.ctx, this.hashKey).Err()
}

// ---------------------------------------------------------
// 设置过期时间
func (this *s_HSet[K, V, Sub]) SetExpireCtx(ctx context.Context, du time.Duration, keys ...K) error {
	strKeys := []string{}
	for _, key := range keys {
		strKeys = append(strKeys, this.sub.keyToString(key))
	}
	return this.client.HExpire(ctx, this.hashKey, du, strKeys...).Err()
}

// 设置过期时间
func (this *s_HSet[K, V, Sub]) SetExpire(du time.Duration, keys ...K) error {
	return this.SetExpireCtx(this.ctx, du, keys...)
}
