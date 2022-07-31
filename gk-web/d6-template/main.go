package main

import (
	"log"
	"net/http"
	"time"

	"gkw"
)

func main() {
	r := gkw.New()

	r.Use(gkw.Logger())

	r.GET("/index", func(c *gkw.Context) {
		c.HTML(http.StatusOK, "<h1>Index Pages</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gkw.Context) {
			c.HTML(http.StatusOK, "<h1>Hello V1 Pages</h1>")
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
	v3 := r.Group("/v3")
	v3.Static("/assets", "./static")

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
