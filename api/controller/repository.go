package controller

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/wizard-personal/internal/repo"
	"github.com/pkg/errors"
)

type RepoController struct {
	cc container.Container
}

func NewRepoController(cc container.Container) web.Controller {
	return &RepoController{cc: cc}
}

func (r RepoController) Register(router *web.Router) {
	router.Group("/repo", func(router *web.Router) {
		router.Post("pull/", r.Pull)
		router.Post("push/", r.Push)
		router.Get("/", r.All)
	})
}

func (r RepoController) Push(req web.Request, repository *repo.Repo) (*OperationResponse, error) {
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	if err := repository.Push(name); err != nil {
		return Error(err, "推送失败")
	}

	return Success()
}

func (r RepoController) Pull(req web.Request, repository *repo.Repo) (*OperationResponse, error) {
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	if err := repository.Pull(name); err != nil {
		return Error(err, "拉取更新失败")
	}

	return Success()
}

type Repository struct {
	Name   string `json:"name"`
	Branch string `json:"branch"`
	URL    string `json:"url"`
	Type   string `json:"type"`
}

func (r RepoController) All(req web.Request, repository *repo.Repo) ([]Repository, error) {
	conf, err := repository.LoadConf()
	if err != nil {
		return nil, err
	}

	repositories := make([]Repository, 0)
	for _, rep := range conf.Repositories {
		repositories = append(repositories, Repository{
			Name:   rep.Name,
			Branch: rep.Branch,
			URL:    rep.URL,
			Type:   string(rep.Type),
		})
	}

	return repositories, nil
}
