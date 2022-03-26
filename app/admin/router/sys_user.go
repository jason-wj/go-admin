package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysUserRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerNoCheckUserRouter)
}

// 需认证的路由代码
func registerSysUserRouter(v1 *gin.RouterGroup) {
	api := apis.SysUser{}
	r := v1.Group("/sys-user").Use(middleware.Auth()).Use(middleware.AuthCheckRole()).Use(actions.PermissionAction())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("", api.Update)
		r.DELETE("", api.Delete)
	}

	user := v1.Group("/user").Use(middleware.Auth()).Use(middleware.AuthCheckRole()).Use(actions.PermissionAction())
	{
		user.GET("/profile", api.GetProfile)
		user.POST("/avatar", api.InsetAvatar)
		user.PUT("/pwd/set", api.UpdatePwd)
		user.PUT("/pwd/reset", api.ResetPwd)
		user.PUT("/status", api.UpdateStatus)
	}
	v1auth := v1.Group("").Use(middleware.Auth())
	{
		v1auth.GET("/getinfo", api.GetInfo)
	}
}

func registerNoCheckUserRouter(v1 *gin.RouterGroup) {
	api := apis.SysUser{}
	r := v1.Group("")
	{
		r.POST("/login", api.Login)
		r.POST("/logout", api.LogOut)
	}

}
