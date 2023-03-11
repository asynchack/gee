为什么要有分组？
	1、方便管理路由
	2、不同组路由应用不同的中间件（方便实现分组管理，比如有的组需要鉴权中间件，有的不需要！）

分组group对象，需要那些些属性？又需要那些能力（方法）？
	1、首先要有路由前缀，prefix
	2、然后要支持分组嵌套，那么得有parent，指向自己的父分组
	3、中间件应该是注册在分组上的，所以得有[]HandlerFunc存储，该分组所有的注册的中间件函数
	4、那需要访问访问Engine.router的能力，那么就得有的*engine 指向全局的Engine，从而获得间接访问router的能力！


业务代码的调用方式是什么样？根据调用方式推测，框架代码需要如何组织？
	1、如如下调用方式推出，全局Engine的实例r，需要有Group方法，得到一个group对象，groupd对象，需要有GET注册路由-和handler的能力
	2、此外r本身还可以继续不分组的注册路由和-handler 
	3、前面：GET、POST等方法都是调用router上的addRoute（）方法，进行的路由注册，
		r本身有router字段，可以间接使用router的addRoute（）方法 
		那么group对象，应该也指向router，这里选择指向Engine实例r，从而间接得到访问router的能力

	4、Group方法，做在Engine上吗？
		要支持分组嵌套，那么v1显然需要有Group方法，所以先做在group对象上，这样可以支持无限嵌套分组，
		那么Engine如何方法Group（）方法，简单：Engine从group对象上继承一些就行了，间接就获得了Group（）方法，Engine继承后，就看成最顶层的分组对象！
		此外：Engine再加个字段groups []*group对象，从而可以得到所有的注册的分组 
r := gee.New()
v1 := r.Group(“v1”)
v1.GET("/", func(c *gee.Context) {
	c.JSON(http.StatusOK, json数据...)
})

r.GET("/login", func(c *gee.Context) {
	c.JSON(http.StatusOK, json数据...)
})



