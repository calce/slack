package main

import (
	"github.com/calce/slack"
	"github.com/calce/slack/engines/denon"
)

func main() {
	s := slack.New("", "8080", "root", "root", "", "")
	s.Register(Denon{
		host: "192.168.1.254",
		port: "23",
	}).
	Serve()
}