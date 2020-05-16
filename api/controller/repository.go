package controller

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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
		router.Get("/", r.All)
		router.Post("/", r.Create)
		router.Delete("/", r.Delete)

		router.Post("/pull/", r.Pull)
		router.Post("/push/", r.Push)
		router.Post("/open/", r.Open)
		router.Post("/reset/", r.Reset)
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

type RepositoryResp struct {
	Repository   Repository `json:"repository"`
	RecentlyLogs []string   `json:"recently_logs"`
}

func (r RepoController) Reset(req web.Request, repository *repo.Repo) (*OperationResponse, error) {
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	if err := repository.ResetRepository(name); err != nil {
		return Error(err, "重置仓库失败")
	}

	return Success()
}

func (r RepoController) Open(req web.Request, repository *repo.Repo) (*RepositoryResp, error) {
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	gitRepo, repoConf, err := repository.OpenRepository(name)
	if err != nil {
		return nil, err
	}

	resp := RepositoryResp{
		Repository: Repository{
			Name:   repoConf.Name,
			Branch: repoConf.Branch,
			URL:    repoConf.URL,
			Type:   string(repoConf.Type),
		},
		RecentlyLogs: make([]string, 0),
	}

	ref, err := gitRepo.Head()
	if err == nil {
		cIter, _ := gitRepo.Log(&git.LogOptions{From: ref.Hash()})
		_ = cIter.ForEach(func(c *object.Commit) error {
			resp.RecentlyLogs = append(resp.RecentlyLogs, fmt.Sprintf("%s %s", c.Author, c.Message))
			if len(resp.RecentlyLogs) > 20 {
				return errors.New("only allow 20 logs")
			}

			return nil
		})
	}

	return &resp, nil
}

func (r RepoController) Create(req web.Request, repository *repo.Repo) (*OperationResponse, error) {
	var repoNew Repository
	if err := req.Unmarshal(&repoNew); err != nil {
		return Error(err, "无法解析请求")
	}

	if repoNew.Name == "" || repoNew.URL == "" {
		return Error(errors.New("仓库名称和地址必填"), "必填字段")
	}

	if err := repository.AddRepository(repo.RepositoryConf{
		Name:   repoNew.Name,
		Branch: repoNew.Branch,
		URL:    repoNew.URL,
		Type:   repo.AuthType(repoNew.Type),
	}); err != nil {
		return Error(err, err.Error())
	}

	return Success()
}

func (r RepoController) Delete(req web.Request, repository *repo.Repo) (*OperationResponse, error) {
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	if err := repository.RemoveRepository(name); err != nil {
		return Error(err, "删除仓库失败")
	}

	return Success()
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
