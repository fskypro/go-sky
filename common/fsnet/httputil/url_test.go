package httputil

import "fmt"
import "testing"
import "fsky.pro/fstest"

func TestJoinURL(t *testing.T) {
	fstest.PrintTestBegin("JoinURL")
	fmt.Println(JoinURL("aa/", "bb", "cc/", "/dd///", "ee"))
	fmt.Println(JoinURL("http:///www.", "aa/", "bb", "cc/", "/dd///", "ee"))
	fmt.Println(JoinURL("https://fsky.pro/lovestar/", "/login/", "index.html?src=L29ubGluZXMvaW5kZXguaHRtbA=="))
	fstest.PrintTestEnd()
}
