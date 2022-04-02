package dto

import (
	"go-admin/common/dto"
)

type SysUserQueryReq struct {
	dto.Pagination `search:"-"`
	UserId         int    `form:"userId" search:"type:exact;column:user_id;table:sys_user" comment:"用户ID"`
	Username       string `form:"username" search:"type:contains;column:username;table:sys_user" comment:"用户名"`
	NickName       string `form:"nickName" search:"type:contains;column:nick_name;table:sys_user" comment:"昵称"`
	Phone          string `form:"phone" search:"type:contains;column:phone;table:sys_user" comment:"手机号"`
	RoleId         string `form:"roleId" search:"type:exact;column:role_id;table:sys_user" comment:"角色ID"`
	Sex            string `form:"sex" search:"type:exact;column:sex;table:sys_user" comment:"性别"`
	Email          string `form:"email" search:"type:contains;column:email;table:sys_user" comment:"邮箱"`
	PostId         string `form:"postId" search:"type:exact;column:post_id;table:sys_user" comment:"岗位"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_user" comment:"状态"`
	DeptJoin       `search:"type:inner;on:dept_id:dept_id;table:sys_user;join:sys_dept"`
	SysUserOrder
}

type SysUserOrder struct {
	UserIdOrder    string `search:"type:order;column:user_id;table:sys_user" form:"userIdOrder"`
	UsernameOrder  string `search:"type:order;column:username;table:sys_user" form:"usernameOrder"`
	StatusOrder    string `search:"type:order;column:status;table:sys_user" form:"statusOrder"`
	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_user" form:"createdAtOrder"`
}

type DeptJoin struct {
	DeptId string `search:"type:contains;column:dept_path;table:sys_dept" form:"deptId"`
}

func (m *SysUserQueryReq) GetNeedSearch() interface{} {
	return *m
}

type ResetSysUserPwdReq struct {
	UserId      int64  `json:"userId" comment:"用户ID" binding:"required"` // 用户ID
	Password    string `json:"password" comment:"密码" binding:"required"`
	CurrAdminId int64  `json:"-" comment:""`
}

type UpdateSysUserAvatarReq struct {
	UserId      int64  `json:"userId" comment:"用户ID" vd:"len($)>0"` // 用户ID
	Avatar      string `json:"avatar" comment:"头像" vd:"len($)>0"`
	CurrAdminId int64  `json:"-" comment:""`
}

type UpdateSysUserStatusReq struct {
	UserId      int64  `json:"userId" comment:"用户ID" vd:"$>0"` // 用户ID
	Status      string `json:"status" comment:"状态" vd:"len($)>0"`
	CurrAdminId int64  `json:"-" comment:""`
}

type SysUserInsertReq struct {
	UserId      int64  `json:"userId" comment:"用户ID"` // 用户ID
	Username    string `json:"username" comment:"用户名" vd:"len($)>0"`
	Password    string `json:"password" comment:"密码"`
	NickName    string `json:"nickName" comment:"昵称" vd:"len($)>0"`
	Phone       string `json:"phone" comment:"手机号" vd:"len($)>0"`
	RoleId      int    `json:"roleId" comment:"角色ID"`
	Avatar      string `json:"avatar" comment:"头像"`
	Sex         string `json:"sex" comment:"性别"`
	Email       string `json:"email" comment:"邮箱" vd:"len($)>0,email"`
	DeptId      int    `json:"deptId" comment:"部门" vd:"$>0"`
	PostId      int    `json:"postId" comment:"岗位"`
	Remark      string `json:"remark" comment:"备注"`
	Status      string `json:"status" comment:"状态" vd:"len($)>0" default:"1"`
	CurrAdminId int64  `json:"-" comment:""`
}

type SysUserUpdateReq struct {
	UserId      int64  `json:"userId" comment:"用户ID"` // 用户ID
	Username    string `json:"username" comment:"用户名"`
	NickName    string `json:"nickName" comment:"昵称"`
	Phone       string `json:"phone" comment:"手机号"`
	RoleId      int    `json:"roleId" comment:"角色ID"`
	Avatar      string `json:"avatar" comment:"头像"`
	Sex         string `json:"sex" comment:"性别"`
	Email       string `json:"email" comment:"邮箱"`
	DeptId      int    `json:"deptId" comment:"部门"`
	PostId      int    `json:"postId" comment:"岗位"`
	Remark      string `json:"remark" comment:"备注"`
	Status      string `json:"status" comment:"状态" default:"1"`
	CurrAdminId int64  `json:"-" comment:""`
	Password    string `json:"password" comment:""`
}

type SysUserById struct {
	dto.ObjectById
	CurrAdminId int64 `json:"-" comment:""`
}

// PassWord 密码
type PassWord struct {
	NewPassword string `json:"newPassword" vd:"len($)>0"`
	OldPassword string `json:"oldPassword" vd:"len($)>0"`
}

type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Code     string `form:"code" json:"code" binding:"required"`
	UUID     string `form:"uuid" json:"uuid" binding:"required"`
}

// SysUserDeleteReq 功能删除请求参数
type SysUserDeleteReq struct {
	Ids int64 `json:"ids"`
}
