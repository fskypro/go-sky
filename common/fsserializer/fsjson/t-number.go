/**
@copyright: fantasysky 2016
@brief: 数值类型值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package fsjson

import (
	"bufio"
	"fmt"
)

// -------------------------------------------------------------------
// S_Int64
// -------------------------------------------------------------------
type S_Int64 struct {
	s_Base
	value int64
}

func NewInt64(value int64) *S_Int64 {
	return &S_Int64{
		s_Base: createBase(TInt64),
		value:  value,
	}
}

func (this *S_Int64) V() int64 {
	return this.value
}

func (this *S_Int64) ToInt() int {
	return int(this.value)
}

func (this *S_Int64) ToInt8() int8 {
	return int8(this.value)
}

func (this *S_Int64) ToInt16() int16 {
	return int16(this.value)
}

func (this *S_Int64) ToInt32() int32 {
	return int32(this.value)
}

// -----------------------------------------------
func (this *S_Int64) ToUInt() uint {
	return uint(this.value)
}

func (this *S_Int64) ToUInt8() uint8 {
	return uint8(this.value)
}

func (this *S_Int64) ToUInt16() uint16 {
	return uint16(this.value)
}

func (this *S_Int64) ToUInt32() uint32 {
	return uint32(this.value)
}

func (this *S_Int64) TUInt64() uint64 {
	return uint64(this.value)
}

// -----------------------------------------------
func (this *S_Int64) ToFloat32() float32 {
	return float32(this.value)
}

func (this *S_Int64) ToFloat64() float64 {
	return float64(this.value)
}

// --------------------------------------------------------
func (this *S_Int64) AsInt64() *S_Int64 {
	return this
}

func (this *S_Int64) AsUInt64() *S_UInt64 {
	return NewUInt64(uint64(this.value))
}

func (this *S_Int64) AsFloat64() *S_Float64 {
	return NewFloat64(float64(this.value))
}

// --------------------------------------------------------
func (this *S_Int64) WriteTo(w *bufio.Writer) (int, error) {
	return w.WriteString(this.String())
}

func (this *S_Int64) String() string {
	return fmt.Sprintf("%v", this.value)
}

func (this *S_Int64) FmtString() string {
	return fmt.Sprintf("%d", this.value)
}

// -------------------------------------------------------------------
// S_UInt64
// -------------------------------------------------------------------
type S_UInt64 struct {
	s_Base
	value uint64
}

func NewUInt64(value uint64) *S_UInt64 {
	return &S_UInt64{
		s_Base: createBase(TUInt64),
		value:  value,
	}
}

func (this *S_UInt64) V() uint64 {
	return uint64(this.value)
}

func (this *S_UInt64) ToUInt() uint {
	return uint(this.value)
}

func (this *S_UInt64) ToUInt8() uint8 {
	return uint8(this.value)
}

func (this *S_UInt64) ToUInt16() uint16 {
	return uint16(this.value)
}

func (this *S_UInt64) ToUInt32() uint32 {
	return uint32(this.value)
}

// -----------------------------------------------
func (this *S_UInt64) ToInt() int {
	return int(this.value)
}

func (this *S_UInt64) ToInt8() int8 {
	return int8(this.value)
}

func (this *S_UInt64) ToInt16() int16 {
	return int16(this.value)
}

func (this *S_UInt64) ToInt32() int32 {
	return int32(this.value)
}

func (this *S_UInt64) ToInt64() int64 {
	return int64(this.value)
}

// -----------------------------------------------
func (this *S_UInt64) ToFloat32() float32 {
	return float32(this.value)
}

func (this *S_UInt64) ToFloat64() float64 {
	return float64(this.value)
}

// --------------------------------------------------------
func (this *S_UInt64) AsUInt64() *S_UInt64 {
	return this
}

func (this *S_UInt64) AsInt64() *S_Int64 {
	return NewInt64(int64(this.value))
}

func (this *S_UInt64) AsFloat64() *S_Float64 {
	return NewFloat64(float64(this.value))
}

// --------------------------------------------------------
func (this *S_UInt64) WriteTo(w *bufio.Writer) (int, error) {
	return w.WriteString(this.String())
}

func (this *S_UInt64) String() string {
	return fmt.Sprintf("%v", this.value)
}

func (this *S_UInt64) FmtString() string {
	return fmt.Sprintf("%d", this.value)
}

// -------------------------------------------------------------------
// S_Float64
// -------------------------------------------------------------------
type S_Float64 struct {
	s_Base
	value float64
}

func NewFloat64(value float64) *S_Float64 {
	return &S_Float64{
		s_Base: createBase(TFloat64),
		value:  value,
	}
}

func (this *S_Float64) V() float64 {
	return this.value
}

func (this *S_Float64) ToFloat32() float32 {
	return float32(this.value)
}

// -----------------------------------------------
func (this *S_Float64) ToInt() int {
	return int(this.value)
}

func (this *S_Float64) ToInt8() int8 {
	return int8(this.value)
}

func (this *S_Float64) ToInt16() int16 {
	return int16(this.value)
}

func (this *S_Float64) ToInt32() int32 {
	return int32(this.value)
}

func (this *S_Float64) ToInt64() int64 {
	return int64(this.value)
}

// -----------------------------------------------
func (this *S_Float64) ToUInt() uint {
	return uint(this.value)
}

func (this *S_Float64) ToUInt8() uint8 {
	return uint8(this.value)
}

func (this *S_Float64) ToUInt16() uint16 {
	return uint16(this.value)
}

func (this *S_Float64) ToUInt32() uint32 {
	return uint32(this.value)
}

func (this *S_Float64) TUInt64() uint64 {
	return uint64(this.value)
}

// --------------------------------------------------------
func (this *S_Float64) AsFloat64() *S_Float64 {
	return this
}

func (this *S_Float64) AsInt64() *S_Int64 {
	return NewInt64(int64(this.value))
}

func (this *S_Float64) AsUInt() *S_UInt64 {
	return NewUInt64(uint64(this.value))
}

// --------------------------------------------------------
func (this *S_Float64) WriteTo(w *bufio.Writer) (int, error) {
	return w.WriteString(this.String())
}

func (this *S_Float64) String() string {
	return fmt.Sprintf("%v", this.value)
}

func (this *S_Float64) FmtString() string {
	return fmt.Sprintf("%f", this.value)
}
