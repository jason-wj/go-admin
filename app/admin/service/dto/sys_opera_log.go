package dto

import (
	"time"

	"go-admin/common/dto"
)

type SysOperaLogQueryReq struct {
	dto.Pagination `search:"-"`
	Title          string `form:"title" search:"type:contains;column:title;table:sys_opera_log" comment:"操作模块"`
	Method         string `form:"method" search:"type:contains;column:method;table:sys_opera_log" comment:"函数"`
	RequestMethod  string `form:"requestMethod" search:"type:contains;column:request_method;table:sys_opera_log" comment:"请求方式"`
	OperUrl        string `form:"operUrl" search:"type:contains;column:oper_url;table:sys_opera_log" comment:"访问地址"`
	OperIp         string `form:"operIp" search:"type:exact;column:oper_ip;table:sys_opera_log" comment:"客户端ip"`
	Status         int    `form:"status" search:"type:exact;column:status;table:sys_opera_log" comment:"状态"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:ctime;table:sys_opera_log" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:ctime;table:sys_opera_log" comment:"创建时间"`
	SysOperaLogOrder
}

type SysOperaLogOrder struct {
	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_opera_log" form:"createdAtOrder"`
}

func (m *SysOperaLogQueryReq) GetNeedSearch() interface{} {
	return *m
}

type SysOperaLogControl struct {
	ID            int64     `uri:"Id" comment:"编码"` // 编码
	Title         string    `json:"title" comment:"操作模块"`
	BusinessType  string    `json:"businessType" comment:"操作类型"`
	BusinessTypes string    `json:"businessTypes" comment:""`
	Method        string    `json:"method" comment:"函数"`
	RequestMethod string    `json:"requestMethod" comment:"请求方式"`
	OperatorType  string    `json:"operatorType" comment:"操作类型"`
	OperName      string    `json:"operName" comment:"操作者"`
	DeptName      string    `json:"deptName" comment:"部门名称"`
	OperUrl       string    `json:"operUrl" comment:"访问地址"`
	OperIp        string    `json:"operIp" comment:"客户端ip"`
	OperLocation  string    `json:"operLocation" comment:"访问位置"`
	OperParam     string    `json:"operParam" comment:"请求参数"`
	Status        string    `json:"status" comment:"操作状态"`
	OperTime      time.Time `json:"operTime" comment:"操作时间"`
	JsonResult    string    `json:"jsonResult" comment:"返回数据"`
	Remark        string    `json:"remark" comment:"备注"`
	LatencyTime   string    `json:"latencyTime" comment:"耗时"`
	UserAgent     string    `json:"userAgent" comment:"ua"`
}

type SysOperaLogGetReq struct {
	Id int64 `uri:"id"`
}

type SysOperaLogInsertReq struct {
	Id            int64      `uri:"id" `
	Title         string     `json:"title"`
	BusinessType  string     `json:"businessType"`
	BusinessTypes string     `json:"businessTypes"`
	Method        string     `json:"method"`
	RequestMethod string     `json:"requestMethod"`
	OperatorType  string     `json:"operatorType"`
	OperName      string     `json:"operName"`
	DeptName      string     `json:"deptName"`
	OperUrl       string     `json:"operUrl"`
	OperIp        string     `json:"operIp"`
	OperLocation  string     `json:"operLocation"`
	OperParam     string     `json:"operParam"`
	Status        string     `json:"status"`
	OperTime      *time.Time `json:"operTime"`
	JsonResult    string     `json:"jsonResult"`
	Remark        string     `json:"remark"`
	LatencyTime   string     `json:"latencyTime"`
	UserAgent     string     `json:"userAgent"`
	CurrUserId    int64      `json:"-" comment:""`
}

// SysOperaLogDeleteReq 功能删除请求参数
type SysOperaLogDeleteReq struct {
	Ids []int64 `json:"ids"`
}
