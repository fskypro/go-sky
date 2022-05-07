/**
* @brief: zipper.go
* @copyright: 2016 fantasysky
* @author: fanky
* @version: 1.0
* @date: 2018-12-31
 */

package fszip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Unzip gzip 解压数据
func Unzip(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	data, err = ioutil.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return data, err
}

// Zip gzip 压缩数据
func Zip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
