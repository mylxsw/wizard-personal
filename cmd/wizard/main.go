//go:generate go run -tags generate ../gen.go

package wizard

import (
	"fmt"
	syslog "log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/glacier/starter/application"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/wizard-personal/api"
	"github.com/mylxsw/wizard-personal/configs"
	"github.com/mylxsw/wizard-personal/internal/repo"
	"github.com/mylxsw/wizard-personal/internal/utils"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

var ListenAddr = "127.0.0.1:0"

func Wizard(version, gitCommit string, mainFunc interface{}) {
	logFormatter := formatter.NewDefaultFormatter(true)
	log.All().LogFormatter(logFormatter)

	syslog.SetOutput(utils.NewLogConverter())
	//log.All().LogWriter(writer.NewDefaultFileWriter("wizard-personal.log"))

	app := application.Create(fmt.Sprintf("%s (%s)", version, gitCommit[:8]))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "use_local_dashboard",
		Usage: "whether using local dashboard, this is used when development",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "work_dir",
		Usage: "Markdown 文档所在目录",
		Value: resolveDefaultWorkDir(),
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "storage_type",
		Usage: "底层存储方式：os/mem",
		Value: "os",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "storage_base",
		Usage: "当使用本地文件系统存储(os)时，使用该目录作为内部存储的根",
		Value: "/tmp/wizard-personal",
	}))

	gl := app.Glacier()
	gl.WithHttpServer(ListenAddr)
	gl.DefaultLogFormatter(logFormatter)
	//gl.UseStackLogger(func(cc container.Container, stackWriter *writer.StackWriter) {
	//	stackWriter.PushWithLevels(writer.NewStdoutWriter())
	//})

	gl.Singleton(func(c glacier.FlagContext) *configs.Config {
		return &configs.Config{
			Listen:            c.String("listen"),
			UseLocalDashboard: c.Bool("use_local_dashboard"),
			WorkDir:           c.String("work_dir"),
			Storage: configs.StorageConfig{
				Type:     c.String("storage_type"),
				BasePath: c.String("storage_base"),
			},
		}
	})

	gl.BeforeInitialize(func(c glacier.FlagContext) error {
		// disable logs for cron
		log.Module("glacier.cron").LogLevel(level.Warning)
		return nil
	}).WebAppExceptionHandler(func(ctx web.Context, err interface{}) web.Response {
		log.Errorf("Err: %v, Stack: %s", err, debug.Stack())
		return ctx.JSONError(fmt.Sprintf("%v", err), http.StatusInternalServerError)
	})

	gl.Main(func(conf *configs.Config, router *mux.Router) {
		log.WithFields(log.Fields{
			"config": conf,
		}).Debug("configuration")

		for _, r := range web.GetAllRoutes(router) {
			log.Debugf("route: %s -> %s | %s | %s", r.Name, r.Methods, r.PathTemplate, r.PathRegexp)
		}

		gl.MustResolve(mainFunc)
	})

	gl.Provider(api.ServiceProvider{})
	gl.Provider(repo.ServiceProvider{})

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}

func resolveDefaultWorkDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}
