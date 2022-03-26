package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/plugins/content/apis"

	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerArticleRouter)
}

// registerArticleRouter
func registerArticleRouter(v1 *gin.RouterGroup) {
	api := apis.Article{}
	r := v1.Group("/article").Use(middleware.Auth()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
}
