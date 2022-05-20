package tools

import (
	"time"

	"gorm.io/gorm"
)

type SysTables struct {
	TableId        int          `gorm:"primaryKey;autoIncrement" json:"tableId"`      //表编码
	TBName         string       `gorm:"column:table_name;size:255;" json:"tableName"` //表名称
	MLTBName       string       `gorm:"-" json:"-"`                                   //表名称
	TableComment   string       `gorm:"size:255;" json:"tableComment"`                //表备注
	ClassName      string       `gorm:"size:255;" json:"className"`                   //类名
	TplCategory    string       `gorm:"size:255;" json:"tplCategory"`                 //
	PackageName    string       `gorm:"size:255;" json:"packageName"`                 //包名
	ModuleName     string       `gorm:"size:255;" json:"moduleName"`                  //go文件名
	BusinessName   string       `gorm:"size:255;" json:"businessName"`                //
	FunctionName   string       `gorm:"size:255;" json:"functionName"`                //功能名称
	FunctionAuthor string       `gorm:"size:255;" json:"functionAuthor"`              //功能作者
	PkColumn       string       `gorm:"size:255;" json:"pkColumn"`
	PkGoField      string       `gorm:"size:255;" json:"pkGoField"`
	PkJsonField    string       `gorm:"size:255;" json:"pkJsonField"`
	Options        string       `gorm:"size:255;" json:"options"`
	IsPlugin       string       `gorm:"size:1;default:1;comment:是否插件 0-否 1-是" json:"isPlugin"`
	Remark         string       `gorm:"size:255;" json:"remark"`
	DataScope      string       `gorm:"-" json:"dataScope"`
	Columns        []SysColumns `gorm:"-" json:"columns"`
	CreateBy       int          `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy       int          `json:"updateBy" gorm:"index;comment:更新者"`
	CreatedAt      *time.Time   `json:"createdAt"`
	UpdatedAt      *time.Time   `json:"updatedAt"`
}

func (SysTables) TableName() string {
	return "sys_tables"
}

func (e *SysTables) GetPage(tx *gorm.DB, pageSize int, pageIndex int) ([]SysTables, int, error) {
	var doc []SysTables

	table := tx.Table("sys_tables")

	if e.TBName != "" {
		table = table.Where("table_name = ?", e.TBName)
	}
	if e.TableComment != "" {
		table = table.Where("table_comment = ?", e.TableComment)
	}

	var count int64

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return doc, int(count), nil
}

func (e *SysTables) Get(tx *gorm.DB, exclude bool) (SysTables, error) {
	var doc SysTables
	var err error
	table := tx.Table("sys_tables")

	if e.TBName != "" {
		table = table.Where("table_name = ?", e.TBName)
	}
	if e.TableId != 0 {
		table = table.Where("table_id = ?", e.TableId)
	}
	if e.TableComment != "" {
		table = table.Where("table_comment = ?", e.TableComment)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	var col SysColumns
	col.TableId = doc.TableId
	if doc.Columns, err = col.GetList(tx, exclude); err != nil {
		return doc, err
	}

	return doc, nil
}

func (e *SysTables) GetTree(tx *gorm.DB) ([]SysTables, error) {
	var doc []SysTables
	var err error
	table := tx.Table("sys_tables")

	if e.TBName != "" {
		table = table.Where("table_name = ?", e.TBName)
	}
	if e.TableId != 0 {
		table = table.Where("table_id = ?", e.TableId)
	}
	if e.TableComment != "" {
		table = table.Where("table_comment = ?", e.TableComment)
	}

	if err := table.Find(&doc).Error; err != nil {
		return doc, err
	}
	for i := 0; i < len(doc); i++ {
		var col SysColumns
		//col.FkCol = append(col.FkCol, SysColumns{ColumnId: 0, ColumnName: "请选择"})
		col.TableId = doc[i].TableId
		if doc[i].Columns, err = col.GetList(tx, false); err != nil {
			return doc, err
		}

	}

	return doc, nil
}

func (e *SysTables) Create(tx *gorm.DB) (SysTables, error) {
	var doc SysTables
	e.CreateBy = 0
	result := tx.Table("sys_tables").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	for i := 0; i < len(e.Columns); i++ {
		e.Columns[i].TableId = doc.TableId

		_, _ = e.Columns[i].Create(tx)
	}

	return doc, nil
}

func (e *SysTables) Update(tx *gorm.DB) (update SysTables, err error) {
	//if err = orm.Eloquent.Table("sys_tables").First(&update, e.TableId).Error; err != nil {
	//	return
	//}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	e.UpdateBy = 0
	if err = tx.Table("sys_tables").Where("table_id = ?", e.TableId).Updates(&e).Error; err != nil {
		return
	}
	for i := 0; i < len(e.Columns); i++ {
		_, _ = e.Columns[i].Update(tx)
	}
	return
}

func (e *SysTables) Delete(db *gorm.DB) (success bool, err error) {
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if err = tx.Table("sys_tables").Delete(SysTables{}, "table_id = ?", e.TableId).Error; err != nil {
		success = false
		return
	}
	if err = tx.Table("sys_columns").Delete(SysColumns{}, "table_id = ?", e.TableId).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *SysTables) BatchDelete(tx *gorm.DB, id []int) (Result bool, err error) {
	if err = tx.Unscoped().Table(e.TableName()).Where(" table_id in (?)", id).Delete(&SysColumns{}).Error; err != nil {
		return
	}
	Result = true
	return
}
