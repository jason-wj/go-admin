package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
)

// SysMenuSearch 列表或者搜索使用结构体
type SysMenuSearch struct {
	dto.Pagination `search:"-"`
	Title          string `form:"title" search:"type:contains;column:title;table:sys_menu" comment:"菜单名称"`  // 菜单名称
	Visible        int    `form:"visible" search:"type:exact;column:visible;table:sys_menu" comment:"显示状态"` // 显示状态
}

func (m *SysMenuSearch) GetNeedSearch() interface{} {
	return *m
}

// SysMenuControl 增、改使用的结构体
type SysMenuControl struct {
	MenuId      int             `uri:"id" comment:"编码"`            // 编码
	MenuName    string          `form:"menuName" comment:"菜单name"` //菜单name
	Title       string          `form:"title" comment:"显示名称"`      //显示名称
	Icon        string          `form:"icon" comment:"图标"`         //图标
	Path        string          `form:"path" comment:"路径"`         //路径
	Paths       string          `form:"paths" comment:"id路径"`      //id路径
	MenuType    string          `form:"menuType" comment:"菜单类型"`   //菜单类型
	SysApi      []models.SysApi `form:"sysApi"`
	Apis        []int           `form:"apis"`
	Action      string          `form:"action" comment:"请求方式"`      //请求方式
	Permission  string          `form:"permission" comment:"权限编码"`  //权限编码
	ParentId    int             `form:"parentId" comment:"上级菜单"`    //上级菜单
	NoCache     bool            `form:"noCache" comment:"是否缓存"`     //是否缓存
	Breadcrumb  string          `form:"breadcrumb" comment:"是否面包屑"` //是否面包屑
	Component   string          `form:"component" comment:"组件"`     //组件
	Sort        int             `form:"sort" comment:"排序"`          //排序
	Visible     string          `form:"visible" comment:"是否显示"`     //是否显示
	IsFrame     string          `form:"isFrame" comment:"是否frame"`  //是否frame
	CurrAdminId int64           `json:"-" comment:""`
	CreateBy    int             `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy    int             `json:"updateBy" gorm:"index;comment:更新者"`
}

// SysMenuById 获取单个或者删除的结构体
type SysMenuById struct {
	Ids []int64 `json:"ids"`
}

type SysMenuByGetId struct {
	Id int64 `json:"id"`
}

type MenuLabel struct {
	Id       int         `json:"id,omitempty" gorm:"-"`
	Label    string      `json:"label,omitempty" gorm:"-"`
	Children []MenuLabel `json:"children,omitempty" gorm:"-"`
}

type MenuRole struct {
	models.SysMenu
	IsSelect bool `json:"is_select" gorm:"-"`
}

type SelectRole struct {
	RoleId int64 `uri:"roleId"`
}
