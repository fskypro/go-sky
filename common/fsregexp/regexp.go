/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: regexp utils
@author: fanky
@version: 1.0
@date: 2023-04-23
**/

package fsregexp

import "bytes"

var escapes = []byte(`[\^\$\.\|\?\*\+\(\)\[\]\{\}\\]`)

func IsEscapeChar(char byte) bool {
	return bytes.Contains(escapes, []byte{char})
}
