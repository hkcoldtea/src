package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hkcoldtea/src/projects/web.go/util"
)

type Router struct {
	tree      map[string]*util.Node
	routerMap map[string]Handler
	handle404 Handler
}

func InitRouter() *Router {
	return &Router{
		tree:      make(map[string]*util.Node),
		routerMap: make(map[string]Handler),
	}
}

func (r *Router) addRoute(t string, url string, handler Handler) {
	parts := util.ParseUrl(url)
	key := t + "-" + url
	_, ok := r.tree[t]
	if !ok {
		r.tree[t] = &util.Node{
			Children: make(map[string]*util.Node),
		}
	}
	r.tree[t].Insert(url, parts, 0)
	r.routerMap[key] = handler
}

func (r *Router) getRoute(t string, url string) (string, map[string]string) {
	inputParts := util.ParseUrl(url)
	pathParams := make(map[string]string)
	root, ok := r.tree[t]
	if !ok {
		return "", nil
	}
	methodUrl := root.Search(inputParts, 0)
	if methodUrl == "" {
		return "", nil
	}
	parts := util.ParseUrl(methodUrl)
	for i, v := range parts {
		if v[0] == ':' {
			pathParams[v[1 : ]] = inputParts[i]
		} else {
			if v[0] == '*' {
				var sl string = ""
				var inputP string = ""
				for i2, v2 := range inputParts[i:] {
					inputP += sl
					inputP += v2
				//	sl = "/v2" + v2 + "v2/"
					sl = "/"
					_ = i2
					_ = v2
				}
				pathParams[v[1 : ]] = inputP
			//	pathParams[v[1 : ]] = inputParts[i]
			}
		}
	}
	return methodUrl, pathParams
}
/*
func (r *Router) default404Handle(c *Context) {
	if err := r.handle404; err != nil {
		r.handle404(c)
		return
	}
	c.String(http.StatusNotFound, "404 NOT FOUND: %s.\n", c.Path)
}
*/
func (r *Router) handle(c *Context) {
	url, pathParams := r.getRoute(c.Method, c.Path)
	if url == "" {
		/*
		r.default404Handle(c)
		return
		*/
		if err := r.handle404; err != nil {
			r.handle404(c)
			return
		}
		c.String(http.StatusNotFound, "404 NOT FOUND: %s.\n", c.Path)
	//	c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound)+": %s.\n", c.Path)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			log.Printf("%s\n", trace(message))
		//	c.Fail(http.StatusInternalServerError, "Internal Server Error")
			c.Fail(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}()
	c.PathParams = pathParams
	key := c.Method + "-" + url
	r.routerMap[key](c)
}
