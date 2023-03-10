package gee

import (
	"log"
	"time"
)

// 全局的Loggger中间件注册函数 ，功能：记录-业务handler所需时间（可能包含它之后执行的其他中间件执行时长）
func Logger() HandleFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next() // 进行后续中间件、或业务handler的执行

		log.Printf("[%d] %s in %v\n", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
