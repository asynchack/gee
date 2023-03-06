package gee

import (
	"fmt"
	"net/http"
)

// 定义一个空的结构体，并为其实现ServeHTTP方法，从而满足了ListenAndServe函数的第一个参数要求，是一个http.Handler 接口

type HandleFunc func(http.ResponseWriter, *http.Request) // 标准库中要求的处理http请求的函数定义！

type Engine struct {
	routers map[string]HandleFunc
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 这里就是Engine中，自定义的处理所有http请求逻辑的总入口！

	// 从routers这个map根据key（路由）取出对应的value就是handler函数，然后调用handler函数

	key := req.Method + "-" + req.URL.Path
	if handler, ok := e.routers[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "[%q]404 not found!\n", req.URL)
	}

}

// Engine应该有1个方法，addRoute用来添加路由到自己所管理的router字段中，需要接收字段有：路由的方法，路由路径的pattern（一类特征的路由），路由的处理函数
func (e *Engine) addRoute(method string, pattern string, handler HandleFunc) {
	key := method + "-" + pattern
	e.routers[key] = handler
}

// 为了简化，添加路由，应该类似gin实现，*Engine.GET() *Engine.POST()等方法，简化传入method这个参数，本质是对addRoute的调用，固化了第一个参数
func (e *Engine) GET(pattern string, handler HandleFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandleFunc) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) DELETE(pattern string, handler HandleFunc) {
	e.addRoute("DELETE", pattern, handler)
}

func (e *Engine) PUT(pattern string, handler HandleFunc) {
	e.addRoute("PUT", pattern, handler)
}

// 为了实现 r = New() r.Run()方法 本质是，对ListenAndServe调用的封装
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func New() *Engine {
	return &Engine{routers: make(map[string]HandleFunc)} // 初始化构造函数，是一个空的路由映射！待之后填充！
}
