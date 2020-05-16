package repo

import (
	"context"

	"github.com/go-git/go-billy/v5"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/wizard-personal/configs"
	"github.com/mylxsw/wizard-personal/pkg/fsadapter"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app container.Container) {
	app.MustSingleton(func(conf *configs.Config) afero.Fs {
		switch conf.Storage.Type {
		case "mem":
			return afero.NewMemMapFs()
		default:
			return afero.NewBasePathFs(afero.NewOsFs(), conf.Storage.BasePath)
		}
	})
	app.MustSingleton(func(aferoFs afero.Fs) billy.Filesystem {
		return fsadapter.New(aferoFs)
	})
	app.MustSingleton(NewRepo)
	app.MustSingleton(NewProgressWriter)
}

func (s ServiceProvider) Boot(app glacier.Glacier) {
	app.MustResolve(func(repo *Repo, fs billy.Filesystem, writer *ProgressWriter) {
		conf, err := repo.LoadConf()
		if err != nil {
			panic(errors.Wrapf(err, "加载仓库配置失败: %v", err))
		}

		for _, r := range conf.Repositories {
			r.Init(conf)
			repository, err := repo.openGitRepository(context.TODO(), fs, &r, writer)
			if err != nil {
				log.Errorf("打开仓库 %s 失败： %v", r.Name, err)
				continue
			}

			log.Debugf("inspect repository %s ...", r.Name)
			repo.Inspection(repository)
		}

		log.Debug("所有仓库加载完毕")
	})
}

type ProgressWriter struct{}

func NewProgressWriter() *ProgressWriter {
	return &ProgressWriter{}
}

func (pw ProgressWriter) Write(p []byte) (n int, err error) {
	log.Debugf("git > %s", string(p))
	return len(p), nil
}
