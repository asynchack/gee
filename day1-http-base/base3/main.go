package main

import (
	"fmt"
	"gee"
	"net/http"
)

/* jason-comment
本版gee实现了：
	1、Engine空结构体，并实现了ServeHTTP方法，可以接收所有http的请求，并在自己的逻辑中处理！
	2、Engine结构体，添加路由字段，是个map[string]HandleFunc的类型，可以记录，注册的路由和对应处理函数的映射管理，
		注册：添加映射
		请求处理：根据路由，在map中查找对应HandleFunc调用处理函数，处理即可！
*/
func main() {

	r := gee.New()
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello")
	})
	// log.Fatal(http.ListenAndServe(":9999", nil))
	// log.Fatal(http.ListenAndServe(":9999", engine))
	r.Run(":9999")
}
