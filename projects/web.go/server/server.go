package server

import (
	"context"
	"html/template"
	"net/http"
	"os"
)

type Handler func(*Context)

type Server struct {
	*Group
	router        *Router
	funcMap       template.FuncMap   // for custom render function
	HTMLTemplates *template.Template // for html render
	srv           *http.Server
}

func InitServer() *Server {
	server := &Server{
		router: InitRouter(),
	}
	server.Group = &Group{server: server}
	return server
}
/*
func (s *Server) Head(url string, handler Handler) {
	s.router.addRoute("HEAD", url, handler)
}
*/
func (s *Server) Get(url string, handler Handler) {
	s.router.addRoute("GET", url, handler)
}

func (s *Server) Post(url string, handler Handler) {
	s.router.addRoute("POST", url, handler)
}
/*
func (s *Server) Delete(url string, handler Handler) {
	s.router.addRoute("DELETE", url, handler)
}
*/
// for custom render function
func (s *Server) SetFuncMap(funcMap template.FuncMap) {
	s.funcMap = funcMap
}

func (s *Server) LoadTemplate(pattern string) {
	s.HTMLTemplates = template.Must(template.New("").Funcs(s.funcMap).ParseGlob(pattern))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := InitContext(w, req, s)
	s.router.handle(context)
}

func (s *Server) Run(address string) error {
	srv := &http.Server{
		Addr: address,
		Handler: s,
	}
	s.srv = srv
	return srv.ListenAndServe()
}

func (s *Server) RunTLS(address, ca, key string) error {
	return http.ListenAndServeTLS(address, ca, key, s)
}

func (s *Server) Shutdown() {
	s.srv.Shutdown(context.Background())
	os.Exit(0)
}

func (s *Server) Set404Handle(h Handler) {
	s.router.handle404 = h
}
