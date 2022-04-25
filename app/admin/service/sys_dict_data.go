package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"go-admin/common/core/sdk"
	"go-admin/common/core/sdk/service"
	"time"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
)

type SysDictData struct {
	service.Service
}

func NewSysDictDataService(s *service.Service) *SysDictData {
	var srv = new(SysDictData)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取列表
func (e *SysDictData) GetPage(c *dto.SysDictDataSearch) ([]models.SysDictData, int64, error) {
	var list []models.SysDictData
	var data models.SysDictData
	var count int64

	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysDictDataService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取对象
func (e *SysDictData) Get(id int64) (*models.SysDictData, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysDictData{}
	err := e.Orm.Model(&models.SysDictData{}).First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// Insert 创建对象
func (e *SysDictData) Insert(c *dto.SysDictDataControl) error {
	if c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}

	var err error
	now := time.Now()
	data := models.SysDictData{}
	data.DictSort = c.DictSort
	data.DictLabel = c.DictLabel
	data.DictValue = c.DictValue
	data.DictType = c.DictType
	data.CssClass = c.CssClass
	data.ListClass = c.ListClass
	data.IsDefault = c.IsDefault
	data.Status = c.Status
	data.Default = c.Default
	data.Remark = c.Remark
	data.CreateBy = c.CurrAdminId
	data.UpdateBy = c.CurrAdminId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysDictDataService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改对象
func (e *SysDictData) Update(c *dto.SysDictDataControl) (bool, error) {
	if c.Id <= 0 || c.CurrAdminId <= 0 {
		return false, errors.New("参数错误")
	}

	var model = models.SysDictData{}
	err := e.Orm.Debug().First(&model, c.Id).Error
	if err != nil {
		return false, errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	updates := map[string]interface{}{}

	if model.DictSort != c.DictSort {
		updates["dict_sort"] = c.DictSort
	}
	/*if model.DictLabel != c.DictLabel {
		updates["dict_label"] = c.DictLabel
	}
	if model.DictValue != c.DictValue {
		updates["dict_value"] = c.DictValue
	}
	if model.DictType != c.DictType {
		updates["dict_type"] = c.DictType
	}
	if model.CssClass != c.CssClass {
		updates["css_class"] = c.CssClass
	}
	if model.IsDefault != c.IsDefault {
		updates["is_default"] = c.IsDefault
	}
	if model.Status != c.Status {
		updates["status"] = c.Status
	}
	if model.Default != c.Default {
		updates["default"] = c.Default
	}*/
	if model.Remark != c.Remark {
		updates["remark"] = c.Remark
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysDictData{}).Where("dict_code=?", c.Id).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysConfigService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除
func (e *SysDictData) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.SysDictData

	err = e.Orm.Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetAll 获取所有
func (e *SysDictData) GetAll(c *dto.SysDictDataSearch) ([]models.SysDictData, error) {
	var err error
	var list []models.SysDictData

	err = e.Orm.Model(&models.SysDictData{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Find(&list).Error
	if err != nil {
		e.Log.Errorf("SysConfigService GetAll error:%s", err)
		return nil, err
	}
	return list, nil
}

// GetValue 根据dict和key获取值
func (e *SysDictData) GetLabel(dict, value string) string {
	if dict == "" || value == "" {
		return ""
	}
	key := dict + value
	v, _ := sdk.Runtime.GetCacheAdapter().Get("", key)
	if v != "" {
		return v
	}

	var data models.SysDictData
	search := dto.SysDictDataSearch{}
	search.DictType = dict
	search.DictValue = value

	result := models.SysDictData{}

	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(search.GetNeedSearch()),
		).First(&result).Error
	if err != nil {
		e.Log.Errorf("SysConfigService GetLabel error:%s", err)
		return ""
	}
	label := result.DictLabel
	//添加缓存
	_ = sdk.Runtime.GetCacheAdapter().Set("", key, label, -1)
	return label
}

// GetExcel 导出SysDictData
func (e *SysDictData) GetExcel(list []models.SysDictData) ([]byte, error) {
	//sheet名称
	sheetName := "DictData"
	xlsx := excelize.NewFile()
	no := xlsx.NewSheet(sheetName)
	//各列间隔
	xlsx.SetColWidth(sheetName, "A", "P", 25)
	//头部描述
	xlsx.SetSheetRow(sheetName, "A1", &[]interface{}{
		"", "", "", "", "", "", "", "",
		"", "", "", "", "", "", ""})

	/*for i, item := range list {
		axis := fmt.Sprintf("A%d", i+2)

		//todo 数据导入逻辑

		//按标签对应输入数据
		xlsx.SetSheetRow(sheetName, axis, &[]interface{}{})
	}*/
	xlsx.SetActiveSheet(no)
	data, _ := xlsx.WriteToBuffer()
	return data.Bytes(), nil
}
