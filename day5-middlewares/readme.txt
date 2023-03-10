为什么要有中间件？
	1、通用类代码，具备通用的处理功能，比如对用户请求鉴权，限流，等
	2、且有全局、部分之分，全局即所有用户请求都要经过我中间件处理，比如日志、错误处理
		部分，一般是作用在分组之上，比如/v1组，的所有请求都要鉴权，认证后才能放行，比如/v2组，的所有请求需要限流，100的req/s每秒等 


中间件本质是什么？处理的是什么数据？
	本质就是和业务处理函数一样的，函数类型，HandleFunc 是一个func(c *Context)类型的函数，显然！它处理的是Context数据！


如何实现中间件？
	1、仿照gin的实现，一般是全局的Engine或group对象上，实现一个Use（）方法，来注册中间件，这里：在RouterGroup对象上实现，Engine自然继承
		1.1、Use()方法的实现逻辑：
		接收 HandleFunc... 类型的参数 
		把HandleFunc这些函数，注册到调用Use（）的group对象的middlewares列表中去 

	2、中间件调用？
		2.1、ServeHTTP方法中实现
			每个协程中，根据w 和req，组成Context一份数据实例，然后，根据本次req.url.path，和全局Engine.groups中所有的分组对象遍历匹配，进行prefix的匹配，如果匹配到了，就把这个分组对象的middlewares，都取出来，放在一个大列表中，
				（中间件handler的顺序，应该可以更精细控制，比如按照前缀匹配最长优先原则，把前缀匹配最长的分组的中间件，放在最前面！）
			最终赋值给c 实例的handlers列表上，
			然后此时交给router的handle（c）处理 

		2.2、在router的handle（c）方法中实现
			handle中，进行路由匹配的判断，如果匹配到了，那么就在c.handlers列表上，追加业务代码注册的业务handle（在最后）；没匹配到，那么就在c.handlers列表上，追加一个默认的404的处理handler（也是在最后）

		然后调用c.Next()——————开始执行，针对这个数据实例c（这一次请求），执行其所有的中间件handler，和业务处理handler函数

中间件，执行流程？业务handler前、和后的动作，在中间件中如何实现？——————c.Next()方法 
	c.Next()就是遍历c上所有的handler，进行调用，那么中间件如何实现，其中的一些动作，分别在业务处理函数前、和业务处理函数后执行呢？
	答：靠c.Next()在中间件handler中定义的位置实现；



如下：

func A(c *Context) {
    part1
    c.Next()
    part2
}
func B(c *Context) {
    part3
    c.Next()
    part4
}

再有一个业务handlerC，

当一个请求过来时，假设它匹配到的中间件有A和B两个,且A在B前，最后：还一个业务handler的C，那么该请求依次经过的处理流程就是：
part1	part3	C part4  part2  

如果有想在业务handler之后执行的操作，比如追加一些额外的响应信息，那么就一定是在中间件handler的c.Next()之后，编写代码

