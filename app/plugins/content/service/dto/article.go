package dto

import (
	"go-admin/common/dto"

	"time"
)

type ArticleQueryReq struct {
	dto.Pagination `search:"-"`
	Id             int64  `form:"id"  search:"type:exact;column:id;table:content_article" comment:"编号"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:content_article" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:content_article" comment:"创建时间"`
	CateId         int64  `form:"cateId"  search:"type:exact;column:cate_id;table:content_article" comment:"分类编号"`
	Name           string `form:"name"  search:"type:contains;column:name;table:content_article" comment:"标题"`
	ArticleOrder
}

type ArticleOrder struct {
	IdOrder        int64     `form:"idOrder"  search:"type:order;column:id;table:content_article"`
	CateIdOrder    int64     `form:"cateIdOrder"  search:"type:order;column:cate_id;table:content_article"`
	NameOrder      string    `form:"nameOrder"  search:"type:order;column:name;table:content_article"`
	ContentOrder   string    `form:"contentOrder"  search:"type:order;column:content;table:content_article"`
	RemarkOrder    string    `form:"remarkOrder"  search:"type:order;column:remark;table:content_article"`
	StatusOrder    string    `form:"statusOrder"  search:"type:order;column:status;table:content_article"`
	CreateByOrder  string    `form:"createByOrder"  search:"type:order;column:create_by;table:content_article"`
	UpdateByOrder  string    `form:"updateByOrder"  search:"type:order;column:update_by;table:content_article"`
	UpdatedAtOrder time.Time `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:content_article"`
	CreatedAtOrder time.Time `form:"createdAtOrder"  search:"type:order;column:created_at;table:content_article"`
}

func (m *ArticleQueryReq) GetNeedSearch() interface{} {
	return *m
}

type ArticleInsertReq struct {
	CateId      int64  `json:"cateId" comment:"分类编号"`
	Name        string `json:"name" comment:"标题"`
	Content     string `json:"content" comment:"文本内容"`
	Remark      string `json:"remark" comment:"备注信息"`
	CurrAdminId int64  `json:"-" comment:""`
}

type ArticleUpdateReq struct {
	Id          int64  `uri:"id" comment:"主键编码"` // 主键编码
	CateId      int64  `json:"cateId" comment:"分类编号"`
	Name        string `json:"name" comment:"标题"`
	Content     string `json:"content" comment:"文本内容"`
	Remark      string `json:"remark" comment:"备注信息"`
	Status      string `json:"status" comment:"状态（0正常 1删除 2停用 3冻结）"`
	CurrAdminId int64  `json:"-" comment:""`
}

// ArticleGetReq 功能获取请求参数
type ArticleGetReq struct {
	Id int64 `uri:"id"`
}

// ArticleDeleteReq 功能删除请求参数
type ArticleDeleteReq struct {
	Ids []int64 `json:"ids"`
}
