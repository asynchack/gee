package main

import (
	"gee"
)

/* jason-comment
	func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")

	为了实现在业务代码中，提供以上gee库的调用方式，本次需要：
	1、将router字段，以及该字段相关的操作方法独立抽离，放在单独的router.go文件中， 且方面以后router的增强（以后Engine就是个框，总入口，各类数据和其操作方法，单独管理）
	2、封装一个*gee.Context字段，类似*gin.Context 是一次http上下文所涉及的数据，主要2大块，还是w http.ResponseWriter req *http.Request，再加一些自定义字段——（本质上，一次http请求和相应，就是一堆上下文数据的处理，每一次都是一个独立的Context，然后在数据上，实现方法，即在Context上实现方法，来操作数据，把特殊操作和数据进行了绑定，面向对象？）
		2.1、Context要实现的方法，PostForm（），Query（）实现快速查询本次Context中的查询参数，表单参数等，（底层是对 req的操作，调用的方法也是req的方法！）
		2.2、JSON HTML Sting等方法，快速针对此次请求，构造响应，然后返回（底层 还是对w的操作，使用的也是w的方法）


	为什么？要实现Context？？？
		web服务，本质上就是一次请求、一次响应，所有的handler函数就是对req *http.Request处理，然后写到w http.ResponseWriter中，返回给客户端；
		那么一次请求、和其响应中涉及到的所有数据，都属于本次处理中相关的“上下文的数据，即Context“，且每个请求各不相同，（随着一个请求产生而产生！随着一个请求结束而销毁！）
		我们所写的handler本质就是，利用w 和 req提供的数据和方法，对数据进行操作，加工后，返回，但是！这2个参数的方法颗粒度太细，实际写业务代码时，重复代码较多，比如每个响应，都需要使用w.setHeader()方法，设置本次响应的头部！ 类似的还要考虑每次的Body，状态码等等
		如果有了框架，就可以把头部，状态码等，数据的设置或填充，在框架代码层面完成，——从而简化了业务代码部分！！！

		此外：一次请求涉及到的Context还不止，基础库提供的w 和req，我们可能还需要一些其他的数据，比如每次Context所需要经过的中间件！如果是解析动态路由/api/v1/:name，那么本次Context解析出来的name的值，是多少，也应该存在Context中，
			总结Context需要管理一次请求-响应中，所有！可能涉及到的数据，以及，实现所有！需要操作这些数据的方法，所谓的中间件、handler就像是一个个流水线工人（函数）逐一对一次Context进行数据加工，（方法作用在数据上！）最终实现一次完整的http请求-响应！
}
*/
func main() {

	r := gee.New()
	r.GET("/hello/:name", func(c *gee.Context) {
		c.JSON(200, gee.H{
			"path":          "hello",
			"name-value-is": c.Param("name"),
		})
	})
	r.POST("/test", func(c *gee.Context) {
		c.HTML(200, "this is test html page !\n")
	})
	// log.Fatal(http.ListenAndServe(":9999", nil))
	// log.Fatal(http.ListenAndServe(":9999", engine))
	r.Run(":9999")
}
