package api

import (
	"github.com/gorilla/mux"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/wizard-personal/configs"
	"net/http"
)

type ServiceProvider struct {}

func (s ServiceProvider) Register(app container.Container) {
}

func (s ServiceProvider) Boot(app glacier.Glacier) {
	app.MustResolve(func(conf *configs.Config) {
		app.WebAppRouter(routers(app.Container()))
		app.WebAppMuxRouter(func(router *mux.Router) {
			router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(FS(false)))).Name("assets")
			router.PathPrefix("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte(IndexFile()))
			})
		})
	})
}

