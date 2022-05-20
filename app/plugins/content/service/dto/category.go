package dto

import (
	"go-admin/common/dto"
	"time"
)

type CategoryQueryReq struct {
	dto.Pagination `search:"-"`
	Id             int64  `form:"id"  search:"type:exact;column:id;table:content_category" comment:"编号"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:content_category" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:content_category" comment:"创建时间"`
	Name           string `form:"name"  search:"type:contains;column:name;table:content_category" comment:"分类名称"`
	NameInner      string `form:"-"  search:"type:exact;column:name;table:content_category" comment:"分类名称"`
	CategoryOrder
}

type CategoryOrder struct {
	IdOrder        int64     `form:"idOrder"  search:"type:order;column:id;table:content_category"`
	NameOrder      string    `form:"nameOrder"  search:"type:order;column:name;table:content_category"`
	CreateByOrder  string    `form:"createByOrder"  search:"type:order;column:create_by;table:content_category"`
	CreatedAtOrder time.Time `form:"createdAtOrder"  search:"type:order;column:created_at;table:content_category"`
	UpdateByOrder  string    `form:"updateByOrder"  search:"type:order;column:update_by;table:content_category"`
	UpdatedAtOrder time.Time `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:content_category"`
	StatusOrder    string    `form:"statusOrder"  search:"type:order;column:status;table:content_category"`
	RemarkOrder    string    `form:"remarkOrder"  search:"type:order;column:remark;table:content_category"`
}

func (m *CategoryQueryReq) GetNeedSearch() interface{} {
	return *m
}

type CategoryInsertReq struct {
	Name       string `json:"name" comment:"分类名称"`
	Status     string `json:"status" comment:"状态（0正常 1删除 2停用 3冻结）"`
	Remark     string `json:"remark" comment:"备注信息"`
	CurrUserId int64  `json:"-" comment:"当前登陆用户"`
}

type CategoryUpdateReq struct {
	Id         int64  `uri:"id" comment:"主键编码"` // 主键编码
	Name       string `json:"name" comment:"分类名称"`
	Status     string `json:"status" comment:"状态（0正常 1删除 2停用 3冻结）"`
	Remark     string `json:"remark" comment:"备注信息"`
	CurrUserId int64  `json:"-" comment:""`
}

// CategoryGetReq 功能获取请求参数
type CategoryGetReq struct {
	Id int64 `uri:"id"`
}

// CategoryDeleteReq 功能删除请求参数
type CategoryDeleteReq struct {
	Ids []int64 `json:"ids"`
}
