package service

import (
	"errors"
	"fmt"
	"go-admin/app/admin/models/tools"
	"strings"
	"time"

	"go-admin/common/core/sdk/pkg"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	cModels "go-admin/common/models"

	"go-admin/common/core/sdk/service"
)

type SysMenu struct {
	service.Service
}

func NewSysMenuService(s *service.Service) *SysMenu {
	var srv = new(SysMenu)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取SysMenu列表
func (e *SysMenu) GetPage(c *dto.SysMenuSearch) ([]models.SysMenu, error) {
	var menus []models.SysMenu
	list, count, err := e.getPage(c)
	if err != nil {
		return nil, nil
	}
	for i := 0; i < int(count); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&list, list[i])
		menus = append(menus, menusInfo)
	}
	return menus, nil
}

// getPage 菜单分页列表
func (e *SysMenu) getPage(c *dto.SysMenuSearch) ([]models.SysMenu, int64, error) {

	var list []models.SysMenu
	var data models.SysMenu
	var count int64

	err := e.Orm.Model(&data).
		Scopes(
			cDto.OrderDest("sort", false),
			cDto.MakeCondition(c.GetNeedSearch()),
		).Preload("SysApi").Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysMenuService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取SysMenu对象
func (e *SysMenu) Get(id int64) (*models.SysMenu, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysMenu{}
	err := e.Orm.Preload("SysApi").First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	apis := make([]int, 0)
	for _, v := range model.SysApi {
		apis = append(apis, v.Id)
	}
	model.Apis = apis
	return model, nil
}

// Insert 创建SysMenu对象
func (e *SysMenu) Insert(c *dto.SysMenuControl) error {
	var err error
	now := time.Now()
	data := models.SysMenu{}
	data.Name = c.Name
	data.Title = c.Title
	data.Icon = c.Icon
	data.Path = c.Path
	data.Paths = c.Paths
	data.MenuType = c.MenuType
	data.Action = c.Action
	data.SysApi = c.SysApi
	data.Permission = c.Permission
	data.ParentId = c.ParentId
	data.KeepAlive = c.KeepAlive
	data.Breadcrumb = c.Breadcrumb
	data.Component = c.Component
	data.Sort = c.Sort
	data.Hidden = c.Hidden
	data.IsFrame = c.IsFrame
	data.CreateBy = 1
	data.UpdateBy = 1
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysMenuService Insert error:%s", err)
		return err
	}
	return nil
}

func (e *SysMenu) InsertConf(id int) error {
	var err error
	tx := e.Orm.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	table := tools.SysTables{}

	table.TableId = id
	tab, _ := table.Get(e.Orm, true)
	tab.MLTBName = strings.Replace(tab.TBName, "_", "-", -1)

	Mmenu := dto.SysMenuControl{}
	Mmenu.Title = tab.TableComment
	Mmenu.Path = "/" + tab.MLTBName
	Mmenu.Icon = "pass"
	Mmenu.MenuType = "M"
	Mmenu.Action = "无"
	Mmenu.ParentId = 0
	Mmenu.KeepAlive = false
	Mmenu.Component = "Layout"
	Mmenu.Sort = 0
	Mmenu.Hidden = false
	Mmenu.IsFrame = "0"
	Mmenu.CreateBy = 1
	err = e.Insert(&Mmenu)
	if err != nil {
		return err
	}

	Cmenu := dto.SysMenuControl{}
	Cmenu.Name = tab.ClassName + "Manage"
	Cmenu.Title = tab.TableComment
	Cmenu.Icon = "pass"
	Cmenu.Path = "/" + tab.PackageName + "/" + tab.MLTBName
	Cmenu.MenuType = "C"
	Cmenu.Action = "无"
	Cmenu.Permission = tab.PackageName + ":" + tab.BusinessName + ":list"
	Cmenu.ParentId = Mmenu.MenuId
	Cmenu.KeepAlive = false
	Cmenu.Component = "/" + tab.PackageName + "/" + tab.MLTBName + "/index"
	Cmenu.Sort = 0
	Cmenu.Hidden = false
	Cmenu.IsFrame = "0"
	Cmenu.CreateBy = 1
	Cmenu.UpdateBy = 1
	err = e.Insert(&Cmenu)
	if err != nil {
		return err
	}

	MList := dto.SysMenuControl{}
	MList.Name = ""
	MList.Title = "分页获取" + tab.TableComment
	MList.Icon = ""
	MList.Path = tab.TBName
	MList.MenuType = "F"
	MList.Action = "无"
	MList.Permission = tab.PackageName + ":" + tab.BusinessName + ":query"
	MList.ParentId = Cmenu.MenuId
	MList.KeepAlive = false
	MList.Sort = 0
	MList.Hidden = false
	MList.IsFrame = "0"
	MList.CreateBy = 1
	MList.UpdateBy = 1
	err = e.Insert(&MList)
	if err != nil {
		return err
	}

	MCreate := dto.SysMenuControl{}
	MCreate.Name = ""
	MCreate.Title = "创建" + tab.TableComment
	MCreate.Icon = ""
	MCreate.Path = tab.TBName
	MCreate.MenuType = "F"
	MCreate.Action = "无"
	MCreate.Permission = tab.PackageName + ":" + tab.BusinessName + ":add"
	MCreate.ParentId = Cmenu.MenuId
	MCreate.KeepAlive = false
	MCreate.Sort = 0
	MCreate.Hidden = false
	MCreate.IsFrame = "0"
	MCreate.CreateBy = 1
	MCreate.UpdateBy = 1
	err = e.Insert(&MCreate)
	if err != nil {
		return err
	}

	MUpdate := dto.SysMenuControl{}
	MUpdate.Name = ""
	MUpdate.Title = "修改" + tab.TableComment
	MUpdate.Icon = ""
	MUpdate.Path = tab.TBName
	MUpdate.MenuType = "F"
	MUpdate.Action = "无"
	MUpdate.Permission = tab.PackageName + ":" + tab.BusinessName + ":edit"
	MUpdate.ParentId = Cmenu.MenuId
	MUpdate.KeepAlive = false
	MUpdate.Sort = 0
	MUpdate.Hidden = false
	MUpdate.IsFrame = "0"
	MUpdate.CreateBy = 1
	MUpdate.UpdateBy = 1
	err = e.Insert(&MUpdate)
	if err != nil {
		return err
	}

	MDelete := dto.SysMenuControl{}
	MDelete.Name = ""
	MDelete.Title = "删除" + tab.TableComment
	MDelete.Icon = ""
	MDelete.Path = tab.TBName
	MDelete.MenuType = "F"
	MDelete.Action = "无"
	MDelete.Permission = tab.PackageName + ":" + tab.BusinessName + ":remove"
	MDelete.ParentId = Cmenu.MenuId
	MDelete.KeepAlive = false
	MDelete.Sort = 0
	MDelete.Hidden = false
	MDelete.IsFrame = "0"
	MDelete.CreateBy = 1
	MDelete.UpdateBy = 1
	err = e.Insert(&MDelete)
	if err != nil {
		return err
	}
	return nil
}

func (e *SysMenu) initPaths(menu *models.SysMenu) error {
	var err error
	var data models.SysMenu
	parentMenu := new(models.SysMenu)
	if menu.ParentId != 0 {
		e.Orm.Model(&data).First(parentMenu, menu.ParentId)
		if parentMenu.Paths == "" {
			return errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
		}
		menu.Paths = parentMenu.Paths + "/" + pkg.IntToString(menu.MenuId)
	} else {
		menu.Paths = "/0/" + pkg.IntToString(menu.MenuId)
	}
	e.Orm.Model(&data).Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths)
	return err
}

// Update 修改SysMenu对象
func (e *SysMenu) Update(c *dto.SysMenuControl) error {
	if c.MenuId <= 0 || c.CurrUserId <= 0 {
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
	var alist = make([]models.SysApi, 0)
	var model = models.SysMenu{}
	tx.Model(&models.SysMenu{}).Preload("SysApi").First(&model, c.MenuId)
	tx.Where("id in ?", c.Apis).Find(&alist)
	err = tx.Model(&model).Association("SysApi").Delete(model.SysApi)
	if err != nil {
		e.Log.Errorf("SysMenuService Update error:%s", err)
		return err
	}

	if c.MenuId > 0 {
		model.MenuId = c.MenuId
	}
	now := time.Now()
	model.Name = c.Name
	model.Title = c.Title
	model.Icon = c.Icon
	model.Path = c.Path
	model.Paths = c.Paths
	model.MenuType = c.MenuType
	model.Action = c.Action
	model.Permission = c.Permission
	model.ParentId = c.ParentId
	model.KeepAlive = c.KeepAlive
	model.Breadcrumb = c.Breadcrumb
	model.Component = c.Component
	model.Sort = c.Sort
	model.Hidden = c.Hidden
	model.IsFrame = c.IsFrame
	model.UpdateBy = c.CurrUserId
	model.UpdatedAt = &now
	model.SysApi = alist
	err = tx.Model(&model).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model).Error
	if err != nil {
		e.Log.Errorf("SysConfigService Update error:%s", err)
		return err
	}
	return nil
}

// Remove 删除SysMenu
func (e *SysMenu) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.SysMenu

	err = e.Orm.Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetList 获取菜单数据
func (e *SysMenu) GetList(c *dto.SysMenuSearch) ([]models.SysMenu, error) {
	var list []models.SysMenu
	var err error
	err = e.Orm.Model(&models.SysMenu{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Find(&list).Error
	if err != nil {
		e.Log.Errorf("SysConfigService GetList error:%s", err)
		return nil, err
	}
	return list, nil
}

// SetLabel 修改角色中 设置菜单基础数据
func (e *SysMenu) SetLabel() (m []dto.MenuLabel, err error) {
	var list []models.SysMenu
	list, err = e.GetList(&dto.SysMenuSearch{})

	m = make([]dto.MenuLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		e := dto.MenuLabel{}
		e.Id = list[i].MenuId
		e.Label = list[i].Title
		deptsInfo := menuLabelCall(&list, e)

		m = append(m, deptsInfo)
	}
	return
}

// GetSysMenuByRoleName 左侧菜单
func (e *SysMenu) GetSysMenuByRoleName(roleName ...string) ([]models.SysMenu, error) {
	var MenuList []models.SysMenu
	var role models.SysRole
	var err error
	admin := false
	for _, s := range roleName {
		if s == "admin" {
			admin = true
		}
	}

	if len(roleName) > 0 && admin {
		var data []models.SysMenu
		err = e.Orm.Where(" menu_type in ('M','C')").
			Order("sort").
			Find(&data).
			Error
		MenuList = data
	} else {
		err = e.Orm.Model(&role).Preload("SysMenu", func(db *gorm.DB) *gorm.DB {
			return db.Where(" menu_type in ('M','C')").Order("sort")
		}).Where("role_name in ?", roleName).Find(&role).
			Error
		MenuList = *role.SysMenu
	}

	if err != nil {
		e.Log.Errorf("db error:%s", err)
	}
	return MenuList, err
}

// menuLabelCall 递归构造组织数据
func menuLabelCall(eList *[]models.SysMenu, dept dto.MenuLabel) dto.MenuLabel {
	list := *eList

	min := make([]dto.MenuLabel, 0)
	for j := 0; j < len(list); j++ {

		if dept.Id != list[j].ParentId {
			continue
		}
		mi := dto.MenuLabel{}
		mi.Id = list[j].MenuId
		mi.Label = list[j].Title
		mi.Children = []dto.MenuLabel{}
		if list[j].MenuType != "F" {
			ms := menuLabelCall(eList, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	if len(min) > 0 {
		dept.Children = min
	} else {
		dept.Children = nil
	}
	return dept
}

// menuCall 构建菜单树
func menuCall(menuList *[]models.SysMenu, menu models.SysMenu) models.SysMenu {
	list := *menuList

	min := make([]models.SysMenu, 0)
	for j := 0; j < len(list); j++ {

		if menu.MenuId != list[j].ParentId {
			continue
		}
		mi := models.SysMenu{}
		mi.MenuId = list[j].MenuId
		mi.Name = list[j].Name
		mi.Title = list[j].Title
		mi.Icon = list[j].Icon
		mi.Path = list[j].Path
		mi.MenuType = list[j].MenuType
		mi.Action = list[j].Action
		mi.Permission = list[j].Permission
		mi.ParentId = list[j].ParentId
		mi.KeepAlive = list[j].KeepAlive
		mi.Breadcrumb = list[j].Breadcrumb
		mi.Component = list[j].Component
		mi.Sort = list[j].Sort
		mi.Hidden = list[j].Hidden
		mi.CreatedAt = list[j].CreatedAt
		mi.SysApi = list[j].SysApi
		mi.Children = []models.SysMenu{}

		if mi.MenuType != cModels.Button {
			ms := menuCall(menuList, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	menu.Children = min
	return menu
}

// SetMenuRole 获取左侧菜单树使用
func (e *SysMenu) SetMenuRole(roleName string) (m []models.SysMenu, err error) {
	menus, err := e.getByRoleName(roleName)
	m = make([]models.SysMenu, 0)
	for i := 0; i < len(menus); i++ {
		if menus[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&menus, menus[i])
		m = append(m, menusInfo)
	}
	return
}

func (e *SysMenu) getByRoleName(roleName string) ([]models.SysMenu, error) {
	var MenuList []models.SysMenu
	var role models.SysRole
	var err error

	if roleName == "admin" {
		var data []models.SysMenu
		err = e.Orm.Where(" menu_type in ('M','C')").Order("sort").Find(&data).Error
		MenuList = data
	} else {
		role.RoleKey = roleName
		err = e.Orm.Debug().Model(&role).Where("role_key = ? ", roleName).Preload("SysMenu", func(db *gorm.DB) *gorm.DB {
			return db.Where(" menu_type in ('M','C')").Order("sort")
		}).Find(&role).Error
		if role.SysMenu != nil {
			MenuList = *role.SysMenu
		}
	}

	if err != nil {
		e.Log.Errorf("db error:%s", err)
	}
	return MenuList, err
}
