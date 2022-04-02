package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/sdk/api"
	_ "go-admin/common/core/sdk/pkg/response"
	"go-admin/common/middleware/auth"
)

type SysApi struct {
	api.Api
}

// GetPage 获取接口管理列表
func (e SysApi) GetPage(c *gin.Context) {
	s := service.SysApi{}
	req := dto.SysApiQueryReq{}
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
	//数据权限检查
	p := actions.GetPermissionFromContext(c)
	list := make([]models.SysApi, 0)
	var count int64
	list, count, err = s.GetPage(&req, p)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取接口管理
func (e SysApi) Get(c *gin.Context) {
	req := dto.SysApiGetReq{}
	s := service.SysApi{}
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
	p := actions.GetPermissionFromContext(c)
	result, err := s.Get(req.Id, p)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.OK(result, "查询成功")
}

// Update 修改接口管理
// @Summary 修改接口管理
// @Description 修改接口管理
// @Tags 接口管理
func (e SysApi) Update(c *gin.Context) {
	req := dto.SysApiUpdateReq{}
	s := service.SysApi{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid
	p := actions.GetPermissionFromContext(c)
	b, err := s.Update(&req, p)
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

// DeleteSysApi 删除接口管理
func (e SysApi) DeleteSysApi(c *gin.Context) {
	req := dto.SysApiDeleteReq{}
	s := service.SysApi{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}
	p := actions.GetPermissionFromContext(c)
	err = s.Remove(req.Ids, p)
	if err != nil {
		e.Error(500, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
