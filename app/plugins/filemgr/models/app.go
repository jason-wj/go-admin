package models

import (
	"time"
)

type App struct {
	Id int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`

	Version      string     `json:"version" gorm:"column:version;type:varchar(100);comment:版本号"`
	Platform     string     `json:"platform" gorm:"column:platform;type:char(1);comment:平台 ( 0:安卓 1:苹果 )"`
	Type         string     `json:"type" gorm:"column:type;type:varchar(255);comment:App类型"`
	LocalAddress string     `json:"localAddress" gorm:"column:local_address;type:varchar(255);comment:本地地址"`
	BucketName   string     `json:"bucketName" gorm:"column:bucket_name;type:varchar(255);comment:bucketName，用于生成下载链接"`
	OssKey       string     `json:"ossKey" gorm:"column:oss_key;type:varchar(255);comment:用于生成下载链接"`
	DownloadNum  int        `json:"downloadNum" gorm:"column:download_num;type:decimal(10,0);comment:下载数量"`
	Status       string     `json:"status" gorm:"column:status;type:char(1);comment:状态（0正常 1删除 2停用 3冻结）"`
	Remark       string     `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注信息"`
	DownloadType string     `json:"downloadType" gorm:"column:download_type;type:char(1);comment:下载类型(0-oss 1-外链)"`
	DownloadUrl  string     `json:"downloadUrl" gorm:"column:download_url;type:varchar(255);comment:下载地址(download_type=1使用)"`
	CreateBy     int64      `json:"createBy" gorm:"column:create_by;type:varchar(64);comment:创建人"`
	UpdateBy     int64      `json:"updateBy" gorm:"column:update_by;type:varchar(64);comment:更新人"`
	CreatedAt    *time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt    *time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
}

func (App) TableName() string {
	return "filemgr_app"
}
