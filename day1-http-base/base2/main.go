package main

import (
	"gee"
	"log"
	"net/http"
)

/* jason-comment
	本版gee实现了： 
		1、Engine空结构体，并实现了ServeHTTP方法，可以接收所有http的请求，并在自己的逻辑中处理！
*/
func main() {

	engine := gee.New()
	// log.Fatal(http.ListenAndServe(":9999", nil))
	log.Fatal(http.ListenAndServe(":9999", engine))
}
