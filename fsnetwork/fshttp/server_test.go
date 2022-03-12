package fshttp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type Service struct {
	*S_Service
}

func (this *Service) OnError(err string, r *http.Request) {
	log.Println(err)
}

func newService() *Service {
	svc := &Service{
		S_Service: NewService(),
	}
	svc.S_Service.OnError = svc.OnError
	return svc
}

func handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	fmt.Println(string(body), err)
	w.Write([]byte("ok!"))
}

func Test(t *testing.T) {
	svr := NewServer("", 8088)
	svr.SetTLS("./certs/server.crt", "./certs/server.key")
	svc := newService()
	svc.AddHandler("/", handle)

	err := svr.Listen()
	if err != nil {
		log.Fatal(err)
	}
	err = svr.ServeTLS(svc)
	if err != nil {
		fmt.Println(err)
	}
}
