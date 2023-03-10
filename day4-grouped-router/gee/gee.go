package gee

import (
	"net/http"
)

type HandleFunc func(*Context) // 标准库中要求的处理http请求的函数定义！
// 处理函数的参数为Context了

// 定义一个空的结构体，并为其实现ServeHTTP方法，从而满足了ListenAndServe函数的第一个参数要求，是一个http.Handler 接口

type Engine struct {
	router *router

	// 继承RouterGroup
	// RouterGroup *RouterGroup
	*RouterGroup // 注意！不能是上一行写法，不然该写法的，结构体嵌套，Engine的实例，无法直接访问Group（）方法了，
	// 所有分组列表
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	parent      *RouterGroup
	middlewares []HandleFunc
	engine      *Engine // 指向全局的Engine，实现得以访问——Engine.router对象的addRoute（）方法，那么group对象的GET（）方法就是对Engine.router的addRoute（）方法的封装
}

// 在某分组下，注册它子分组
/* jason-comment
当是Engine的实例，调用group（）方法时， 因为有继承，所以可以访问该方法，且，group就是指向了Engine的实例
*/
func (group *RouterGroup) Group(prefix string) *RouterGroup {

	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine, // 都指向全局的Engine
	}

	// 添加到全局Engine的groups列表中
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

// RouterGroup实现注册路由-handler的方法
func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.engine.router.addRoute("GET", group.prefix+pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.engine.router.addRoute("POST", group.prefix+pattern, handler)
}

func New() *Engine {
	engine := &Engine{router: newRouter()}

	// 定义了分组之后，Engine的RouterGroup字段，要实例一个出来，填充进去

	newGroup := &RouterGroup{engine: engine}
	engine.RouterGroup = newGroup
	// engine.groups = append(engine.groups, newGroup) 应该不能直接append，此时还是nil
	engine.groups = []*RouterGroup{newGroup} // 字面量赋值并初始化

	return engine // 初始化构造函数，是一个空的路由映射！待之后填充！
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 这里就是Engine中，自定义的处理所有http请求逻辑的总入口！

	// 总入口这里逻辑：根据本次请求，找到对应的handler，然后由对应handler进行业务处理，只不过业务函数此时接收的是*Context类型，
	// 所以：还需要对w和req进行包装，成Context，

	context := newContext(w, req)

	// 然后将context，根据路由判断（此时这里是method-path），由对应路由的handler进行业务处理！
	// 此时处理方法，应该做在Engine.router上，由它，判断context中的Method和Path进行判断，判断自己所管理的route和handler的映射中，有无匹配的route，如果有就交由
	// 对应的handler处理，若没有，就报404

	e.router.handle(context)

}

// 为了简化，添加路由，应该类似gin实现，*Engine.GET() *Engine.POST()等方法，简化传入method这个参数，本质是对addRoute的调用，固化了第一个参数
func (e *Engine) GET(pattern string, handler HandleFunc) {
	e.router.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandleFunc) {
	e.router.addRoute("POST", pattern, handler)
}

func (e *Engine) DELETE(pattern string, handler HandleFunc) {
	e.router.addRoute("DELETE", pattern, handler)
}

func (e *Engine) PUT(pattern string, handler HandleFunc) {
	e.router.addRoute("PUT", pattern, handler)
}

// 为了实现 r = New() r.Run()方法 本质是，对ListenAndServe调用的封装
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
