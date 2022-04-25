package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	adminService "go-admin/app/admin/service"
	"go-admin/common/core/sdk/api"
	"go-admin/common/utils/dateUtils"
	"strconv"
	"time"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysLoginLog struct {
	api.Api
}

// GetPage 登录日志列表
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
		e.Error(500, err.Error())
		return
	}
	list := make([]models.SysLoginLog, 0)
	var count int64
	list, count, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 登录日志通过id获取
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
		e.Error(500, err.Error())
		return
	}
	result, err := s.Get(req.Id)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.OK(result, "查询成功")
}

// Delete 登录日志删除
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
		e.Error(500, err.Error())
		return
	}
	err = s.Remove(req.Ids)
	if err != nil {
		e.Error(500, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

// Export 导出日志
func (e SysLoginLog) Export(c *gin.Context) {
	req := dto.SysLoginLogSearch{}
	s := service.SysLoginLog{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	var sysConfService = new(adminService.SysConfig)
	sysConfService.Orm = s.Orm
	sysConfService.Log = s.Log

	//最小导出数据量
	maxSize, err := strconv.Atoi(sysConfService.GetWithKeyStr("max_export_size", "1000"))
	if err != nil {
		e.Error(500, fmt.Sprintf("配置读取异常：%s", err.Error()))
		return
	}

	list := make([]models.SysLoginLog, 0)
	req.PageIndex = 1
	req.PageSize = maxSize
	list, _, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	data, _ := s.GetExcel(list)
	fileName := "loginlog_" + dateUtils.ConvertToStr(time.Now(), 3) + ".xlsx"
	e.DownloadExcel(fileName, data)
}
