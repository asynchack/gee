package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/* jason-comment
管理一次会话中所有上下文数据Context，及其上的操作方法！
*/

type H map[string]interface{} // 给json数据结构，起个别名，方便业务代码使用，类似gin.H{"a": a-value, "b": b-value}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int

	// 添加个字段，记录每个方法中，动态路由匹配出来的参数 比如：/p/:lang/doc 中，可能的匹配就有 lang=golang
	Params map[string]string

	// 存储此次请求，需要用到哪些handler
	handlers []HandleFunc
	// 记录此次请求，该执行哪个handler了，是handlers的索引，初始化时，应初始化为-1
	index int
}

// Next()方法，实现调用业务handler和中间件handler的地方
func (c *Context) Next() {
	c.index++

	if c.index < len(c.handlers) { // 做个索引地址保护
		c.handlers[c.index](c)
	}

}

// 提供个方法，获取url匹配后的参数
func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

// Context构造函数，每次请求来时，都要被调用，根据w和req生成新的Context
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// 下面3个函数，根据code，调用w中的方法，写入本次响应的code状态码；
// 根据不同的数据格式，写入不同的响应头中的content-type字段；
// 根据不同的数据格式，调用w的方法，写入不同的数据，返回客户端！
func (c *Context) JSON(code int, data interface{}) {
	c.Status(code)
	c.SetHeader("Context-Type", "application/json")

	// ？？？
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Context-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) String(code int, format string, data ...interface{}) {
	c.SetHeader("Context-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, data...)))
}

// 给的StatusCode字段，设置code状态码数字
func (c *Context) Status(code int) {
	c.StatusCode = code
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}
