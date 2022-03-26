package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis"

	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysOperaLogRouter)
}

// 需认证的路由代码
func registerSysOperaLogRouter(v1 *gin.RouterGroup) {
	api := apis.SysOperaLog{}
	r := v1.Group("/sys-opera-log").Use(middleware.Auth()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.DELETE("", api.Delete)
	}
}
