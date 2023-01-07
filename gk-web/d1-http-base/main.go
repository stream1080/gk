package main

import (
	"fmt"
	"net/http"

	"github.com/stream1080/gk/gk-web/d1-http-base/gkw"
)

func main() {
	r := gkw.New()

	r.GET("/index", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":8080")

}
