package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.Default() // 内部使用Use（）注册了2个全局中间件handler
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "hello shuaiwangsoserious!")
	})

	r.GET("/panic", func(c *gee.Context) {
		names := []string{"wang"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":19999")
}
