package models

import "gorm.io/gorm/schema"

type ActiveRecord interface {
	schema.Tabler
	SetCreateBy(createBy int64)
	SetUpdateBy(updateBy int64)
	Generate() ActiveRecord
	GetId() interface{}
}

type BaseActiveRecord interface {
	schema.Tabler
}
