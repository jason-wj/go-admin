package dto

import (
	"go-admin/common/dto"
)

// SysDeptSearch 列表或者搜索使用结构体
type SysDeptSearch struct {
	dto.Pagination `search:"-"`
	DeptId         int    `form:"deptId" search:"type:exact;column:dept_id;table:sys_dept" comment:"id"`       //id
	ParentId       int    `form:"parentId" search:"type:exact;column:parent_id;table:sys_dept" comment:"上级部门"` //上级部门
	DeptPath       string `form:"deptPath" search:"type:exact;column:dept_path;table:sys_dept" comment:""`     //路径
	DeptName       string `form:"deptName" search:"type:exact;column:dept_name;table:sys_dept" comment:"部门名称"` //部门名称
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_dept" comment:"排序"`            //排序
	Leader         string `form:"leader" search:"type:exact;column:leader;table:sys_dept" comment:"负责人"`       //负责人
	Phone          string `form:"phone" search:"type:exact;column:phone;table:sys_dept" comment:"手机"`          //手机
	Email          string `form:"email" search:"type:exact;column:email;table:sys_dept" comment:"邮箱"`          //邮箱
	Status         string `form:"status" search:"type:exact;column:status;table:sys_dept" comment:"状态"`        //状态
}

func (m *SysDeptSearch) GetNeedSearch() interface{} {
	return *m
}

// SysDeptControl 增、改使用的结构体
type SysDeptControl struct {
	DeptId      int64  `uri:"id" comment:"编码"`                                         // 编码
	ParentId    int64  `json:"parentId" comment:"上级部门" vd:"?"`                         //上级部门
	DeptPath    string `json:"deptPath" comment:""`                                    //路径
	DeptName    string `json:"deptName" comment:"部门名称" vd:"len($)>0"`                  //部门名称
	Sort        int    `json:"sort" comment:"排序" vd:"?"`                               //排序
	Leader      string `json:"leader" comment:"负责人" vd:"@:len($)>0; msg:'leader不能为空'"` //负责人
	Phone       string `json:"phone" comment:"手机" vd:"?"`                              //手机
	Email       string `json:"email" comment:"邮箱" vd:"?"`                              //邮箱
	Status      int    `json:"status" comment:"状态" vd:"$>0"`                           //状态
	CurrAdminId int64  `json:"-" comment:""`
}

// SysDeptById 获取单个或者删除的结构体
type SysDeptById struct {
	Id  int64   `uri:"id"`
	Ids []int64 `json:"ids"`
}

type DeptLabel struct {
	Id       int64       `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLabel `gorm:"-" json:"children"`
}
