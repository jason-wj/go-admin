package dto

import (
	"go-admin/common/dto"
)

type SysDictTypeSearch struct {
	dto.Pagination `search:"-"`
	DictId         []int  `form:"dictId" search:"type:in;column:dict_id;table:sys_dict_type"`
	DictName       string `form:"dictName" search:"type:icontains;column:dict_name;table:sys_dict_type"`
	DictType       string `form:"dictType" search:"type:icontains;column:dict_type;table:sys_dict_type"`
	Status         int    `form:"status" search:"type:exact;column:status;table:sys_dict_type"`
}

type SysDictTypeOrder struct {
	DictIdOrder string `search:"type:order;column:dict_id;table:sys_dict_type" form:"dictIdOrder"`
}

func (m *SysDictTypeSearch) GetNeedSearch() interface{} {
	return *m
}

type SysDictTypeControl struct {
	Id          int    `uri:"id"`
	DictName    string `json:"dictName"`
	DictType    string `json:"dictType"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
	CurrAdminId int64  `json:"-" comment:""`
}

type SysDictTypeById struct {
	dto.ObjectById
	CurrAdminId int64 `json:"-" comment:""`
}

// SysDictDeleteReq 功能删除请求参数
type SysDictrDeleteReq struct {
	Ids []int64 `json:"ids"`
}
