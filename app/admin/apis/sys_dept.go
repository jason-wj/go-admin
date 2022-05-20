package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/sdk/api"
	"go-admin/common/core/sdk/pkg"
	_ "go-admin/common/core/sdk/pkg/response"
	"go-admin/common/middleware/auth"
)

type SysDept struct {
	api.Api
}

// GetPage
func (e SysDept) GetPage(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptSearch{}
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
	list := make([]models.SysDept, 0)
	list, err = s.SetDeptPage(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

// Get
func (e SysDept) Get(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptById{}
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
	p := actions.GetPermissionFromContext(c)
	result, err := s.Get(req.Id, p)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}

	e.OK(result, "查询成功")
}

// Insert 添加部门
func (e SysDept) Insert(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptControl{}
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
	req.CurrUserId = uid
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, "创建失败")
		return
	}
	e.OK(req.DeptId, "创建成功")
}

// Update
func (e SysDept) Update(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptControl{}
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
	err = s.Update(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	e.OK(req.DeptId, "更新成功")
}

// Delete
func (e SysDept) Delete(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptById{}
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

// Get2Tree 用户管理 左侧部门树
func (e SysDept) Get2Tree(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptSearch{}
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
	list := make([]dto.DeptLabel, 0)
	list, err = s.SetDeptTree(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.OK(list, "")
}

// GetDeptTreeRoleSelect TODO: 此接口需要调整不应该将list和选中放在一起
func (e SysDept) GetDeptTreeRoleSelect(c *gin.Context) {
	s := service.SysDept{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	id, err := pkg.StringToInt64(c.Param("roleId"))
	result, err := s.SetDeptLabel()
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	menuIds := make([]int64, 0)
	if id != 0 {
		menuIds, err = s.GetWithRoleId(id)
		if err != nil {
			e.Error(500, err.Error())
			return
		}
	}
	e.OK(gin.H{
		"depts":       result,
		"checkedKeys": menuIds,
	}, "")
}
