package models

import (
	"time"
)

type SysRole struct {
	RoleId    int64      `json:"roleId" gorm:"primaryKey;autoIncrement"` // 角色编码
	RoleName  string     `json:"roleName" gorm:"size:128;"`              // 角色名称
	Status    string     `json:"status" gorm:"size:4;"`                  //
	RoleKey   string     `json:"roleKey" gorm:"size:128;"`               //角色代码
	RoleSort  int        `json:"roleSort" gorm:""`                       //角色排序
	Flag      string     `json:"flag" gorm:"size:128;"`                  //
	Remark    string     `json:"remark" gorm:"size:255;"`                //备注
	Admin     bool       `json:"admin" gorm:"size:4;"`
	DataScope string     `json:"dataScope" gorm:"size:128;"`
	Params    string     `json:"params" gorm:"-"`
	MenuIds   []int      `json:"menuIds" gorm:"-"`
	DeptIds   []int      `json:"deptIds" gorm:"-"`
	CreateBy  int64      `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy  int64      `json:"updateBy" gorm:"index;comment:更新者"`
	CreatedAt *time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	SysDept   []SysDept  `json:"sysDept" gorm:"many2many:sys_role_dept;foreignKey:RoleId;joinForeignKey:role_id;references:DeptId;joinReferences:dept_id;"`
	SysMenu   *[]SysMenu `json:"sysMenu" gorm:"many2many:sys_role_menu;foreignKey:RoleId;joinForeignKey:role_id;references:MenuId;joinReferences:menu_id;"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
