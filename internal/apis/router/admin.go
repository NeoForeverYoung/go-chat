package router

import (
	"go-chat/internal/apis/handler/admin"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoute 注册管理后台相关的路由
// 参数说明:
// - secret: JWT token的密钥
// - router: Gin的路由引擎实例
// - handler: 管理后台的请求处理器
// - storage: 用于存储的中间件接口
func RegisterAdminRoute(
	secret string,
	router *gin.Engine,
	handler *admin.Handler,
	storage middleware.IStorage) {

	// 创建授权中间件，用于验证请求的JWT token
	// 参数"admin"表示这是管理后台的验证
	authorize := middleware.Auth(secret, "admin", storage)

	// 创建 v1 版本的路由组，基础路径为 "/admin/v1"
	v1 := router.Group("/admin/v1")
	{
		// 首页相关路由组
		index := v1.Group("/index")
		{
			// 注册首页路由
			// GET /admin/v1/index
			index.GET("", core.HandlerFunc(handler.V1.Index.Index))
		}

		// 认证相关路由组
		// TODO 为什么admin的登录接口和web的登录接口是一样的，但还是要单独出来？
		auth := v1.Group("/auth")
		{
			// amdin登录接口,只是用于管理员登录的接口
			// POST /admin/v1/auth/login
			auth.POST("/login", core.HandlerFunc(handler.V1.Auth.Login))

			// 获取验证码接口
			// GET /admin/v1/auth/captcha
			auth.GET("/captcha", core.HandlerFunc(handler.V1.Auth.Captcha))

			// 登出接口，需要授权验证
			// GET /admin/v1/auth/logout
			auth.GET("/logout", authorize, core.HandlerFunc(handler.V1.Auth.Logout))

			// 刷新token接口，需要授权验证
			// POST /admin/v1/auth/refresh
			auth.POST("/refresh", authorize, core.HandlerFunc(handler.V1.Auth.Refresh))
		}
	}
}
