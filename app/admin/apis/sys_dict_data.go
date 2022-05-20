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
	_ "go-admin/common/core/sdk/pkg/response"
	"go-admin/common/middleware/auth"
	"go-admin/common/utils/dateUtils"
	"strconv"
	"time"
)

type SysDictData struct {
	api.Api
}

// GetPage
func (e SysDictData) GetPage(c *gin.Context) {
	s := service.SysDictData{}
	req := dto.SysDictDataSearch{}
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

	list := make([]models.SysDictData, 0)
	var count int64
	list, count, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get
func (e SysDictData) Get(c *gin.Context) {
	s := service.SysDictData{}
	req := dto.SysDictDataById{}
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
		e.Logger.Warnf("Get error: %s", err.Error())
		e.Error(500, "查询失败")
		return
	}

	e.OK(result, "查询成功")
}

// Insert
func (e SysDictData) Insert(c *gin.Context) {
	s := service.SysDictData{}
	req := dto.SysDictDataControl{}
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
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrUserId = uid
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, "创建失败")
		return
	}

	e.OK(req.Id, "创建成功")
}

// Update
func (e SysDictData) Update(c *gin.Context) {
	s := service.SysDictData{}
	req := dto.SysDictDataControl{}
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
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrUserId = uid
	b, err := s.Update(&req)
	if err != nil {
		e.Error(500, fmt.Sprintf("%s", err.Error()))
		return
	}
	if !b {
		e.OK(req.Id, "未修改任何信息")
		return
	}
	e.OK(req.Id, "修改成功")
}

// Delete
func (e SysDictData) Delete(c *gin.Context) {
	s := service.SysDictData{}
	req := dto.SysDictDataById{}
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
	e.OK(req.Id, "删除成功")
}

// GetSysDictDataAll 数据字典根据key获取 业务页面使用
func (e SysDictData) GetSysDictDataAll(c *gin.Context) {
	s := service.SysDictData{}
	req := dto.SysDictDataSearch{}
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
	list, err := s.GetAll(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

// Export 导出字典
func (e SysDictData) Export(c *gin.Context) {
	req := dto.SysDictDataSearch{}
	s := service.SysDictData{}
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

	list := make([]models.SysDictData, 0)
	req.PageIndex = 1
	req.PageSize = maxSize
	list, _, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	data, _ := s.GetExcel(list)
	fileName := "dictdata_" + dateUtils.ConvertToStr(time.Now(), 3) + ".xlsx"
	e.DownloadExcel(fileName, data)
}
