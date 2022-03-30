package models

import "time"

type SysMenu struct {
	MenuId     int        `json:"menuId" gorm:"primaryKey;autoIncrement"`
	Name       string     `json:"name" gorm:"size:128;"`
	Title      string     `json:"title" gorm:"size:128;"`
	Icon       string     `json:"icon" gorm:"size:128;"`
	Path       string     `json:"path" gorm:"size:128;"`
	Paths      string     `json:"paths" gorm:"size:128;"`
	MenuType   string     `json:"menuType" gorm:"size:1;"`
	Action     string     `json:"action" gorm:"size:16;"`
	Permission string     `json:"permission" gorm:"size:255;"`
	ParentId   int        `json:"parentId" gorm:"size:11;"`
	NoCache    bool       `json:"noCache" gorm:"size:8;"`
	Breadcrumb string     `json:"breadcrumb" gorm:"size:255;"`
	Component  string     `json:"component" gorm:"size:255;"`
	Sort       int        `json:"sort" gorm:"size:4;"`
	Visible    string     `json:"visible" gorm:"size:1;"`
	IsFrame    string     `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
	SysApi     []SysApi   `json:"sysApi" gorm:"many2many:sys_menu_api_rule"`
	CreateBy   int64      `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy   int64      `json:"updateBy" gorm:"index;comment:更新者"`
	CreatedAt  *time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt  *time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	Apis       []int      `json:"apis" gorm:"-"`
	DataScope  string     `json:"dataScope" gorm:"-"`
	Params     string     `json:"params" gorm:"-"`
	RoleId     int        `gorm:"-"`
	Children   []SysMenu  `json:"children,omitempty" gorm:"-"`
	IsSelect   bool       `json:"is_select" gorm:"-"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
