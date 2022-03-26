package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"go-admin/common/core/sdk/api"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysLoginLog struct {
	api.Api
}

// GetPage 登录日志列表
// @Summary 登录日志列表
// @Description 获取JSON
// @Tags 登录日志
// @Param username query string false "用户名"
// @Param ipaddr query string false "ip地址"
// @Param loginLocation  query string false "归属地"
// @Param status query string false "状态"
// @Param beginTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/sys-login-log [get]
// @Security Bearer
func (e SysLoginLog) GetPage(c *gin.Context) {
	s := service.SysLoginLog{}
	req := dto.SysLoginLogSearch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.SysLoginLog, 0)
	var count int64
	list, count, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 登录日志通过id获取
// @Summary 登录日志通过id获取
// @Description 获取JSON
// @Tags 登录日志
// @Param id path string false "id"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/sys-login-log/{id} [get]
// @Security Bearer
func (e SysLoginLog) Get(c *gin.Context) {
	s := service.SysLoginLog{}
	req := dto.SysLoginLogGetReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	result, err := s.Get(req.Id)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(result, "查询成功")
}

// Delete 登录日志删除
// @Summary 登录日志删除
// @Description 登录日志删除
// @Tags 登录日志
// @Param data body authdto.SysLoginLogById true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/sys-login-log [delete]
// @Security Bearer
func (e SysLoginLog) Delete(c *gin.Context) {
	s := service.SysLoginLog{}
	req := dto.SysLoginLogDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	err = s.Remove(req.Ids)
	if err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
