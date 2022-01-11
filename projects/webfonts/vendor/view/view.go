package view

import (
	"embed"
	"html/template"

	"github.com/hkcoldtea/src/projects/web.go/server"
)

// Res holds our static web server content.
//go:embed templates/*.tmpl
var Res embed.FS

func LoadTemplate(s *server.Server) *server.Server {
	var err error
	s.HTMLTemplates, err = template.New("goHome.tmpl").Funcs(TemplateFuncs).ParseFS(Res, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}
	return s
}
