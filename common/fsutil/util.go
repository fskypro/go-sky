/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: util
@author: fanky
@version: 1.0
@date: 2021-04-02
**/

package fsutil

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// 深拷贝对象
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return fmt.Errorf("encode src instance error: %v", err)
	}
	reader := bytes.NewReader(buf.Bytes())
	if err := gob.NewDecoder(reader).Decode(dst); err != nil {
		return fmt.Errorf("decode memory buffer error: %v", err)
	}
	return nil
}
