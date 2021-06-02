package fsenv

import "runtime"

// 换行符
var Endline = "\n"

func init() {
	switch runtime.GOOS {
	case "Linux":
		Endline = "\n"
	case "Windows":
		Endline = "\r\n"
	case "darwin":
		Endline = "\n"
	default:
		Endline = "\n"
	}
}
