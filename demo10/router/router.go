package router

import (
	"net/http"

	"yannotes.cn/apiserver_demos/demo10/handler/user"

	"yannotes.cn/apiserver_demos/demo10/handler/sd"

	"github.com/gin-gonic/gin"
	"yannotes.cn/apiserver_demos/demo10/router/middleware"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares 全局中间件
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	// 404 Handler
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// pprof router 性能测试
	//pprof.Register(g)

	//用于身份验证功能的api
	g.POST("/login", user.Login)

	u := g.Group("/v1/user")
	u.Use(middleware.JWTAuth()) //群组中间件
	{
		u.POST("/", user.Create)      //// 创建用户
		u.DELETE("/:id", user.Delete) // 删除用户
		u.PUT("/:id", user.Update)    // 更新用户
		u.GET("", user.List)          //// 用户列表
		u.GET("/:username", user.Get) // 获取指定用户的详细信息
	}

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
