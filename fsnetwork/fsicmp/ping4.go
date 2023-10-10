/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: ping
@author: fanky
@version: 1.0
@date: 2023-09-26
**/

package fsicmp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// -------------------------------------------------------------------
// PingInfo
// -------------------------------------------------------------------
type S_PingInfo struct {
	Delay  time.Duration // 延时
	Sends  int           // 发包数
	Replys int           // 收包数
}

// 丢包数
func (this *S_PingInfo) Losts() int {
	return this.Sends - this.Replys
}

// 丢包率
func (this *S_PingInfo) LostRate() float64 {
	if this.Sends == 0 { return 0.0 }
	return float64(this.Sends-this.Replys) / float64(this.Sends)
}

// -------------------------------------------------------------------
// PingV4
// -------------------------------------------------------------------
type S_Ping4 struct {
	localIP string
	addr    string
	ip      *net.IPAddr
	msg     *icmp.Message
	conn    *icmp.PacketConn
	timeout time.Duration
	ctx     context.Context
	cancel  func()
}

func NewPing4(localIP string, addr string, timeout time.Duration) *S_Ping4 {
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  0,
			Data: []byte("ping"),
		},
	}
	if localIP == "" {
		localIP = "0.0.0.0"
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &S_Ping4{
		localIP: localIP,
		addr:    addr,
		msg:     msg,
		timeout: timeout,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (this *S_Ping4) Lisen() error {
	ip, err := net.ResolveIPAddr("ip4", this.addr)
	if err != nil {
		return fmt.Errorf("resolving IP address fial, %v", err)
	}
	this.ip = ip

	conn, err := icmp.ListenPacket("ip4:icmp", this.localIP)
	if err != nil {
		return fmt.Errorf("Error opening connection: %v", err)
	}
	this.conn = conn
	return nil
}

func (this *S_Ping4) Ping(seq int) (delay time.Duration, err error) {
	if this.conn == nil {
		err = errors.New("invalid connection, call Listen at first")
		return
	}

	this.msg.Body.(*icmp.Echo).Seq = seq
	data, e := this.msg.Marshal(nil)
	if e != nil {
		err = fmt.Errorf("marshal ICMP message fail, %v", e)
		return
	}

	start := time.Now()
	_, err = this.conn.WriteTo(data, this.ip)
	if err != nil {
		err = fmt.Errorf("sending ICMP message fail, %v", err)
		return
	}

	err = this.conn.SetReadDeadline(time.Now().Add(this.timeout))
	if err != nil {
		err = fmt.Errorf("setting read deadline, %v", err)
		return
	}

	reply := make([]byte, 1024)
	n, addr, e := this.conn.ReadFrom(reply)
	if e != nil {
		// 接受数据失败，可能超时
		err = ErrTimeout
		return
	}

	if n == 0 || addr.String() != this.ip.IP.String() {
		// 脏数据包
		err = ErrInvalidPackage
		return
	}

	delay = time.Since(start)
	return
}

// 循环 ping
func (this *S_Ping4) CyclePing(interval time.Duration, fun func(*S_PingInfo, error) bool) {
	if this.conn == nil {
		panic("invalid connection, call Listen at first")
	}
	pinInfo := &S_PingInfo{}
	for {
		select {
		case <-this.ctx.Done():
			return
		default:
			pinInfo.Sends += 1
			delay, err := this.Ping(pinInfo.Sends)
			if err == nil {
				pinInfo.Replys++
				pinInfo.Delay = (pinInfo.Delay + delay) / 2
			}
			if !fun(pinInfo, err) {
				this.Close()
				return
			}
			time.Sleep(interval)
		}
	}
}

func (this *S_Ping4) Close() error {
	this.cancel()
	if this.conn == nil {
		return nil
	}
	return this.conn.Close()
}

// -------------------------------------------------------------------
// public functions
// -------------------------------------------------------------------
func Ping4(laddr, raddr string, seq int, timeout time.Duration) (delay time.Duration, err error) {
	ip, e := net.ResolveIPAddr("ip4", raddr)
	if e != nil {
		err = fmt.Errorf("resolving IP address fial, %v", e)
		return
	}

	if laddr == "" { laddr = "0.0.0.0" }
	conn, e := icmp.ListenPacket("ip4:icmp", laddr)
	if e != nil {
		err = fmt.Errorf("Error opening connection: %v", e)
		return
	}
	defer conn.Close()

	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: []byte("ping"),
		},
	}

	data, e := msg.Marshal(nil)
	if e != nil {
		err = fmt.Errorf("marshal ICMP message fail, %v", e)
		return
	}

	start := time.Now()
	_, err = conn.WriteTo(data, ip)
	if err != nil {
		err = fmt.Errorf("sending ICMP message fail, %v", err)
		return
	}

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		err = fmt.Errorf("setting read deadline, %v", err)
		return
	}

	reply := make([]byte, 1024)
	n, addr, e := conn.ReadFrom(reply)
	if e != nil {
		// 接收数据失败，可能超时
		err = ErrTimeout
		return
	}

	if n == 0 || addr.String() != ip.IP.String() {
		// 脏数据包
		err = ErrInvalidPackage
		return
	}

	delay = time.Since(start)
	return
}
