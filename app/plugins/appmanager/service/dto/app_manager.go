package dto

import (
	"github.com/shopspring/decimal"
	"go-admin/common/dto"
	"time"
)

type AppManagerQueryReq struct {
	dto.Pagination `search:"-"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:app_manager" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:app_manager" comment:"创建时间"`
	Version        string `form:"version"  search:"type:exact;column:version;table:app_manager" comment:"版本号"`
	Platform       string `form:"platform"  search:"type:exact;column:platform;table:app_manager" comment:"平台 ( 0:安卓 1:苹果 )"`
	Type           string `form:"type"  search:"type:exact;column:type;table:app_manager" comment:"App类型"`
	Status         string `form:"status"  search:"type:exact;column:status;table:app_manager" comment:"app状态 0:已发布 1：待发布"`
	DownloadType   string `form:"downloadType"  search:"type:exact;column:download_type;table:app_manager" comment:"下载类型(0-oss 1-外链)"`
	DownloadUrl    string `form:"downloadUrl"  search:"type:exact;column:download_url;table:app_manager" comment:"下载地址(download_type=1使用)"`
	AppManagerOrder
}

type AppManagerOrder struct {
	IdOrder           int64           `form:"idOrder"  search:"type:order;column:id;table:app_manager"`
	VersionOrder      string          `form:"versionOrder"  search:"type:order;column:version;table:app_manager"`
	PlatformOrder     string          `form:"platformOrder"  search:"type:order;column:platform;table:app_manager"`
	TypeOrder         string          `form:"typeOrder"  search:"type:order;column:type;table:app_manager"`
	LocalAddressOrder string          `form:"localAddressOrder"  search:"type:order;column:local_address;table:app_manager"`
	BucketNameOrder   string          `form:"bucketNameOrder"  search:"type:order;column:bucket_name;table:app_manager"`
	OssKeyOrder       string          `form:"ossKeyOrder"  search:"type:order;column:oss_key;table:app_manager"`
	DownloadNumOrder  decimal.Decimal `form:"downloadNumOrder"  search:"type:order;column:download_num;table:app_manager"`
	StatusOrder       string          `form:"statusOrder"  search:"type:order;column:status;table:app_manager"`
	CreateByOrder     string          `form:"createByOrder"  search:"type:order;column:create_by;table:app_manager"`
	CreatedAtOrder    *time.Time      `form:"createdAtOrder"  search:"type:order;column:created_at;table:app_manager"`
	UpdateByOrder     string          `form:"updateByOrder"  search:"type:order;column:update_by;table:app_manager"`
	UpdatedAtOrder    *time.Time      `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:app_manager"`
	RemarkOrder       string          `form:"remarkOrder"  search:"type:order;column:remark;table:app_manager"`
	CorpCodeOrder     string          `form:"corpCodeOrder"  search:"type:order;column:corp_code;table:app_manager"`
	CorpNameOrder     string          `form:"corpNameOrder"  search:"type:order;column:corp_name;table:app_manager"`
	DownloadTypeOrder string          `form:"downloadTypeOrder"  search:"type:order;column:download_type;table:app_manager"`
	DownloadUrlOrder  string          `form:"downloadUrlOrder"  search:"type:order;column:download_url;table:app_manager"`
}

func (m *AppManagerQueryReq) GetNeedSearch() interface{} {
	return *m
}

type AppManagerInsertReq struct {
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

type AppManagerUpdateReq struct {
	Id          int64  `uri:"id" comment:"主键"` // 主键
	Status      string `json:"status" comment:"状态（0正常 1删除 2停用 3冻结）"`
	CurrAdminId int64  `json:"-" comment:""`
}

// AppManagerGetReq 功能获取请求参数
type AppManagerGetReq struct {
	Id int64 `uri:"id"`
}

// AppManagerDeleteReq 功能删除请求参数
type AppManagerDeleteReq struct {
	Ids []int64 `json:"ids"`
}
