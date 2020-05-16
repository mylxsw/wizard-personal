package controller

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/wizard-personal/internal/repo"
	"github.com/mylxsw/wizard-personal/pkg/dtree"
	"github.com/pkg/errors"
)

type TreeController struct {
	cc container.Container
}

func NewTreeController(cc container.Container) web.Controller {
	return &TreeController{cc: cc}
}

func (t *TreeController) Register(router *web.Router) {
	router.Group("/tree", func(router *web.Router) {
		router.Get("files/", t.FileTree)
	})
}

func (t *TreeController) FileTree(req web.Request, repo *repo.Repo) (*dtree.VueTree, error) {
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	tree, err := repo.RepositoryTree(name, ".md")
	if err != nil {
		return nil, errors.Wrap(err, "查询目录失败")
	}

	if tree.Title == "" {
		tree.Title = "根目录"
		tree.Expand = true
	}

	return tree, nil
}
