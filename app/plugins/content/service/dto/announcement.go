package dto

import (
	"go-admin/common/dto"
	"time"
)

type AnnouncementQueryReq struct {
	dto.Pagination `search:"-"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:content_announcement" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:content_announcement" comment:"创建时间"`
	Title          string `form:"title"  search:"type:contains;column:title;table:content_announcement" comment:"标题"`
	Status         string `form:"status"  search:"type:exact;column:status;table:content_announcement" comment:"状态（0正常 1删除 2停用 3冻结）"`
	AnnouncementOrder
}

type AnnouncementOrder struct {
	IdOrder        int64      `form:"idOrder"  search:"type:order;column:id;table:content_announcement"`
	TitleOrder     string     `form:"titleOrder"  search:"type:order;column:title;table:content_announcement"`
	ContentOrder   string     `form:"contentOrder"  search:"type:order;column:content;table:content_announcement"`
	NumOrder       int64      `form:"numOrder"  search:"type:order;column:num;table:content_announcement"`
	RemarkOrder    string     `form:"remarkOrder"  search:"type:order;column:remark;table:content_announcement"`
	StatusOrder    string     `form:"statusOrder"  search:"type:order;column:status;table:content_announcement"`
	CreateByOrder  string     `form:"createByOrder"  search:"type:order;column:create_by;table:content_announcement"`
	UpdateByOrder  string     `form:"updateByOrder"  search:"type:order;column:update_by;table:content_announcement"`
	UpdatedAtOrder *time.Time `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:content_announcement"`
	CreatedAtOrder *time.Time `form:"createdAtOrder"  search:"type:order;column:created_at;table:content_announcement"`
}

func (m *AnnouncementQueryReq) GetNeedSearch() interface{} {
	return *m
}

type AnnouncementInsertReq struct {
	Title       string `json:"title" comment:"标题"`
	Content     string `json:"content" comment:"内容"`
	Num         int    `json:"num" comment:"阅读次数"`
	Remark      string `json:"remark" comment:"备注信息"`
	CurrAdminId int64  `json:"-" comment:""`
}

type AnnouncementUpdateReq struct {
	Id          int64  `uri:"id" comment:"主键编码"` // 主键编码
	Title       string `json:"title" comment:"标题"`
	Content     string `json:"content" comment:"内容"`
	Num         int    `json:"num" comment:"阅读次数"`
	Remark      string `json:"remark" comment:"备注信息"`
	Status      string `json:"status" comment:"状态（0正常 1删除 2停用 3冻结）"`
	CurrAdminId int64  `json:"-" comment:""`
}

// AnnouncementGetReq 功能获取请求参数
type AnnouncementGetReq struct {
	Id int64 `uri:"id"`
}

// AnnouncementDeleteReq 功能删除请求参数
type AnnouncementDeleteReq struct {
	Ids []int64 `json:"ids"`
}
