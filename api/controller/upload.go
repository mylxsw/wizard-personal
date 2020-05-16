package controller

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/go-uuid"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/wizard-personal/internal/repo"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type UploadController struct {
	cc container.Container
}

func NewUploadController(cc container.Container) web.Controller {
	return UploadController{
		cc: cc,
	}
}

func (c UploadController) Register(router *web.Router) {
	router.Group("/upload", func(router *web.Router) {
		router.Post("/images/", c.UploadImages)
	})
}

type ImageUploadResp struct {
	URL string `json:"url"`
}

func (c UploadController) UploadImages(req web.Request, repository *repo.Repo) (*ImageUploadResp, error) {
	name := req.Input("name")
	if name == "" {
		return nil, errors.New("name is required")
	}

	file, err := req.File("image")
	if err != nil {
		return nil, errors.Wrap(err, "获取上传的文件失败")
	}

	if !funk.Contains([]string{"jpg", "jpeg", "png", "gif"}, file.Extension()) {
		return nil, errors.New("不支持该类型的图片")
	}

	data, err := ioutil.ReadFile(file.GetTempFilename())
	if err != nil {
		return nil, errors.Wrap(err, "读取上传文件内容失败")
	}

	id, _ := uuid.GenerateUUID()
	relativePath := filepath.ToSlash(filepath.Join("/assets", id+"."+file.Extension()))

	if err := repository.SaveFile(name, relativePath, "", data); err != nil {
		return nil, err
	}

	return &ImageUploadResp{
		URL: relativePath,
	}, nil
}
