package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/core/sdk/api"
)

type SysOperaLog struct {
	api.Api
}

// GetPage 操作日志列表
// @Summary 操作日志列表
// @Description 获取JSON
// @Tags 操作日志
// @Param title query string false "title"
// @Param method query string false "method"
// @Param requestMethod  query string false "requestMethod"
// @Param operUrl query string false "operUrl"
// @Param operIp query string false "operIp"
// @Param status query string false "status"
// @Param beginTime query string false "beginTime"
// @Param endTime query string false "endTime"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/sys-opera-log [get]
// @Security Bearer
func (e SysOperaLog) GetPage(c *gin.Context) {
	s := service.SysOperaLog{}
	req := new(dto.SysOperaLogQueryReq)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	list := make([]models.SysOperaLog, 0)
	var count int64

	list, count, err = s.GetPage(req)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 操作日志通过id获取
// @Summary 操作日志通过id获取
// @Description 获取JSON
// @Tags 操作日志
// @Param id path string false "id"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/sys-opera-log/{id} [get]
// @Security Bearer
func (e SysOperaLog) Get(c *gin.Context) {
	s := new(service.SysOperaLog)
	req := dto.SysOperaLogGetReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
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

// Delete 操作日志删除
// DeleteSysMenu 操作日志删除
// @Summary 删除操作日志
// @Description 删除数据
// @Tags 操作日志
// @Param data body authdto.SysOperaLogById true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/sys-opera-log [delete]
// @Security Bearer
func (e SysOperaLog) Delete(c *gin.Context) {
	s := new(service.SysOperaLog)
	req := dto.SysOperaLogDeleteReq{}
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

	err = s.Remove(req.Ids)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf("删除失败！错误详情：%s", err.Error()))
		return
	}
	e.OK(req.Ids, "删除成功")
}
