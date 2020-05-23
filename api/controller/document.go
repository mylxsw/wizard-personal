package controller

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/wizard-personal/internal/repo"
	"github.com/pkg/errors"
)

type DocumentController struct {
	cc container.Container
}

func (d DocumentController) Register(router *web.Router) {
	router.Group("/document", func(router *web.Router) {
		router.Get("/", d.View)
		router.Get("/assets/", d.Assets)
		router.Post("/", d.Save)
		router.Delete("/", d.Delete)
	})
}

func NewDocumentController(cc container.Container) web.Controller {
	return &DocumentController{cc: cc}
}

type Document struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

func (d DocumentController) Delete(req web.Request, repository *repo.Repo) (*OperationResponse, error) {
	title := req.Input("title")
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	filename := filepath.ToSlash(title + ".md")
	if err := repository.DeleteFile(name, filename); err != nil {
		return Error(err, "删除文档失败")
	}

	return Success()
}

func (d DocumentController) Save(req web.Request, repository *repo.Repo) (*OperationResponse, error) {
	title := req.Input("title")
	originalTitle := req.Input("original_title")
	content := req.Input("content")

	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	filename := filepath.ToSlash(title + ".md")
	var originalFilename string
	if originalTitle != "" {
		originalFilename = filepath.ToSlash(originalTitle + ".md")
	}

	if err := repository.SaveFile(name, filename, originalFilename, []byte(content)); err != nil {
		return Error(err, "保存文档失败")
	}

	return Success()
}

func (d DocumentController) Assets(req web.Request, webCtx web.Context, repository *repo.Repo) (web.Response, error) {
	filename := req.Input("filename")
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	data, err := repository.GetFile(name, filename)
	if err != nil {
		return nil, err
	}

	webCtx.Response().SetContent(data)
	webCtx.Response().SetCode(http.StatusOK)
	return webCtx.NewRawResponse(), nil
}

func (d DocumentController) View(req web.Request, repository *repo.Repo) (*Document, error) {
	filename := strings.TrimSuffix(req.Input("filename"), ".md") + ".md"
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	data, err := repository.GetFile(name, filename)
	if err != nil {
		return nil, errors.Wrap(err, "读取文件失败")
	}

	return &Document{
		Content: string(data),
		Title:   strings.TrimSuffix(filepath.Base(filename), ".md"),
	}, nil
}
