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

type SysDictType struct {
	api.Api
}

// GetPage 字典类型列表数据
func (e SysDictType) GetPage(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeSearch{}
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
	list := make([]models.SysDictType, 0)
	var count int64
	list, count, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 字典类型通过字典id获取
func (e SysDictType) Get(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeById{}
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

//Insert 字典类型创建
func (e SysDictType) Insert(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeControl{}
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
	req.CurrAdminId = uid
	err = s.Insert(&req)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, fmt.Sprintf(" 创建字典类型失败，详情：%s", err.Error()))
		return
	}
	e.OK(req.DictId, "创建成功")
}

// Update
func (e SysDictType) Update(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(500, err.Error())
		e.Logger.Error(err)
		return
	}
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid
	b, err := s.Update(&req)
	if err != nil {
		e.Error(500, fmt.Sprintf("%s", err.Error()))
		return
	}
	if !b {
		e.OK(req.DictId, "未修改任何信息")
		return
	}
	e.OK(req.DictId, "修改成功")
}

// Delete
func (e SysDictType) Delete(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeById{}
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
		e.Error(500, err.Error())
		return
	}
	e.OK(req.Id, "删除成功")
}

// GetAll
func (e SysDictType) GetAll(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeSearch{}
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
	list := make([]models.SysDictType, 0)
	list, err = s.GetAll(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	e.OK(list, "查询成功")
}

// Export 导出通告
func (e SysDictType) Export(c *gin.Context) {
	req := dto.SysDictTypeSearch{}
	s := service.SysDictType{}
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

	list := make([]models.SysDictType, 0)
	req.PageIndex = 1
	req.PageSize = maxSize
	list, _, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	data, _ := s.GetExcel(list)
	fileName := "dicttype_" + dateUtils.ConvertToStr(time.Now(), 3) + ".xlsx"
	e.DownloadExcel(fileName, data)
}
