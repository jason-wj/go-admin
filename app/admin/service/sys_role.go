package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"time"

	"github.com/casbin/casbin/v2"

	"go-admin/common/core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
)

type SysRole struct {
	service.Service
}

func NewSysRoleService(s *service.Service) *SysRole {
	var srv = new(SysRole)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取SysRole列表
func (e *SysRole) GetPage(c *dto.SysRoleSearch) ([]models.SysRole, int64, error) {
	var list []models.SysRole
	var data models.SysRole
	var count int64

	err := e.Orm.Model(&data).Preload("SysMenu").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysRoleService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取SysRole对象
func (e *SysRole) Get(id int64) (*models.SysRole, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysRole{}
	err := e.Orm.First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	model.MenuIds, err = e.GetRoleMenuId(model.RoleId)
	if err != nil {
		e.Log.Errorf("SysRoleService Get error:%s", err)
		return nil, err
	}
	return model, nil
}

// Insert 创建SysRole对象
func (e *SysRole) Insert(c *dto.SysRoleControl, cb *casbin.SyncedEnforcer) error {
	if c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var dataMenu []models.SysMenu
	err = e.Orm.Preload("SysApi").Where("menu_id in ?", c.MenuIds).Find(&dataMenu).Error
	if err != nil {
		e.Log.Errorf("SysRoleService Insert error:%s", err)
		return err
	}
	c.SysMenu = dataMenu
	now := time.Now()
	data := models.SysRole{}
	data.RoleName = c.RoleName
	data.Status = c.Status
	data.RoleKey = c.RoleKey
	data.RoleSort = c.RoleSort
	data.Flag = c.Flag
	data.Remark = c.Remark
	data.Admin = c.Admin
	data.DataScope = c.DataScope
	data.SysMenu = &c.SysMenu
	data.SysDept = c.SysDept
	data.CreateBy = c.CurrAdminId
	data.UpdateBy = c.CurrAdminId
	data.CreatedAt = &now
	data.UpdatedAt = &now

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysRoleService Insert error:%s", err)
		return err
	}

	for _, menu := range dataMenu {
		for _, api := range menu.SysApi {
			_, err = cb.AddNamedPolicy("p", data.RoleKey, api.Path, api.Action)
		}
	}
	_ = cb.SavePolicy()
	//if len(c.MenuIds) > 0 {
	//	s := SysRoleMenu{}
	//	s.Orm = e.Orm
	//	s.Log = e.Log
	//	err = s.ReloadRule(tx, c.RoleId, c.MenuIds)
	//	if err != nil {
	//		e.Log.Errorf("reload casbin rule error, %", err.Error())
	//		return err
	//	}
	//}
	return nil
}

// Update 修改SysRole对象
func (e *SysRole) Update(c *dto.SysRoleControl, cb *casbin.SyncedEnforcer) error {
	if c.RoleId <= 0 || c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var model = models.SysRole{}
	var mlist = make([]models.SysMenu, 0)
	tx.Preload("SysMenu").First(&model, c.RoleId)
	tx.Preload("SysApi").Where("menu_id in ?", c.MenuIds).Find(&mlist)
	err = tx.Model(&model).Association("SysMenu").Delete(model.SysMenu)
	if err != nil {
		e.Log.Errorf("SysRoleService Update error:%s", err)
		return errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	if c.RoleId > 0 {
		model.RoleId = c.RoleId
	}
	now := time.Now()
	model.RoleName = c.RoleName
	model.Status = c.Status
	model.RoleKey = c.RoleKey
	model.RoleSort = c.RoleSort
	model.Flag = c.Flag
	model.Remark = c.Remark
	model.Admin = c.Admin
	model.DataScope = c.DataScope
	model.SysMenu = &c.SysMenu
	model.UpdatedAt = &now
	model.UpdateBy = c.CurrAdminId
	model.SysDept = c.SysDept

	model.SysMenu = &mlist
	err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model).Error
	if err != nil {
		e.Log.Errorf("SysRoleService Update error:%s", err)
		return err
	}

	_, err = cb.RemoveFilteredPolicy(0, model.RoleKey)
	if err != nil {
		e.Log.Errorf("SysRoleService Update error:%s", err)
		return err
	}

	for _, menu := range mlist {
		for _, api := range menu.SysApi {
			_, err = cb.AddNamedPolicy("p", model.RoleKey, api.Path, api.Action)
		}
	}
	_ = cb.SavePolicy()
	return nil
}

// Remove 删除SysRole
func (e *SysRole) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	tx := e.Orm.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var model = models.SysRole{}
	err = tx.Preload("SysMenu").Preload("SysDept").First(&model, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	err = tx.Select(clause.Associations).Delete(&model).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetRoleMenuId 获取角色对应的菜单ids
func (e *SysRole) GetRoleMenuId(roleId int64) ([]int, error) {
	menuIds := make([]int, 0)
	model := models.SysRole{}
	model.RoleId = roleId
	if err := e.Orm.Model(&model).Preload("SysMenu").First(&model).Error; err != nil {
		return nil, err
	}
	l := *model.SysMenu
	for i := 0; i < len(l); i++ {
		menuIds = append(menuIds, l[i].MenuId)
	}
	return menuIds, nil
}

func (e *SysRole) UpdateDataScope(c *dto.RoleDataScopeReq) error {
	var err error
	tx := e.Orm.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var dlist = make([]models.SysDept, 0)
	var model = models.SysRole{}
	tx.Preload("SysDept").First(&model, c.RoleId)
	tx.Where("id in ?", c.DeptIds).Find(&dlist)
	err = tx.Model(&model).Association("SysDept").Delete(model.SysDept)
	if err != nil {
		e.Log.Errorf("SysRoleService UpdateDataScope error:%s", err)
		return err
	}
	if c.RoleId > 0 {
		model.RoleId = c.RoleId
	}
	model.DataScope = c.DataScope
	model.DeptIds = c.DeptIds
	model.SysDept = dlist
	err = tx.Model(&model).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model).Error
	if err != nil {
		e.Log.Errorf("SysRoleService UpdateDataScope error:%s", err)
		return err
	}
	return nil
}

// UpdateStatus 修改SysRole对象status
func (e *SysRole) UpdateStatus(c *dto.UpdateStatusReq) error {
	var err error
	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var model = models.SysRole{}
	err = tx.First(&model, c.RoleId).Error
	if err != nil {
		e.Log.Errorf("SysRoleService UpdateStatus error:%s", err)
		return errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}
	model.Status = c.Status
	err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model).Error
	if err != nil {
		e.Log.Errorf("SysRoleService UpdateStatus error:%s", err)
		return errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}
	return nil
}

// GetWithName 获取SysRole对象
func (e *SysRole) GetWithName(d *dto.SysRoleByName, model *models.SysRole) error {
	err := e.Orm.Where("role_name = ?", d.RoleName).First(model).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	model.MenuIds, err = e.GetRoleMenuId(model.RoleId)
	if err != nil {
		e.Log.Errorf("SysRoleService GetWithName error:%s", err)
		return err
	}
	return nil
}

// GetById 获取SysRole对象
func (e *SysRole) GetById(roleId int64) ([]string, error) {
	permissions := make([]string, 0)
	model := models.SysRole{}
	model.RoleId = roleId
	if err := e.Orm.Model(&model).Preload("SysMenu").First(&model).Error; err != nil {
		return nil, err
	}
	l := *model.SysMenu
	for i := 0; i < len(l); i++ {
		permissions = append(permissions, l[i].Permission)
	}
	return permissions, nil
}
