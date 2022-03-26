package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/plugins/appmanager/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerAppManagerRouter)
}

// registerAppManagerRouter
func registerAppManagerRouter(v1 *gin.RouterGroup) {
	api := apis.AppManager{}
	r := v1.Group("/app-manager").Use(middleware.Auth()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.POST("/upload", api.Upload)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
}
