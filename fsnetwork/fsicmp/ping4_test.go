package fsicmp

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPing4(t *testing.T) {
	ping := NewPing4("", "11.39.205.119", time.Second*3)
	err := ping.Lisen()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 10)
		ping.Close()
	}()
	ping.CyclePing(time.Second, func(pinfo *S_PingInfo, err error) bool {
		fmt.Println(pinfo, err)
		return true
	})
}
