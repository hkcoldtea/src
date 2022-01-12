package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hkcoldtea/src/projects/web.go/server"
	"github.com/hkcoldtea/src/projects/web.go/util"
)

type student struct {
	Name string
	Age  int
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

	s := server.InitServer()

	var funcMap template.FuncMap
	funcMap = template.FuncMap{
		"formatAsDate": util.FormatAsDate,
		"trim": strings.TrimSpace,
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	}
	s.SetFuncMap(funcMap)
	s.LoadTemplate("test/templates/*")
	s.StaticResource("/static/css", "test/static")

	stu1 := &student{Name: "Mary", Age: 10}
	stu2 := &student{Name: "Peter", Age: 11}

	s.Get("/student", func(c *server.Context) {
		c.HTML(http.StatusOK, "test.tmpl", server.Content{
			"title":    "Golang",
			"students": [2]*student{stu1, stu2},
		})
	})

	s.Get("/s", func(c *server.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", server.Content{
			"title":    "Students",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	s.Get("/", func(c *server.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", server.Content{
			"title":	"Everybody",
			"now":      time.Now(),
		})
	})

	s.Get("/css", func(c *server.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	group := s.SetGroup("/pre")
	{

		group.Get("/:path/bbb/:yan", func(c *server.Context) {
			c.JSON(http.StatusOK, server.Content{
				"t1":       c.PathParams["yan"],
				"text":     c.PathParams["path"],
				"username": "yanyibin",
				"password": "yyb",
			})
		})
	}

	s.Get("/:path/bbb/:yan", func(c *server.Context) {
		c.JSON(http.StatusOK, server.Content{
			"t1":       c.PathParams["yan"],
			"text":     c.PathParams["path"],
			"username": "yanyibin",
			"password": "yyb",
		})
	})

	s.Get("/e", func(c *server.Context) {
		i := 10
		b := 0
		a := i / b
		fmt.Println(a)
	})

	s.Get("/j", func(c *server.Context) {
		c.HTML(http.StatusOK, "no-such-file.tmpl", nil)
	})

	s.Get("/b", func(c *server.Context) {
		files := []string{
			"test/static/h.html",
		}
		c.ParseFiles(http.StatusOK, files, nil)
	})

	s.Post("/stop", func(c *server.Context) {
		c.ExecuteInOrder(PreProcess, AProcess, PostProcess)
		if c.StatusCode < 100 {
			c.Shutdown()
		}
	})

	s.Run("localhost:9999")
}

func PreProcess(c *server.Context) {
	c.SetHeader("X-Header-A", "PreProcess")
	username := c.GetPostValue("username")
	if username == "" || username != os.Getenv("USER") {
		c.String(302, "missing data")
	}
}

func AProcess(c *server.Context) {
	username := os.Getenv("USER")
	password := c.GetPostValue("password")
	if username != "" && password != "" {
		res := util.PamAuthenticate(username, password)
		if res > 0 {
			return
		}
		c.String(404, "Not found")
		return
	} else {
		c.String(302, "missing data")
	}
}

func PostProcess(c *server.Context) {
	c.SetHeader("X-Header-B", "PostProcess")
}
