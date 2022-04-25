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
	"go-admin/common/middleware/auth"
	"go-admin/common/utils/dateUtils"
	"strconv"
	"time"
)

type SysConfig struct {
	api.Api
}

// GetPage 获取配置管理列表
func (e SysConfig) GetPage(c *gin.Context) {
	s := service.SysConfig{}
	req := dto.SysConfigSearch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	list := make([]models.SysConfig, 0)
	var count int64
	list, count, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取配置管理
func (e SysConfig) Get(c *gin.Context) {
	req := dto.SysConfigById{}
	s := service.SysConfig{}
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
	result, err := s.Get(req.Id)
	if err != nil {
		e.Error(500, err.Error())
		return
	}

	e.OK(result, "查询成功")
}

// Insert 创建配置管理
func (e SysConfig) Insert(c *gin.Context) {
	s := service.SysConfig{}
	req := dto.SysConfigControl{}
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
		e.Error(500, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

// Update 修改配置管理
func (e SysConfig) Update(c *gin.Context) {
	s := service.SysConfig{}
	req := dto.SysConfigControl{}
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
	req.CurrAdminId = uid
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

// Delete 删除配置管理
func (e SysConfig) Delete(c *gin.Context) {
	s := service.SysConfig{}
	req := dto.SysConfigById{}
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

// GetSysConfigBySysApp 获取系统配置信息
func (e SysConfig) GetSysConfigBySysApp(c *gin.Context) {
	req := dto.SysConfigSearch{}
	s := service.SysConfig{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}
	// 控制只读前台的数据
	req.IsFrontend = 1
	list := make([]models.SysConfig, 0)
	list, err = s.GetWithKeyList(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	mp := make(map[string]string)
	for i := 0; i < len(list); i++ {
		key := list[i].ConfigKey
		if key != "" {
			mp[key] = list[i].ConfigValue
		}
	}
	e.OK(mp, "查询成功")
}

// Get2Set 获取配置
func (e SysConfig) Get2Set(c *gin.Context) {
	s := service.SysConfig{}
	req := make([]dto.GetSetSysConfigReq, 0)
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}
	err = s.GetForSet(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	m := make(map[string]interface{}, 0)
	for _, v := range req {
		m[v.ConfigKey] = v.ConfigValue
	}
	e.OK(m, "查询成功")
}

// Update2Set 设置配置
// @Summary 设置配置
// @Description 界面操作设置配置值
// @Tags 配置管理
func (e SysConfig) Update2Set(c *gin.Context) {
	s := service.SysConfig{}
	req := make([]dto.GetSetSysConfigReq, 0)
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

	err = s.UpdateForSet(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}

	e.OK("", "更新成功")
}

// GetSysConfigByKEYForService 根据Key获取SysConfig的Service
func (e SysConfig) GetSysConfigByKEYForService(c *gin.Context) {
	var s = new(service.SysConfig)
	var req = new(dto.SysConfigByKeyReq)
	var resp = new(dto.GetSysConfigByKEYForServiceResp)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	resp, err = s.GetWithKey(req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	e.OK(resp, s.Msg)
}

// Export 导出字典
func (e SysConfig) Export(c *gin.Context) {
	req := dto.SysConfigSearch{}
	s := service.SysConfig{}
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

	list := make([]models.SysConfig, 0)
	req.PageIndex = 1
	req.PageSize = maxSize
	list, _, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	data, _ := s.GetExcel(list)
	fileName := "config_" + dateUtils.ConvertToStr(time.Now(), 3) + ".xlsx"
	e.DownloadExcel(fileName, data)
}
