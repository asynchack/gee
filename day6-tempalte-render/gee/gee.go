package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
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

	htmlTemplates *template.Template
	funcMap       template.FuncMap // 记录所有的funcMap
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

type RouterGroup struct {
	prefix      string
	parent      *RouterGroup
	middlewares []HandleFunc
	engine      *Engine // 指向全局的Engine，实现得以访问——Engine.router对象的addRoute（）方法，那么group对象的GET（）方法就是对Engine.router的addRoute（）方法的封装
}

func (group *RouterGroup) createStaticHandler(relatePath string, filePath string) HandleFunc {
	// 本质是借助http.FileSever实现，只需要提供它完整的url路径、和完全的文件系统路径即可

	fs := http.Dir(filePath) // 转为Dir类型，该类型实现了http.FileSystem接口，即具有Open方法

	absolutePath := path.Join(group.prefix, relatePath)              // 获得完整的url路径
	fileSever := http.StripPrefix(absolutePath, http.FileServer(fs)) // 根据完整的url路径，和完全文件系统路径，创建handler

	return func(c *Context) {
		fileName := c.Param("fileath") //

		if _, err := fs.Open(fileName); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileSever.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relatePath string, filePath string) {
	// 接收2个参数，实现url的路径，和 文件系统路径的的映射，比如 "/assets" ./static的映射

	handler := group.createStaticHandler(relatePath, filePath)
	pattern := relatePath + "/*filepath"

	group.GET(pattern, handler) // 无需考虑分组的前缀，内部自动添加前缀（注册路由时）

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

// RouterGroup添加Use（）方法 ,用于向分组对象，注册路由！
func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	for _, middleware := range middlewares {
		group.middlewares = append(group.middlewares, middleware)
	}
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

	handlers := make([]HandleFunc, 0)

	// 调用router的handle方法前，应该先根据url判断它需要应用哪些中间件！
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			handlers = append(handlers, group.middlewares...)
		}
	}

	context.handlers = handlers

	context.engine = e // 构造context实例时，让其指向全局engine！
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
