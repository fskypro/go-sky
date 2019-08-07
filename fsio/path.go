/**
@copyright: fantasysky 2016
@brief: 实现路径相关功能
@author: fanky
@version: 1.0
@date: 2019-01-06
**/

package fsio

import "os"
import "os/exec"
import "path/filepath"

// IsPathExists 判断路径是否存在(包括文件和文件夹)
// 注意：
//	如果参数 path 为空字符串，则返回 false
func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDir 判断指定路径是否是已经存在的文件夹
// 注意：
//	如果 path 为空字符串，则返回 false，而不会认为是当前路径
func IsDirExists(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断指定的路径是否是已经存在的文件
func IsFileExists(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// CurrentDir 获取可执行程序当前路径
func CurrentDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// -------------------------------------------------------------------
// GetBinPath 获取可执行程序所在绝对路径
func GetBinPath() (string, error) {
	rbin, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	abin, err := filepath.Abs(rbin)
	if err != nil {
		return "", err
	}
	dir, _ := filepath.Split(abin)
	return dir, nil
}

// GetFullPathToBin 根据可执行程序的相对路径，获取绝对路径
// 如果传入的路径以路径分隔符（linux 为“/”，windows 为“\”）开头，则，则直接返回 path
func GetFullPathToBin(path string) (string, error) {
	if path[0] == filepath.Separator {
		return path, nil
	}
	binPath, err := GetBinPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(binPath, path), nil
}
