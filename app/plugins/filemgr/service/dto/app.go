package dto

import (
	"github.com/shopspring/decimal"
	"go-admin/common/dto"
	"time"
)

type AppQueryReq struct {
	dto.Pagination `search:"-"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:filemgr_app" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:filemgr_app" comment:"创建时间"`
	Version        string `form:"version"  search:"type:exact;column:version;table:filemgr_app" comment:"版本号"`
	Platform       string `form:"platform"  search:"type:exact;column:platform;table:filemgr_app" comment:"平台 ( 0:安卓 1:苹果 )"`
	Type           string `form:"type"  search:"type:exact;column:type;table:filemgr_app" comment:"App类型"`
	Status         string `form:"status"  search:"type:exact;column:status;table:filemgr_app" comment:"app状态 0:已发布 1：待发布"`
	DownloadType   string `form:"downloadType"  search:"type:exact;column:download_type;table:filemgr_app" comment:"下载类型(0-oss 1-外链)"`
	DownloadUrl    string `form:"downloadUrl"  search:"type:exact;column:download_url;table:filemgr_app" comment:"下载地址(download_type=1使用)"`
	AppOrder
}

type AppOrder struct {
	IdOrder           int64           `form:"idOrder"  search:"type:order;column:id;table:filemgr_app"`
	VersionOrder      string          `form:"versionOrder"  search:"type:order;column:version;table:filemgr_app"`
	PlatformOrder     string          `form:"platformOrder"  search:"type:order;column:platform;table:filemgr_app"`
	TypeOrder         string          `form:"typeOrder"  search:"type:order;column:type;table:filemgr_app"`
	LocalAddressOrder string          `form:"localAddressOrder"  search:"type:order;column:local_address;table:filemgr_app"`
	BucketNameOrder   string          `form:"bucketNameOrder"  search:"type:order;column:bucket_name;table:filemgr_app"`
	OssKeyOrder       string          `form:"ossKeyOrder"  search:"type:order;column:oss_key;table:filemgr_app"`
	DownloadNumOrder  decimal.Decimal `form:"downloadNumOrder"  search:"type:order;column:download_num;table:filemgr_app"`
	StatusOrder       string          `form:"statusOrder"  search:"type:order;column:status;table:filemgr_app"`
	CreateByOrder     string          `form:"createByOrder"  search:"type:order;column:create_by;table:filemgr_app"`
	CreatedAtOrder    *time.Time      `form:"createdAtOrder"  search:"type:order;column:created_at;table:filemgr_app"`
	UpdateByOrder     string          `form:"updateByOrder"  search:"type:order;column:update_by;table:filemgr_app"`
	UpdatedAtOrder    *time.Time      `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:filemgr_app"`
	RemarkOrder       string          `form:"remarkOrder"  search:"type:order;column:remark;table:filemgr_app"`
	CorpCodeOrder     string          `form:"corpCodeOrder"  search:"type:order;column:corp_code;table:filemgr_app"`
	CorpNameOrder     string          `form:"corpNameOrder"  search:"type:order;column:corp_name;table:filemgr_app"`
	DownloadTypeOrder string          `form:"downloadTypeOrder"  search:"type:order;column:download_type;table:filemgr_app"`
	DownloadUrlOrder  string          `form:"downloadUrlOrder"  search:"type:order;column:download_url;table:filemgr_app"`
}

func (m *AppQueryReq) GetNeedSearch() interface{} {
	return *m
}

type AppInsertReq struct {
	Version      string          `json:"version" comment:"版本号"`
	Platform     string          `json:"platform" comment:"平台 ( 0:安卓 1:苹果 )"`
	Type         string          `json:"type" comment:"App类型"`
	LocalAddress string          `json:"localAddress" comment:"本地地址"`
	DownloadNum  decimal.Decimal `json:"downloadNum" comment:"下载数量"`
	Remark       string          `json:"remark" comment:"备注信息"`
	DownloadType string          `json:"downloadType" comment:"下载类型(0-oss 1-外链)"`
	DownloadUrl  string          `json:"downloadUrl" comment:"下载地址(download_type=1使用)"`
	CurrAdminId  int64           `json:"-" comment:""`
}

type AppUpdateReq struct {
	Id          int64  `uri:"id" comment:"主键"` // 主键
	Status      string `json:"status" comment:"状态（0正常 1删除 2停用 3冻结）"`
	CurrAdminId int64  `json:"-" comment:""`
}

// AppGetReq 功能获取请求参数
type AppGetReq struct {
	Id int64 `uri:"id"`
}

// AppDeleteReq 功能删除请求参数
type AppDeleteReq struct {
	Ids []int64 `json:"ids"`
}
