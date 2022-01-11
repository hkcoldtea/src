package server

import (
	"net/http"
	"path"
)

type Group struct {
	server *Server
	urlPrefix string
}

func (g *Group) SetGroup(prefix string) *Group {
	s := g.server
	newGroup := &Group{
		server: s,
		urlPrefix: prefix,
	}
	return newGroup
}

func (g *Group) addRoute(t, postUrl string, handler Handler) {
	url := g.urlPrefix + postUrl
	g.server.router.addRoute(t, url, handler)
}

func (g *Group) Head(url string, handler Handler) {
	g.addRoute("HEAD", url, handler)
}

func (g *Group) Get(url string, handler Handler) {
	g.addRoute("GET", url, handler)
}

func (g *Group) Post(url string, handler Handler) {
	g.addRoute("POST", url, handler)
}

func (g *Group) Delete(url string, handler Handler) {
	g.addRoute("DELETE", url, handler)
}

func (g *Group) StaticResource(relativePath, root string) {
	absolutePath := path.Join(g.urlPrefix, relativePath)
	fs := http.Dir(root)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	handler := func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
	url := path.Join(relativePath, "/:filepath")
	g.Get(url, handler)
}

func (g *Group) StaticResourceFile(url, root string) {
	relativePath := path.Dir(url)
	absolutePath := path.Join(g.urlPrefix, relativePath)
	fs := http.Dir(root)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	handler := func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
	g.Get(url, handler)
}
