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

	// 添加一个指向Engine实例

	engine *Engine
}

// Next()方法，实现调用业务handler和中间件handler的地方
func (c *Context) Next() {
	// // 写法1：
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ { // 依赖这里的c.index++实现，能不管中间件有无调用c.Next()方法，都能保障遍历所有的handler，包含业务handler！
		/* jason-comment
		假设1：所有中间件都没调用Next()
			那么：该for循环能确保，所有handler得到依次执行，且都是上一个完全支持完毕之后，再执行后一个，不存在递归调用 ，在第一层的Next（）中的for循环中实现所有handler的调用，只有一层next（）
		假设2：有至少一个中间件调用了Next（） ，假设是第二个中间件handler-B调用了next（）
			那么： 第二个中间件B进入递归，它的Next（）调用，进入一个新的Next函数，函数第一行的++可以保证随后执行它其后的一个handler，假设是C， 那么调用C的Next的for循环中，保证执行后面所有的handler了，
				随后函数退出，回到B，B执行Next（）后的部分，
				B退出，回到第一层的Next（），此时后续的c.index++再自增一，判断，由于已经执行了所有的handler，index已经到了最后一个，再++，也就是退出
			如果有2个中间件调用了next()呢？
				假设分别A、B、C、业务handlerD ，B和C都调用了next（）
				1、B开始，执行一半，调用第二层next，进入后首先，根据for循环，要调用C，
				2、C开始，c执行一半，遇到next（），又调了第三层的next，进入后首先，根据for循环，要调用d，此时d已经是最后一个index了，d是业务代码，没有next了一般，执行完，return
				3、业务代码return，回到第三层的next（）函数，next此时，已经无需执行后续遍历了，因为（公用的一个index，index已经指向最后了，再++即不满足条件）
				4、三层next（）return，回到c调用，c执行后半部分，return
				5、c return到第2层的next，此时next同理，for循环条件不满足，return
				6、2层next（）return，回到b，b执行完后半部分，b return
				7、b return到第1层的next（）此时next同理，for循环条件不满足，return，
				8、第1层next（）return后，就是到router的 handle了
			总结c的处理流程：
				A的全部-> B一半-> C的一半 -> d业务代码 -> C的剩下 -> B的剩下 -> 处理完毕!

			业务代码handler，是在，它紧前面一个调用了next（）方法的中间件的for循环中，被调用的
		*/
		c.handlers[c.index](c)
	}

	// c.index++
	// s := len(c.handlers)
	// for ; c.index < s; c.index++ {
	// 	c.handlers[c.index](c)
	// }

	// 写法2：该方法，无法执行到业务handler，为什么？？
	/* jason-comment
	当每个中间件都有c.Next()时，可以执行到业务handler
	当有一个中间件没有c.Next()时，就不行！
	原因是：
		写法2：默认认为，每个中间件都会调用c.Next()方法，从而在Next（）方法中通过index++，来执行后续的handler，但是！一旦有一个中间件没有调用c.Next()方法，那么该中间件后面的中间件、包含业务handler都无法被调用到！！！
			链条断了！
		所以：不能依赖中间件中调用c.Next()方法，应当确保在Next（）方法中能遍历所有handlers，因此要有for遍历，就是写法1
	*/
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

// 添加name，支持输入模版的名字
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Context-Type", "text/html")
	c.Status(code)

	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.String(http.StatusNotFound, "sorry template render error! ")
	}

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
