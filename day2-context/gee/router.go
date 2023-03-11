package gee

type router struct {
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandleFunc)}
}

// Engine应该有1个方法，addRoute用来添加路由到自己所管理的router字段中，需要接收字段有：路由的方法，路由路径的pattern（一类特征的路由），路由的处理函数
func (r *router) addRoute(method string, pattern string, handler HandleFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		// c.JSON() // 调用c 的方法，进行返回，404 not found等信息
		c.String(200, "404 not found!")
	}
}
