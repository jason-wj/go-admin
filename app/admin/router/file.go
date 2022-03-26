package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerFileRouter)
}

// 需认证的路由代码
func registerFileRouter(v1 *gin.RouterGroup) {
	var api = apis.File{}
	r := v1.Group("").Use(middleware.Auth())
	{
		r.POST("/public/uploadFile", api.UploadFile)
	}
}
