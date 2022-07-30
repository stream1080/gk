package main

import (
	"net/http"

	"gkw"
)

func main() {
	r := gkw.New()

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

	r.Run(":8080")
}
