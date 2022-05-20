package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
)

// SysRoleSearch 列表或者搜索使用结构体
type SysRoleSearch struct {
	dto.Pagination `search:"-"`

	RoleId    int    `form:"roleId" search:"type:exact;column:role_id;table:sys_role" comment:"角色编码"`     // 角色编码
	RoleName  string `form:"roleName" search:"type:exact;column:role_name;table:sys_role" comment:"角色名称"` // 角色名称
	Status    string `form:"status" search:"type:exact;column:status;table:sys_role" comment:"状态"`        // 状态
	RoleKey   string `form:"roleKey" search:"type:exact;column:role_key;table:sys_role" comment:"角色代码"`   // 角色代码
	RoleSort  int    `form:"roleSort" search:"type:exact;column:role_sort;table:sys_role" comment:"角色排序"` // 角色排序
	Flag      string `form:"flag" search:"type:exact;column:flag;table:sys_role" comment:"标记"`            // 标记
	Remark    string `form:"remark" search:"type:exact;column:remark;table:sys_role" comment:"备注"`        // 备注
	Admin     bool   `form:"admin" search:"type:exact;column:admin;table:sys_role" comment:"是否管理员"`
	DataScope string `form:"dataScope" search:"type:exact;column:data_scope;table:sys_role" comment:"是否管理员"`
}

type SysRoleOrder struct {
	RoleIdOrder    string `search:"type:order;column:role_id;table:sys_role" form:"roleIdOrder"`
	RoleNameOrder  string `search:"type:order;column:role_name;table:sys_role" form:"roleNameOrder"`
	RoleSortOrder  string `search:"type:order;column:role_sort;table:sys_role" form:"usernameOrder"`
	StatusOrder    string `search:"type:order;column:status;table:sys_role" form:"statusOrder"`
	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_role" form:"createdAtOrder"`
}

func (m *SysRoleSearch) GetNeedSearch() interface{} {
	return *m
}

type SysRoleControl struct {
	RoleId     int64            `uri:"id" comment:"角色编码"`        // 角色编码
	RoleName   string           `form:"roleName" comment:"角色名称"` // 角色名称
	Status     string           `form:"status" comment:"状态"`     // 状态
	RoleKey    string           `form:"roleKey" comment:"角色代码"`  // 角色代码
	RoleSort   int              `form:"roleSort" comment:"角色排序"` // 角色排序
	Flag       string           `form:"flag" comment:"标记"`       // 标记
	Remark     string           `form:"remark" comment:"备注"`     // 备注
	Admin      bool             `form:"admin" comment:"是否管理员"`
	DataScope  string           `form:"dataScope"`
	SysMenu    []models.SysMenu `form:"sysMenu"`
	MenuIds    []int            `form:"menuIds"`
	SysDept    []models.SysDept `form:"sysDept"`
	DeptIds    []int            `form:"deptIds"`
	CurrUserId int64            `json:"-" comment:""`
}

type UpdateStatusReq struct {
	RoleId     int64  `form:"roleId" comment:"角色编码"` // 角色编码
	Status     string `form:"status" comment:"状态"`   // 状态
	CurrUserId int64  `json:"-" comment:""`
}

type SysRoleByName struct {
	RoleName string `form:"role"` // 角色编码
}

// SysRoleById 获取单个或者删除的结构体
type SysRoleById struct {
	Id  int64   `uri:"id"`
	Ids []int64 `json:"ids"`
}

// RoleDataScopeReq 角色数据权限修改
type RoleDataScopeReq struct {
	RoleId    int64  `json:"roleId" binding:"required"`
	DataScope string `json:"dataScope" binding:"required"`
	DeptIds   []int  `json:"deptIds"`
}

type DeptIdList struct {
	DeptId int64 `json:"DeptId"`
}

// SysRoleDeleteReq 功能删除请求参数
type SysRoleDeleteReq struct {
	Ids []int64 `json:"ids"`
}
