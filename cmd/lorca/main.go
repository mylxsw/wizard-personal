package main

import (
	"fmt"
	"net"
	"net/url"

	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/graceful"
	"github.com/mylxsw/wizard-personal/cmd/wizard"
	"github.com/pkg/errors"
	"github.com/zserge/lorca"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	ui, err := lorca.New("data:text/html,"+url.PathEscape(`
<html>
		<head><title>Loading</title></head>
		<body>
			<h1>Loading ...</h1>
		</body>
	</html>
	`), "", 1280, 760)
	if err != nil {
		panic(errors.Wrap(err, "启动图形界面失败"))
	}

	defer func() { _ = ui.Close() }()

	wizard.ListenAddr = "127.0.0.1:20101"
	wizard.Wizard(Version, GitCommit, func(listener net.Listener, gf *graceful.Graceful) {
		if err := ui.Load(fmt.Sprintf("http://%s", listener.Addr())); err != nil {
			panic(errors.Wrap(err, "加载资源失败"))
		}
		if err := ui.Bind("serverIP", func() string {
			return listener.Addr().String()
		}); err != nil {
			log.Errorf("绑定 serverIP 函数失败: %v", err)
		}

		<-ui.Done()
		gf.Shutdown()
	})
}
