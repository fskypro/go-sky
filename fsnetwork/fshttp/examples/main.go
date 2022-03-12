package main

import (
	"fmt"
	"fshttp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type Service struct {
	*fshttp.S_Service
}

func (this *Service) onError(err string, r *http.Request) {
	log.Println("Error:", err)
}

func (this *Service) handle(w http.ResponseWriter, r *http.Request) {
	root := "./webroot"
	fpath := path.Join(root, r.URL.Path)
	state, err := os.Stat(fpath)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("path %q is not exists.", r.URL.Path)))
		return
	}
	if state.IsDir() {
		files, err := ioutil.ReadDir(fpath)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("read files fail, error: %v", err.Error())))
			return
		}
		for _, file := range files {
			w.Write([]byte(file.Name()))
			w.Write([]byte("\n"))
		}
		return
	}
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("read file %q error: %v", r.URL.Path, err)))
		return
	}
	w.Write(bs)
}

func newService() *Service {
	service := &Service{
		S_Service: fshttp.NewService(),
	}
	service.S_Service.OnError = service.onError
	service.S_Service.OnPanic = service.onError
	service.AddHandler("/", service.handle)
	return service
}

func main() {
	server := fshttp.NewServer("", 443)
	server.SetTLS("./certs/fsky.pro.pem", "./certs/fsky.pro.key")
	if err := server.Listen(); err != nil {
		log.Fatal("Error:", err)
	}

	service := newService()
	if err := server.ServeTLS(service); err != nil {
		log.Println("Error:", err)
	}
}
