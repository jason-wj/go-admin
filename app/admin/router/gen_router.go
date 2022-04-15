package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis"
	"go-admin/app/admin/apis/tools"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, sysNoCheckRoleRouter, registerDBRouter, registerSysTableRouter)
}

func sysNoCheckRoleRouter(v1 *gin.RouterGroup) {
	r1 := v1.Group("")
	{
		sys := apis.System{}
		r1.GET("/captcha", sys.GenerateCaptchaHandler)
	}

	r := v1.Group("").Use(middleware.Auth())
	{
		gen := tools.Gen{}
		r.GET("/gen/preview/:tableId", gen.Preview)
		r.GET("/gen/toproject/:tableId", gen.GenCode)
		r.GET("/gen/downloadCode/:tableId", gen.DownloadCode)
		r.GET("/gen/todb/:tableId", gen.GenMenuAndApi)
		sysTable := tools.SysTable{}
		r.GET("/gen/tabletree", sysTable.GetSysTablesTree)
	}
}

func registerDBRouter(v1 *gin.RouterGroup) {
	db := v1.Group("/db").Use(middleware.Auth())
	{
		gen := tools.Gen{}
		db.GET("/tables/page", gen.GetDBTableList)
		db.GET("/columns/page", gen.GetDBColumnList)
	}
}

func registerSysTableRouter(v1 *gin.RouterGroup) {
	tables := v1.Group("/sys/tables")
	{
		sysTable := tools.SysTable{}
		tables.Group("").Use(middleware.Auth()).GET("/page", sysTable.GetPage)
		tablesInfo := tables.Group("/info").Use(middleware.Auth())
		{
			tablesInfo.POST("", sysTable.Insert)
			tablesInfo.PUT("", sysTable.Update)
			tablesInfo.DELETE("/:tableId", sysTable.Delete)
			tablesInfo.GET("/:tableId", sysTable.Get)
			tablesInfo.GET("", sysTable.GetSysTablesInfo)
		}
	}
}
