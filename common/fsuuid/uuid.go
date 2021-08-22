/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: implement uuid generators
@author: fanky
@version: 1.0
@date: 2021-08-19
**/

package fsuuid

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
	clockSeq = rand.Uint32()
}

// -------------------------------------------------------------------
// S_UUID
// -------------------------------------------------------------------
type S_UUID struct {
	seg1 uint32
	seg2 uint16
	seg3 uint16
	seg4 uint16
	seg5 uint64
}

func (this *S_UUID) Version() uint8 {
	return uint8((this.seg3 & 0xf000) >> 12)
}

func (this *S_UUID) Variant() uint8 {
	return uint8(this.seg4 & 0xe000 >> 13)
}

func (this *S_UUID) String() string {
	return fmt.Sprintf("%08X-%04X-%04X-%04X-%012X", this.seg1, this.seg2, this.seg3, this.seg4, this.seg5)
}

func (this *S_UUID) LowerString() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", this.seg1, this.seg2, this.seg3, this.seg4, this.seg5)
}

func (this *S_UUID) ShortString() string {
	return fmt.Sprintf("%08X%04X%04X%04X%012X", this.seg1, this.seg2, this.seg3, this.seg4, this.seg5)
}

func (this *S_UUID) LowerShortString() string {
	return fmt.Sprintf("%08x%04x%04x%04x%012x", this.seg1, this.seg2, this.seg3, this.seg4, this.seg5)
}

// -----------------------------------------------
func (this *S_UUID) setVersion(v uint8) {
	this.seg3 = (this.seg3 & 0x0fff) | (uint16(v) << 12)
}

// bits 表示 varient 占几个二进制位
func (this *S_UUID) setVarient(v uint8, bits int) {
	if bits == 2 {
		this.seg4 = (this.seg4 & 0x3fff) | (uint16(v) << 14)
	} else {
		this.seg4 = (this.seg4 & 0x1fff) | (uint16(v) << 13)
	}
}

// -------------------------------------------------------------------
// uuidv1
// -------------------------------------------------------------------
var baseTime = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC).Unix()
var clockSeq uint32

// uuid v1 版本
func NewV1() *S_UUID {
	clock := atomic.AddUint32(&clockSeq, 1)
	uuid := new(S_UUID)
	now := time.Now().UTC()
	// 1582 年到现在的纳米数除以 100，作为 UUID 的时间戳
	timeStamp := uint64(now.Unix()-baseTime)*1e7 + uint64(now.Nanosecond()/100)
	uuid.seg1 = uint32(timeStamp & 0xffffffff)       // 时间戳的低 32 位
	uuid.seg2 = uint16((timeStamp >> 32) & 0xffff)   // 时间戳的中间 16 位
	uuid.seg3 = uint16(timeStamp >> 48)              // 时间戳的高 16 位
	uuid.seg4 = uint16(clock % 0x10000)              // 时钟序列，这里不设置 variant（variant，可以在生成 UUID 对象后，在 UUID 对象中设置）
	uuid.seg5 = uint64(rand.Int63n(0x1000000000000)) // node 部分，这里使用随机数，而不是网卡 MAC 地址
	uuid.setVersion(uint8(0x1))                      // version 设置为二进制的：0001
	uuid.setVarient(uint8(0x2), 2)                   // varient 设置为 RFC4122 标准（即二进制：10X。X为任意值，所以只要设置前面两位为 10 即可）
	return uuid
}

// uuid v4 版本
func NewV4() *S_UUID {
	uuid := new(S_UUID)
	high := rand.Uint64()
	uuid.seg1 = uint32(high >> 32)
	uuid.seg2 = uint16((high >> 16) & 0xffff)
	uuid.seg3 = uint16(high & 0xffff)
	low := rand.Uint64()
	uuid.seg4 = uint16(low >> 48)
	uuid.seg5 = low & 0xffffffffffff
	uuid.setVersion(uint8(0x4))
	uuid.setVarient(uint8(0x2), 2)
	return uuid
}
