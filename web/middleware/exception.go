package middleware

import (
	"runtime/debug"

	"gitee.com/rachel_os/fastsearch/web"

	"github.com/gin-gonic/gin"
)

// Exception 处理异常
func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				c.JSON(200, web.Error(err.(error).Error()))
			}
			c.Abort()
		}()
		c.Next()
	}
}
