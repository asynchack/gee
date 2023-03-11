package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 全局中间件，用于错误处理

func Recover() HandleFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.String(http.StatusInternalServerError, "internal server error!")
			}

		}()

		c.Next() // 后续handler执行后，回到这里，然后期间如果有panic，这里注册的defer中，会利用recover使得程序避免崩溃
	}
}

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
