package middleware

import (
	"github.com/gin-gonic/gin"
	"yannotes.cn/apiserver_demos/demo10/handler"
	"yannotes.cn/apiserver_demos/demo10/pkg/errno"
	"yannotes.cn/apiserver_demos/demo10/pkg/token"
)

//token 校验中间件
//TODO:该方法存在重大漏洞，自己实现JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
