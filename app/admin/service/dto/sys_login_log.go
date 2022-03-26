package dto

import (
	"time"

	"go-admin/common/dto"
)

type SysLoginLogSearch struct {
	dto.Pagination `search:"-"`
	Username       string `form:"username" search:"type:exact;column:username;table:sys_login_log" comment:"用户名"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_login_log" comment:"状态"`
	Ipaddr         string `form:"ipaddr" search:"type:exact;column:ipaddr;table:sys_login_log" comment:"ip地址"`
	LoginLocation  string `form:"loginLocation" search:"type:exact;column:login_location;table:sys_login_log" comment:"归属地"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:ctime;table:sys_login_log" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:ctime;table:sys_login_log" comment:"创建时间"`
	SysLoginLogOrder
}

type SysLoginLogOrder struct {
	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_login_log" form:"createdAtOrder"`
}

func (m *SysLoginLogSearch) GetNeedSearch() interface{} {
	return *m
}

type SysLoginLogControl struct {
	ID            int64     `uri:"Id" comment:"主键"` // 主键
	Username      string    `json:"username" comment:"用户名"`
	Status        string    `json:"status" comment:"状态"`
	Ipaddr        string    `json:"ipaddr" comment:"ip地址"`
	LoginLocation string    `json:"loginLocation" comment:"归属地"`
	Browser       string    `json:"browser" comment:"浏览器"`
	Os            string    `json:"os" comment:"系统"`
	Platform      string    `json:"platform" comment:"固件"`
	LoginTime     time.Time `json:"loginTime" comment:"登录时间"`
	Remark        string    `json:"remark" comment:"备注"`
	Msg           string    `json:"msg" comment:"信息"`
}

type SysLoginLogGetReq struct {
	Id int64 `uri:"id"`
}

// SysLoginLogDeleteReq 功能删除请求参数
type SysLoginLogDeleteReq struct {
	Ids []int64 `json:"ids"`
}
