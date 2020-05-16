package main

import (
	"fmt"
	"net"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/graceful"
	"github.com/mylxsw/wizard-personal/cmd/wizard"
)

var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

var w *astilectron.Window

func main() {

	wizard.Wizard(Version, GitCommit, func(listener net.Listener, gf *graceful.Graceful) {
		defer gf.Shutdown()
		if err := bootstrap.Run(bootstrap.Options{
			Asset:    Asset,
			AssetDir: AssetDir,
			AstilectronOptions: astilectron.Options{
				AppName:            AppName,
				AppIconDefaultPath: "resources/wizard-dark.png",
				AppIconDarwinPath: "resources/wizard.icns",
				SingleInstance:     true,
				VersionAstilectron: VersionAstilectron,
				VersionElectron:    VersionElectron,
			},
			Debug:  true,
			Logger: ElectronLogger{},
			OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
				w = ws[0]
				go func() {
					time.Sleep(5 * time.Second)
					if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
						log.Errorf("sending check.out.menu event failed: %w", err)
					}
				}()
				return nil
			},
			RestoreAssets: RestoreAssets,
			Windows: []*bootstrap.Window{{
				Homepage:       fmt.Sprintf("http://%s", listener.Addr()),
				MessageHandler: handleMessages,
				Options: &astilectron.WindowOptions{
					BackgroundColor: astikit.StrPtr("#333"),
					Center:          astikit.BoolPtr(true),
					Height:          astikit.IntPtr(760),
					Width:           astikit.IntPtr(1280),
				},
			}},
		}); err != nil {
			panic(err)
		}
	})
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	log.Debugf("received message, name=%s, payload=%v", m.Name, m.Payload)
	return
}

type ElectronLogger struct{}

func (e ElectronLogger) Print(v ...interface{}) {
	log.Debugf("%v", v)
}

func (e ElectronLogger) Printf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}
