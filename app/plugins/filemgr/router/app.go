package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/plugins/filemgr/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerAppRouter)
}

// registerAppRouter
func registerAppRouter(v1 *gin.RouterGroup) {
	api := apis.App{}
	r := v1.Group("/filemgr/app").Use(middleware.Auth()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.POST("/upload", api.Upload)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
		r.GET("/export", api.Export)
	}
}
