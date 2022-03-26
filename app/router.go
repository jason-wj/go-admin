// Package app
// @Description:路由汇总，如无特殊必要，勿操作本页代码
package app

import (
	adminRouter "go-admin/app/admin/router"
	appManagerRouter "go-admin/app/plugins/appmanager/router"
	contentRouter "go-admin/app/plugins/content/router"
)

//
//  AllRouter
//  @Description: 汇总各大板块接口
//  @return []func()
//
func AllRouter() []func() {
	var routers []func()
	//管理服务
	routers = append(routers, adminRouter.InitRouter)
	//内容板块
	routers = append(routers, contentRouter.InitRouter)
	//app管理
	routers = append(routers, appManagerRouter.InitRouter)
	return routers
}
