package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/stream1080/gk/gk-web/d6-template/gkw"
)

func main() {
	r := gkw.New()

	r.Use(gkw.Logger())

	r.GET("/index", func(c *gkw.Context) {
		c.HTML(http.StatusOK, "<h1>Index Pages</h1>", nil)
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gkw.Context) {
			c.HTML(http.StatusOK, "<h1>Hello V1 Pages</h1>", nil)
		})

		v1.GET("/hello", func(c *gkw.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gkw.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/assets/*filepath", func(c *gkw.Context) {
			c.JSON(http.StatusOK, gkw.H{"filepath": c.Param("filepath")})
		})

		v2.POST("/login", func(c *gkw.Context) {
			c.JSON(http.StatusOK, gkw.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	// 静态文件访问

	stu1 := &student{Name: "Pony", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	r.Static("/assets", "./static")
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gkw.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.GET("/students", func(c *gkw.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gkw.H{
			"title":  "gkw",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *gkw.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gkw.H{
			"title": "gkw",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":8080")
}

func onlyForV2() gkw.HandlerFunc {
	return func(c *gkw.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
