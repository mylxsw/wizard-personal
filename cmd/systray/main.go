package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"

	"github.com/mylxsw/wizard-personal/resources/icon"

	"github.com/getlantern/systray"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/graceful"
	"github.com/mylxsw/wizard-personal/cmd/wizard"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {

	var stop *graceful.Graceful

	systray.Run(func() {
		wizard.Wizard(Version, GitCommit, func(listener net.Listener, gf *graceful.Graceful) {
			stop = gf

			systray.SetTemplateIcon(icon.Data, icon.Data)
			systray.SetTitle("W")
			systray.SetTooltip("Wizard Personal")
			openMenuItem := systray.AddMenuItem("打开浏览器", "打开应用浏览器窗口")
			aboutMenuItem := systray.AddMenuItem("关于", "关于")
			quitMenuItem := systray.AddMenuItem("退出", "退出应用")

			gf.AddShutdownHandler(func() {
				log.Error("收到退出信号，3s 后强制退出")
				time.AfterFunc(3*time.Second, systray.Quit)
			})

			go func() {
				for {
					select {
					case <-quitMenuItem.ClickedCh:
						gf.Shutdown()
						break
					case <-openMenuItem.ClickedCh:
						openBrowser(fmt.Sprintf("http://%s", listener.Addr()))
					case <-aboutMenuItem.ClickedCh:
						openBrowser("https://github.com/mylxsw/wizard-personal")
					}
				}
			}()
		})
	}, func() {
		log.Debug("应用已退出")
	})

}

func openBrowser(url string) {

	cmd := ""
	args := make([]string, 0)

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)
	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Errorf("无法打开浏览器: %v", err)
	}
}
