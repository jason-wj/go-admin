package models

import (
	"time"
)

type Category struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`

	Name      string     `json:"name" gorm:"column:name;type:varchar(255);comment:分类名称"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;type:datetime;comment:创建时间"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;type:datetime;comment:更新时间"`
	Status    string     `json:"status" gorm:"column:status;type:char(1);comment:状态（0正常 1删除 2停用 3冻结）"`
	Remark    string     `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注信息"`
	CreateBy  int64      `json:"createBy" gorm:"column:create_by;type:varchar(64);comment:创建人"`
	UpdateBy  int64      `json:"updateBy" gorm:"column:update_by;type:varchar(64);comment:更新人"`
}

func (Category) TableName() string {
	return "content_category"
}
