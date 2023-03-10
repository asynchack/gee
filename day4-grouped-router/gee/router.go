package gee

import (
	"net/http"
	"strings"
)

type router struct {
	// 添加一个roots字段
	roots    map[string]*Trie // 每个方法，一个根“树”
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{roots: make(map[string]*Trie), handlers: make(map[string]HandleFunc)}
}

// 吧注册的pattern切成一个个端，比如/p/:lang/doc => p :lang doc
func parsePattern(pattern string) []string {
	subPatterns := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, subPattern := range subPatterns {
		if subPattern != "" {
			parts = append(parts, subPattern)

			// 如果遇到是*开头，后面就不用看了
			if subPattern[0] == '*' {
				break
			}
		}
	}

	return parts
}

// Engine应该有1个方法，addRoute用来添加路由到自己所管理的router字段中，需要接收字段有：路由的方法，路由路径的pattern（一类特征的路由），路由的处理函数
func (r *router) addRoute(method string, pattern string, handler HandleFunc) {

	// 用method和 pattern做key，记录对应handler的地址
	// 当路由匹配时，匹配到的叶子节点里有pattern，此时再根据请求的method和patter，做key，找到对应的handler！
	key := method + "-" + pattern
	r.handlers[key] = handler

	// 根据注册时的方法，注册到不同的树上，第一次注册的方法，新建立一个树
	if root, ok := r.roots[method]; !ok {
		root = &Trie{children: make([]*Trie, 0)}
		r.roots[method] = root
	}

	parts := parsePattern(pattern)

	// 这里不能直接用root.insert()应该是root作用于仅限于if块的原因
	r.roots[method].insert(pattern, parts, 0)
}

// func (r *router) handle(c *Context) {
// 	key := c.Method + "-" + c.Path
// 	if handler, ok := r.handlers[key]; ok {
// 		handler(c)
// 	} else {
// 		// c.JSON() // 调用c 的方法，进行返回，404 not found等信息
// 		c.String(200, "404 not found!")
// 	}
// }

func (r *router) getRoute(method string, path string) (*Trie, map[string]string) {
	// 请求路由时，要进行匹配，根据传入的方法、path在现有树中进行搜索，能搜索到， 返回匹配到的叶子节点，和期间处理得到的参数，（参数和值，要加在Context上了）
	//
	/* jason-comment

	 */

	/* jason-comment
	要么是nil和nil，找到了就是对应的叶子节点和其参数对
	如何找？


	找到对应method的根树，找到了，继续，否则，直接nil ，nil返回

	有对应树，继续：
		先path按照/分割，得到parts列表；
		然后利用node的search方法，进行搜索
		node := root.search(parts, 0)
		if node != nil {
			node.pattern = /p/:lang/doc parts是p golang doc，或者： /p/*filepath /p/image/1.png
			for index, part := range node.pattern 遍历
			if part[0] == ":" {
				params[part[1:]] = parts[index]
			}
			if part[0] == "*" {
				params[part[1:]] = strings.json(parts[index:], "/")
			}
		}

	*/

	if _, ok := r.roots[method]; !ok {
		return nil, nil //连对应方法树都没，直接return
	}

	parts := parsePattern(path)

	node := r.roots[method].search(parts, 0)

	if node == nil {
		return nil, nil
	}

	params := make(map[string]string)

	for index, part := range parsePattern(node.pattern) {
		if part[0] == ':' {
			params[part[1:]] = parts[index]

		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(parts[index:], "/")
			break
		}
	}

	return node, params
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node == nil {
		c.String(http.StatusNotFound, "404 not found! %q\n", c.Path)

	}

	// 匹配到路由，将解析得到的路由参数，赋值给c，Context中已经定义了该字段，用于接收
	c.Params = params

	key := c.Method + "-" + node.pattern // 注意这是是node的pattern，而不是path！！，不然就找不到对应的handler函数了！
	// key := c.Method + "-" + c.Path
	r.handlers[key](c)
}

// // Engine应该有1个方法，addRoute用来添加路由到自己所管理的router字段中，需要接收字段有：路由的方法，路由路径的pattern（一类特征的路由），路由的处理函数
// func (r *router) addRoute(method string, pattern string, handler HandleFunc) {
// 	key := method + "-" + pattern
// 	r.handlers[key] = handler
// }

// func (r *router) handle(c *Context) {
// 	key := c.Method + "-" + c.Path
// 	if handler, ok := r.handlers[key]; ok {
// 		handler(c)
// 	} else {
// 		// c.JSON() // 调用c 的方法，进行返回，404 not found等信息
// 		c.String(200, "404 not found!")
// 	}
// }
