package main

import (
	"net"

	"github.com/mylxsw/graceful"
	"github.com/mylxsw/wizard-personal/cmd/wizard"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	wizard.ListenAddr = "127.0.0.1:20101"
	wizard.Wizard(Version, GitCommit, func(listener net.Listener, gf *graceful.Graceful) {})
}
