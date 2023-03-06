package gee

import (
	"fmt"
	"net/http"
)

// 定义一个空的结构体，并为其实现ServeHTTP方法，从而满足了ListenAndServe函数的第一个参数要求，是一个http.Handler 接口

type Engine struct{}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 这里就是Engine中，自定义的处理所有http请求逻辑的总入口！

	switch req.URL.Path {
	case "/hello":
		fmt.Fprintf(w, "url.path is: [%q]\n", req.URL.Path)

	default:
		fmt.Fprintf(w, "other path, can not handle!")
	}

}

func New() *Engine {
	return &Engine{}
}
