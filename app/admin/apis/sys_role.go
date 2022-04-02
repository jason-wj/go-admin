package apis

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"go-admin/common/core/sdk"
	"go-admin/common/middleware/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/core/sdk/api"
	_ "go-admin/common/core/sdk/pkg/response"
	"go-admin/common/global"
)

type SysRole struct {
	api.Api
}

// GetPage 角色列表数据
func (e SysRole) GetPage(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleSearch{}
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

	list := make([]models.SysRole, 0)
	var count int64

	list, count, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Role数据
func (e SysRole) Get(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, fmt.Sprintf(" %s ", err.Error()))
		return
	}

	result, err := s.Get(req.Id)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, "查询失败")
		return
	}

	e.OK(result, "查询成功")
}

// Insert 创建角色
func (e SysRole) Insert(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleControl{}
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

	// 设置创建人
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid
	if req.Status == "" {
		req.Status = "2"
	}
	cb := sdk.Runtime.GetCasbinKey(c.Request.Host)
	err = s.Insert(&req, cb)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, "创建失败")
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(500, "")
		return
	}
	e.OK(req.RoleId, "创建成功")
}

// Update 修改用户角色
func (e SysRole) Update(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}
	cb := sdk.Runtime.GetCasbinKey(c.Request.Host)

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid

	err = s.Update(&req, cb)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(500, "")
		return
	}
	e.OK(req.RoleId, "更新成功")
}

// Delete 删除用户角色
func (e SysRole) Delete(c *gin.Context) {
	s := new(service.SysRole)
	req := dto.SysRoleById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, fmt.Sprintf("删除角色 %v 失败，\r\n失败信息 %s", req.Ids, err.Error()))
		return
	}

	err = s.Remove(req.Ids)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, "")
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(500, fmt.Sprintf("删除角色 %v 失败，失败信息 %s", req.Id, err.Error()))
		return
	}
	e.OK(req.Id, fmt.Sprintf("删除角色角色 %v 状态成功！", req.Id))
}

// Update2Status 修改用户角色状态
func (e SysRole) Update2Status(c *gin.Context) {
	s := service.SysRole{}
	req := dto.UpdateStatusReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, fmt.Sprintf("更新角色状态失败，失败原因：%s ", err.Error()))
		return
	}
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid
	err = s.UpdateStatus(&req)
	if err != nil {
		e.Error(500, fmt.Sprintf("更新角色状态失败，失败原因：%s ", err.Error()))
		return
	}
	e.OK(req.RoleId, fmt.Sprintf("更新角色 %v 状态成功！", req.RoleId))
}

// Update2DataScope 更新角色数据权限
func (e SysRole) Update2DataScope(c *gin.Context) {
	s := service.SysRole{}
	req := dto.RoleDataScopeReq{}
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
	data := &models.SysRole{
		RoleId:    req.RoleId,
		DataScope: req.DataScope,
		DeptIds:   req.DeptIds,
	}
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	data.UpdateBy = uid
	err = s.UpdateDataScope(&req)
	if err != nil {
		e.Error(500, fmt.Sprintf("更新角色数据权限失败！错误详情：%s", err.Error()))
		return
	}
	e.OK(nil, "操作成功")
}
