package dto

import (
	"go-admin/common/dto"
)

type SysDictDataSearch struct {
	dto.Pagination `search:"-"`
	Id             int    `form:"id" search:"type:exact;column:dict_code;table:sys_dict_data" comment:""`
	DictLabel      string `form:"dictLabel" search:"type:contains;column:dict_label;table:sys_dict_data" comment:""`
	DictValue      string `form:"dictValue" search:"type:contains;column:dict_value;table:sys_dict_data" comment:""`
	DictType       string `form:"dictType" search:"type:contains;column:dict_type;table:sys_dict_data" comment:""`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_dict_data" comment:""`
}

func (m *SysDictDataSearch) GetNeedSearch() interface{} {
	return *m
}

type SysDictDataControl struct {
	Id          int    `uri:"dictCode" comment:""`
	DictSort    int    `json:"dictSort" comment:""`
	DictLabel   string `json:"dictLabel" comment:""`
	DictValue   string `json:"dictValue" comment:""`
	DictType    string `json:"dictType" comment:""`
	CssClass    string `json:"cssClass" comment:""`
	ListClass   string `json:"listClass" comment:""`
	IsDefault   string `json:"isDefault" comment:""`
	Status      int    `json:"status" comment:""`
	Default     string `json:"default" comment:""`
	Remark      string `json:"remark" comment:""`
	CurrAdminId int64  `json:"-" comment:""`
}

type SysDictDataById struct {
	Id  int64   `uri:"dictCode"`
	Ids []int64 `json:"ids"`
}
