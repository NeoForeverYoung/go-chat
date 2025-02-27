package router

import (
	"net/http"

	"go-chat/config"
	"go-chat/internal/apis/handler"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

// 最开始是wire_gen.go调用的这个NewRouter
// TODO 2.16 不明白这个NewRouter具体做了哪些事情
// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handler *handler.Handler, session *cache.JwtTokenStorage) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Cors(conf.Cors))

	if conf.Log.AccessLog {
		// 创建访问过滤规则
		accessFilterRule := middleware.NewAccessFilterRule()
		// 排除talk记录路由
		accessFilterRule.Exclude("/api/v1/talk/records")
		// 排除talk历史路由
		accessFilterRule.Exclude("/api/v1/talk/history")
		// 排除talk转发路由
		accessFilterRule.Exclude("/api/v1/talk/forward")
		// 排除talk发布路由
		accessFilterRule.Exclude("/api/v1/talk/publish")
		// 添加登录路由的过滤规则
		accessFilterRule.AddRule("/api/v1/auth/login", func(info *middleware.RequestInfo) {
			info.RequestBody, _ = sjson.Set(info.RequestBody, `password`, "过滤敏感字段")
		})
		// 使用访问过滤规则创建访问日志中间件
		router.Use(middleware.AccessLog(
			logger.CreateFileWriter(conf.Log.LogFilePath("access.log")),
			accessFilterRule,
		))
	}

	router.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"code": 500, "msg": "系统错误，请重试!!!"})
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]any{"code": 200, "message": "hello world"})
	})

	router.GET("/health/check", func(c *gin.Context) {
		c.JSON(200, map[string]any{"status": "ok"})
	})

	// 注册 Web 路由
	RegisterWebRoute(conf.Jwt.Secret, router, handler.Api, session)

	// 注册 Admin 路由
	RegisterAdminRoute(conf.Jwt.Secret, router, handler.Admin, session)

	// 注册 Open 路由
	RegisterOpenRoute(router, handler.Open)

	// 注册 debug 路由
	if conf.Debug() {
		RegisterDebugRoute(router)
	}

	// 注册404路由
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, map[string]any{"code": 404, "message": "请求地址不存在"})
	})

	return router
}
