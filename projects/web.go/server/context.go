package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Content map[string]interface{}

type Context struct {
	// origin objects
	Writer     http.ResponseWriter
	Req        *http.Request

	// request info
	Path       string
	Method     string
	PathParams map[string]string

	// response info
	StatusCode int

	// server pointer
	server     *Server
}

func InitContext(w http.ResponseWriter, req *http.Request, s *Server) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		server: s,
	}
}

func (c *Context) GetPostValue(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) GetAttribute(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
/*
func NoContent(w http.ResponseWriter, r *http.Request) {
	// Set up any headers you want here.
	w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}
*/
func (c *Context) Fail(code int, err string) {
	c.JSON(code, Content{"message": err})
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML template render
func (c *Context) HTML(code int, uri string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)

	if err := c.server.HTMLTemplates.ExecuteTemplate(c.Writer, uri, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

func (c *Context) Redirect(uri string) {
	http.Redirect(c.Writer, c.Req, uri, http.StatusSeeOther)
}

func (c *Context) MethodNotAllowed() {
	http.Error(c.Writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (c *Context) ParseFiles(code int, files []string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)

	ts := template.Must(template.ParseFiles(files...))
	var err error
	err = ts.Execute(c.Writer, data)
	if err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

// Get returns the value for the given key
func (c *Context) Param(key string) string {
	value, _ := c.PathParams[key]
	return value
}

func (c *Context) Shutdown() {
	c.server.Shutdown()
}
