package controller

import (
	"log"
	"net/http"

	"config"
	"model"
	"view"

	"github.com/hkcoldtea/src/projects/web.go/server"
)

func Run() error {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

//	model.InitialMigration()
	model.BuildFontInventory()

	s := server.InitServer()

	s.Set404Handle(methodNotAllow)
	s.SetFuncMap(view.TemplateFuncs)

//	s.LoadTemplate("vendor/view/templates/*.tmpl")
	view.LoadTemplate(s)

	s.StaticResource("/static/css", "vendor/view/static/css")
	s.StaticResourceFile("/favicon.ico", "vendor/view/static/img/")
	s.StaticResourceFile("/robots.txt", "vendor/view/static/")

	s.Get("/", get_index)
	s.Head("/", methodNotAllow)
	s.Post("/", methodNotAllow)

	s.Get("/css", cssHandler)
	s.Get("/fonts/:family/:filepath", fontsHandler)

	Url := config.Config().Url
	return s.Run(Url)
}

func get_index(c *server.Context) {
//	c.HTML(http.StatusOK, "404.base.tmpl", nil)
	c.HTML(http.StatusOK, "demo.base.tmpl", nil)
}

func BadRequest(c *server.Context) {
	c.Writer.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}

func methodNotAllow(c *server.Context) {
	c.MethodNotAllowed()
}
