/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: redis-list
@author: fanky
@version: 1.0
@date: 2025-06-03
**/

package list

import (
	"context"
	"fmt"

	"fsky.pro/fsredis/codec/jsoncodec"

	"fsky.pro/fsredis/codec"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type E_Direct string

const (
	D_LEFT  = "LEFT"  // 列表左侧
	D_RIGHT = "RIGHT" // 列表右侧
)

// ----------------------------------------------------------------------------
// 返回值对象
// ----------------------------------------------------------------------------
type S_RetValue[T any] struct {
	Exists    bool   // 值是否存在
	Value     T      // 元素值
	DecodeErr error  // 解码值对象错误
	Data      []byte // 存储数据
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------
type S_List[T any] struct {
	ctx    context.Context
	client *redis.Client
	key    string
	codec  codec.I_Codec[T]
}

func NewList[T any](client *redis.Client, key string, codec codec.I_Codec[T]) *S_List[T] {
	return NewListContext(context.Background(), client, key, codec)
}

func NewListContext[T any](ctx context.Context, client *redis.Client, key string, codec codec.I_Codec[T]) *S_List[T] {
	if codec == nil {
		codec = &jsoncodec.S_JsonCodec[T]{}
	}
	return &S_List[T]{
		ctx:    ctx,
		client: client,
		key:    key,
		codec:  codec,
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 获取数组长度
func (this *S_List[T]) Count() (int64, error) {
	return this.client.LLen(this.ctx, this.key).Result()
}

// 返回指定元素的索引
// 如果不存在，则返回第一个参数为 -1
func (this *S_List[T]) Index(value T) (int64, error) {
	bvalue, err := this.codec.Encode(value)
	if err != nil {
		return -1, fmt.Errorf("encode the after item fail, %v", err)
	}
	index, err := this.client.LPos(this.ctx, this.key, string(bvalue), redis.LPosArgs{}).Result()
	if err == redis.Nil {
		return -1, nil
	} else if err != nil {
		return -1, err
	}
	return index, nil
}

// 获取指定索引处的元素
// 如果索引不存在，则第一个返回值的 Exists 字段为 false
func (this *S_List[T]) Get(index int64) (S_RetValue[T], error) {
	ret := S_RetValue[T]{Data: []byte{}}
	data, err := this.client.LIndex(this.ctx, this.key, index).Result()
	if err == redis.Nil {
		return ret, nil
	}
	if err != nil {
		return ret, err
	}

	ret.Exists = true
	ret.Data = []byte(data)
	if err := this.codec.Decode(ret.Data, &ret.Value); err != nil {
		ret.DecodeErr = fmt.Errorf("decode value fail, %v", err)
	}
	return ret, nil
}

// 获取指定范围的元素，获取范围是 start~stop（包含 stop）
// stop 支持负索引，-1 表示最后一个元素
func (this *S_List[T]) GetRange(start, stop int64) ([]S_RetValue[T], error) {
	datas, err := this.client.LRange(this.ctx, this.key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	values := []S_RetValue[T]{}
	for _, data := range datas {
		ret := S_RetValue[T]{
			Exists: true,
			Data:   []byte(data),
		}
		ret.DecodeErr = this.codec.Decode(ret.Data, &ret.Value)
		values = append(values, ret)
	}
	return values, nil
}

// 遍历所有元素
func (this *S_List[T]) Fetch(fun func(S_RetValue[T]) bool) error {
	const gap = 100
	var index int64 = 0
	for {
		rets, err := this.GetRange(index, index+gap)
		if err != nil {
			return err
		}
		if len(rets) == 0 {
			return nil
		}
		for _, ret := range rets {
			if !fun(ret) {
				return nil
			}
		}
		index += gap + 1
	}
}

// -------------------------------------------------------------------
// 往数组左边添加元素
// 返回的第一个参数为数组的新长度
func (this *S_List[T]) LPush(values ...T) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}
	bvalues := make([]any, 0, len(values))
	for _, value := range values {
		bvalue, err := this.codec.Encode(value)
		if err != nil {
			return 0, fmt.Errorf("encode value fail, %v", err)
		}
		bvalues = append(bvalues, string(bvalue))
	}
	return this.client.LPush(this.ctx, this.key, bvalues...).Result()
}

// 往数组右边添加元素
// 返回的第一个参数为数组的新长度
func (this *S_List[T]) RPush(values ...T) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}
	bvalues := make([]any, 0, len(values))
	for _, value := range values {
		bvalue, err := this.codec.Encode(value)
		if err != nil {
			return 0, fmt.Errorf("encode value fail, %v", err)
		}
		bvalues = append(bvalues, string(bvalue))
	}
	return this.client.RPush(this.ctx, this.key, bvalues...).Result()
}

// 弹出左边第一个元素
// 如果 List 为空，则第一个返回值的 Exists 字段值为 false
func (this *S_List[T]) LPop() (S_RetValue[T], error) {
	ret := S_RetValue[T]{}
	data, err := this.client.LPop(this.ctx, this.key).Result()
	if err == redis.Nil {
		return ret, nil
	}
	if err != nil {
		return ret, err
	}
	ret.Data = []byte(data)
	ret.Exists = true
	err = this.codec.Decode(ret.Data, &ret.Value)
	if err != nil {
		ret.DecodeErr = err
	}
	return ret, nil
}

// 弹出最右边一个元素
// 如果 List 为空，则第一个返回值的 Exists 字段值为 false
func (this *S_List[T]) RPop() (S_RetValue[T], error) {
	ret := S_RetValue[T]{}
	data, err := this.client.RPop(this.ctx, this.key).Result()
	if err == redis.Nil {
		return ret, nil
	}
	if err != nil {
		return ret, err
	}
	ret.Data = []byte(data)
	ret.Exists = true
	err = this.codec.Decode(ret.Data, &ret.Value)
	if err != nil {
		ret.DecodeErr = err
	}
	return ret, nil
}

// 在指定元素前面插入一个新的元素
// 如果插入成功，则第一个返回值为新 List 的长度，否则第一个返回值为 -1
func (this *S_List[T]) InsertBefore(back T, value T) (int64, error) {
	bback, err := this.codec.Encode(back)
	if err != nil {
		return -1, fmt.Errorf("encode the after item fail, %v", err)
	}
	bvalue, err := this.codec.Encode(value)
	if err != nil {
		return -1, fmt.Errorf("encode value fail, %v", err)
	}
	return this.client.LInsertBefore(this.ctx, this.key, bback, string(bvalue)).Result()
}

// 在指定元素后面插入一个新的元素
// 如果插入成功，则第一个返回值为新 List 的长度，否则第一个返回值为 -1
func (this *S_List[T]) InsertAfter(ahead T, value T) (int64, error) {
	bahead, err := this.codec.Encode(ahead)
	if err != nil {
		return -1, fmt.Errorf("encode the after item fail, %v", err)
	}
	bvalue, err := this.codec.Encode(value)
	if err != nil {
		return -1, fmt.Errorf("encode value fail, %v", err)
	}
	return this.client.LInsertAfter(this.ctx, this.key, bahead, string(bvalue)).Result()
}

// -------------------------------------------------------------------
// 设置指定位置的元素
func (this *S_List[T]) Set(index int64, value T) error {
	bvalue, err := this.codec.Encode(value)
	if err != nil {
		return fmt.Errorf("encode value fail, %v", err)
	}
	err = this.client.LSet(this.ctx, this.key, index, string(bvalue)).Err()
	if err != nil {
		return fmt.Errorf("index %d is out of range", index)
	}
	return err
}

// 删除指定位置的值
func (this *S_List[T]) Remove(index int64) (bool, error) {
	marker := "__MARKER__" + uuid.NewString()
	script := `
		local idx = tonumber(ARGV[1])
		local marker = ARGV[2]
		local val = redis.call("LINDEX", KEYS[1], idx)
		if not val then return 0 end
		redis.call("LSET", KEYS[1], idx, marker)
		redis.call("LREM", KEYS[1], 1, marker)
		return 1
`
	res, err := this.client.Eval(this.ctx, script, []string{this.key}, index, marker).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// 删除指定值
// 如果要删除的元素不存在，则返回的第一个参数为 false
func (this *S_List[T]) RemoveValue(value T) (bool, error) {
	bvalue, err := this.codec.Encode(value)
	if err != nil {
		return false, fmt.Errorf("encode value fail, %v", err)
	}
	n, err := this.client.LRem(this.ctx, this.key, 0, string(bvalue)).Result()
	return n > 0, nil
}

// 除了指定范围内的元素保留以外，其他元素全部删除
func (this *S_List[T]) Trim(start, stop int64) error {
	return this.client.LTrim(this.ctx, this.key, start, stop).Err()
}

// -------------------------------------------------------------------
// 删除整个列表
func (this *S_List[T]) Clear() error {
	return this.client.Del(this.ctx, this.key).Err()
}

// -------------------------------------------------------------------
// 移动元素到另一个列表中
// 如果移动成功，则第一个返回值的 Exists == true
func (this *S_List[T]) MoveTo(dst *S_List[T], srcDiirect, dstDirect E_Direct) (S_RetValue[T], error) {
	ret := S_RetValue[T]{}
	data, err := this.client.LMove(this.ctx, this.key, dst.key, string(srcDiirect), string(dstDirect)).Result()
	if err == redis.Nil {
		return ret, nil
	}
	if err != nil {
		return ret, err
	}
	ret.Exists = true
	ret.Data = []byte(data)
	if err := this.codec.Decode(ret.Data, &ret.Value); err != nil {
		ret.DecodeErr = err
	}
	return ret, nil
}
