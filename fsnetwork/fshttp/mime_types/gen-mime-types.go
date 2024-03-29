/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: mime types generator
@author: fanky
@version: 1.0
@date: 2022-09-18
**/

// 读取 nginx 的 mime.types 文件生成 mime-types 映射
// 在当前目录下，执行以下命令，生成 ../mime-types.go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var goFile = `package fshttp

import "strings"

// generated by package ./mime_types project in this folder
// date time: %s
var mime_types = map[string]string{
%s}

func GetMimeType(ext string) string {
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	return mime_types[ext]
}
`

func main() {
	file, err := os.Open("./mime.types")
	if err != nil {
		log.Fatalf("load file mime.type fail, %v", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	_, err = reader.ReadBytes('{')
	if err != nil {
		log.Fatal("error mime file type!")
	}

	buff := strings.Builder{}
	for {
		bs, err := reader.ReadBytes(';')
		if err != nil {
			break
		}
		item := strings.TrimSpace(string(bs))
		if item == "" {
			continue
		}
		segs := strings.Fields(item[:len(item)-1])
		if len(segs) < 2 {
			continue
		}
		mimeType := segs[0]
		for _, seg := range segs[1:] {
			item := fmt.Sprintf("\t%q: %q,\n", seg, mimeType)
			buff.WriteString(item)
		}
	}

	mimeTypes := fmt.Sprintf(goFile, time.Now().Format("2006-01-02 15:04:05"), buff.String())
	ioutil.WriteFile("../mime-types.go", []byte(mimeTypes), 0660)

	fmt.Println("generate successfully!")
}
