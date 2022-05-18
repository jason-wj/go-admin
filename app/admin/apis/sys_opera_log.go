package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	adminService "go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/core/sdk/api"
	"go-admin/common/utils/dateUtils"
	"strconv"
	"time"
)

type SysOperaLog struct {
	api.Api
}

// GetPage 操作日志列表
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
		e.Error(500, err.Error())
		return
	}

	list := make([]models.SysOperaLog, 0)
	var count int64

	list, count, err = s.GetPage(req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 操作日志通过id获取
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

// Delete 操作日志删除
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
		e.Error(500, err.Error())
		return
	}

	err = s.Remove(req.Ids)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, fmt.Sprintf("删除失败！错误详情：%s", err.Error()))
		return
	}
	e.OK(req.Ids, "删除成功")
}

// Export 导出日志
func (e SysOperaLog) Export(c *gin.Context) {
	req := dto.SysOperaLogQueryReq{}
	s := service.SysOperaLog{}
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
	maxSize, err := strconv.Atoi(sysConfService.GetWithKeyStr("sys_max_export_size", "1000"))
	if err != nil {
		e.Error(500, fmt.Sprintf("配置读取异常：%s", err.Error()))
		return
	}

	list := make([]models.SysOperaLog, 0)
	req.PageIndex = 1
	req.PageSize = maxSize
	list, _, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	data, _ := s.GetExcel(list)
	fileName := "operalog_" + dateUtils.ConvertToStr(time.Now(), 3) + ".xlsx"
	e.DownloadExcel(fileName, data)
}
