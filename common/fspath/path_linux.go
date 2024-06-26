/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: linux 下的目录文件处理函数
@author: fanky
@version: 1.0
@date: 2022-08-28
**/

package fspath

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"fsky.pro/fserror"
)

// 获取软连接目录/文件的源目录/文件
// 允许父目录是软连接目录
func SLinkSrcPath(p string) string {
	linkPath := func(path string) string {
		src, err := filepath.EvalSymlinks(path)
		if err == nil {
			return src
		} else if fserror.IsError[*os.PathError](err) {
			re := regexp.MustCompile(`\s+.+$`)
			return re.ReplaceAllString(err.(*os.PathError).Path, "")
		} else {
			return src
		}
	}

	segs := strings.Split(path.Clean(p), string(os.PathSeparator))
	path := ""
	if segs[0] != "" {
		path = linkPath(segs[0])
	}
	for _, seg := range segs[1:] {
		path += "/" + seg
		path = linkPath(path)
	}
	return path
}
