/**
@copyright: fantasysky 2016
@brief: URL 解释相关工具
@author: fanky
@version: 1.0
@date: 2019-03-15
**/

package httputil

import "regexp"
import "strings"

func JoinURL(elems ...string) string {
	url := strings.Join(elems, "/")
	reg := regexp.MustCompile("//+")
	url = reg.ReplaceAllString(url, "/")
	return strings.Replace(url, ":/", "://", 1)
}
