package api

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/wizard-personal/api/controller"
)

func routers(cc container.Container) func(router *web.Router, mw web.RequestMiddleware) {
	return func(router *web.Router, mw web.RequestMiddleware) {
		mws := make([]web.HandlerDecorator, 0)
		mws = append(mws, mw.AccessLog(), mw.CORS("*"))

		router.WithMiddleware(mws...).Controllers(
			"/api",
			controller.NewWelcomeController(cc),
			controller.NewTreeController(cc),
			controller.NewDocumentController(cc),
			controller.NewRepoController(cc),
			controller.NewUploadController(cc),
		)
	}
}
