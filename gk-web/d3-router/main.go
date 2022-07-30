package main

import (
	"net/http"

	"gkw"
)

func main() {
	r := gkw.New()

	r.GET("/index", func(c *gkw.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gkw.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gkw.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gkw.Context) {
		c.JSON(http.StatusOK, gkw.H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *gkw.Context) {
		c.JSON(http.StatusOK, gkw.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":8080")
}
