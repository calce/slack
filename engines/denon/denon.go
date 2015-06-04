// turn on and off a Denon receiver
// 
package denon

import (
	"net"
	"net/url"
	"time"
	"strconv"
)

const (
	PowerOn = "power-on"
	StandBy = "standby"
)

type Denon struct {
	host string
	port string
}

func (this *Denon) GetName() string {
	return "denon"
}

func (this *Denon) Do(action string, params url.Values) (interface{}, error) {

	cmd := ""

	switch action {
		case PowerOn: cmd = "PWON"
		case StandBy: cmd = "STANDBY"
		default: return nil, errors.New("Action needed")
	}

	return this.send(cmd + "\r")
}

func (this *Denon) send(cmd) (string, error) {

	addr, _ := net.ResolveTCPAddr("tcp", this.host + ":" + this.port)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil { return "cannot connect to denon://" + host + ":" + this.port }

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	rw.SetWriteDealine(time.Now().Add(time.Second * 1))
	rw.WriteString(cmd)

	rw.SetReadDealine(time.Now().Add(time.Second * 1))
	return rw.ReadString(13)
}