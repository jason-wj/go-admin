package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/core/sdk/service"
	cDto "go-admin/common/dto"
	"time"
)

type SysDictType struct {
	service.Service
}

func NewSysDictTypeService(s *service.Service) *SysDictType {
	var srv = new(SysDictType)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取列表
func (e *SysDictType) GetPage(c *dto.SysDictTypeSearch) ([]models.SysDictType, int64, error) {
	var list []models.SysDictType
	var data models.SysDictType
	var count int64
	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysDictTypeService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取对象
func (e *SysDictType) Get(id int64) (*models.SysDictType, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysDictType{}
	err := e.Orm.First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// Insert 创建对象
func (e *SysDictType) Insert(c *dto.SysDictTypeControl) error {
	if c.CurrUserId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	now := time.Now()
	data := models.SysDictType{}
	data.DictName = c.DictName
	data.DictType = c.DictType
	data.Status = c.Status
	data.Remark = c.Remark
	data.CreateBy = c.CurrUserId
	data.UpdateBy = c.CurrUserId
	data.CreatedAt = &now
	data.UpdatedAt = &now

	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysDictTypeService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改对象
func (e *SysDictType) Update(c *dto.SysDictTypeControl) (bool, error) {
	if c.DictId <= 0 || c.CurrUserId <= 0 {
		return false, errors.New("参数错误")
	}

	var model = models.SysDictType{}
	err := e.Orm.Debug().First(&model, c.DictId).Error
	if err != nil {
		return false, errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	updates := map[string]interface{}{}

	if model.DictName != c.DictName {
		updates["dict_name"] = c.DictName
	}
	//不得修改类型
	/*if model.DictType != c.DictType {
		updates["dict_type"] = c.DictType
	}*/
	if model.Status != c.Status {
		updates["status"] = c.Status
	}
	if model.Remark != c.Remark {
		updates["remark"] = c.Remark
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrUserId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysDictType{}).Where("dict_id=?", c.DictId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysDictTypeService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除
func (e *SysDictType) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.SysDictType

	err = e.Orm.Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetAll 获取所有
func (e *SysDictType) GetAll(c *dto.SysDictTypeSearch) ([]models.SysDictType, error) {
	var err error
	var data models.SysDictType
	var list []models.SysDictType

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Find(&list).Error
	if err != nil {
		e.Log.Errorf("SysDictTypeService GetAll error:%s", err)
		return nil, err
	}
	return list, nil
}

// GetExcel 导出Category
func (e *SysDictType) GetExcel(list []models.SysDictType) ([]byte, error) {
	//sheet名称
	sheetName := "DictType"
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
