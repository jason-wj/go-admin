package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/core/sdk/api"
	"go-admin/common/middleware/auth"
)

type SysMenu struct {
	api.Api
}

// GetPage Menu列表数据
func (e SysMenu) GetPage(c *gin.Context) {
	s := service.SysMenu{}
	req := dto.SysMenuSearch{}
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
	var list = make([]models.SysMenu, 0)
	list, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

// Get 获取菜单详情
func (e SysMenu) Get(c *gin.Context) {
	req := dto.SysMenuByGetId{}
	s := new(service.SysMenu)
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

// Insert 创建菜单
func (e SysMenu) Insert(c *gin.Context) {
	req := dto.SysMenuControl{}
	s := new(service.SysMenu)
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
	e.OK(req.MenuId, "创建成功")
}

// Update 修改菜单
func (e SysMenu) Update(c *gin.Context) {
	req := dto.SysMenuControl{}
	s := new(service.SysMenu)
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
	err = s.Update(&req)
	if err != nil {
		e.Error(500, "更新失败")
		return
	}
	e.OK(req.MenuId, "更新成功")
}

// Delete 删除菜单
func (e SysMenu) Delete(c *gin.Context) {
	control := new(dto.SysMenuById)
	s := new(service.SysMenu)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}
	err = s.Remove(control.Ids)
	if err != nil {
		e.Logger.Errorf("RemoveSysMenu error, %s", err)
		e.Error(500, "删除失败")
		return
	}
	e.OK(control.Ids, "删除成功")
}

// GetMenuRole 根据登录角色名称获取菜单列表数据（左菜单使用）
func (e SysMenu) GetMenuRole(c *gin.Context) {
	s := new(service.SysMenu)
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	result, err := s.SetMenuRole(auth.GetRoleName(c))

	if err != nil {
		e.Error(500, "查询失败")
		return
	}

	e.OK(result, "")
}

//// GetMenuIDS 获取角色对应的菜单id数组
//func (e SysMenu) GetMenuIDS(c *gin.Context) {
//	s := new(service.SysMenu)
//	r := service.SysRole{}
//	m := authdto.SysRoleByName{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&m, binding.JSON).
//		MakeService(&s.Service).
//		MakeService(&r.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500,  err.Error())
//		return
//	}
//	var data models.SysRole
//	err = r.GetWithName(&m, &data).Error
//
//	//data.RoleName = c.GetString("role")
//	//data.UpdateBy = auth.GetUserId(c)
//	//result, err := data.GetIDS(s.Orm)
//
//	if err != nil {
//		e.Logger.Errorf("GetIDS error, %s", err.Error())
//		e.Error(500,  "获取失败")
//		return
//	}
//	e.OK(result, "")
//}

// GetMenuTreeSelect 根据角色ID查询菜单下拉树结构
func (e SysMenu) GetMenuTreeSelect(c *gin.Context) {
	m := service.SysMenu{}
	r := service.SysRole{}
	req := dto.SelectRole{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&m.Service).
		MakeService(&r.Service).
		Bind(&req, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	result, err := m.SetLabel()
	if err != nil {
		e.Error(500, "查询失败")
		return
	}

	menuIds := make([]int, 0)
	if req.RoleId != 0 {
		menuIds, err = r.GetRoleMenuId(req.RoleId)
		if err != nil {
			e.Error(500, "")
			return
		}
	}
	e.OK(gin.H{
		"menus":       result,
		"checkedKeys": menuIds,
	}, "获取成功")
}
