/**
@copyright: fantasysky 2016
@brief: 公共接口
@author: fanky
@version: 1.0
@date: 2019-05-30
**/

package freejson

import "os"
import "io"
import "io/ioutil"
import "bytes"
import "fsky.pro/fsstr/convert"

// ------------------------------------------------------------------
// 读接口
// ------------------------------------------------------------------
// Load 加载 json 文件
func Load(path string) (root *S_Object, err error) {
	fi, err := os.Open(path)
	if err != nil {
		return
	}
	defer fi.Close()
	jbytes, err := ioutil.ReadAll(fi)
	if err != nil {
		return
	}
	root, err = FromBytes(jbytes)
	if err != nil {
		err = newJsonFileError(path, err)
	}
	return
}

// Read 读取 json 流数据
func Read(r io.Reader) (root *S_Object, err error) {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return FromBytes(buf.Bytes())
}

// FromString 解释 json 字符串
func FromString(jstr string) (root *S_Object, err error) {
	jbytes := convert.String2Bytes(jstr)
	return FromBytes(jbytes)
}

// FromBytes 解释字节形式的 json 字符串
func FromBytes(jbytes []byte) (*S_Object, error) {
	return newParser(jbytes).parse()
}

// ------------------------------------------------------------------
// 写接口
// ------------------------------------------------------------------
// 写出 json 文件
// fmtInfo 为 nil 则不缩进
func Save(obj *S_Object, path string, fmtInfo *S_FmtInfo) error {
	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()
	return newWriter(fi, fmtInfo).Write(obj)
}

// 把 json 对象写入到流中
// fmtInfo 为 nil 则不缩进
func Write(w io.Writer, obj *S_Object, fmtInfo *S_FmtInfo) error {
	return newWriter(w, fmtInfo).Write(obj)
}

// 把 json 对象转换为字符串形式
// fmtInfo 为 nil 则不缩进
func ToString(obj *S_Object, fmtInfo *S_FmtInfo) (string, error) {
	bs, err := ToBytes(obj, fmtInfo)
	if err != nil {
		return "", err
	}
	return convert.Bytes2String(bs), nil
}

// 把 json 对象转换为字节数组
// fmtInfo 为 nil 则不缩进
func ToBytes(obj *S_Object, fmtInfo *S_FmtInfo) ([]byte, error) {
	var buff bytes.Buffer
	err := newWriter(&buff, fmtInfo).Write(obj)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
