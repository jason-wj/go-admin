package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/shopspring/decimal"
	"time"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"

	"go-admin/common/core/sdk/service"
)

type SysConfig struct {
	service.Service
}

func NewSysConfigService(s *service.Service) *SysConfig {
	var srv = new(SysConfig)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取SysConfig列表
func (e *SysConfig) GetPage(c *dto.SysConfigSearch) ([]models.SysConfig, int64, error) {
	var list []models.SysConfig
	var count int64
	err := e.Orm.Scopes(
		cDto.MakeCondition(c.GetNeedSearch()),
		cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
	).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysConfigService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取SysConfig对象
func (e *SysConfig) Get(id int64) (*models.SysConfig, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysConfig{}
	err := e.Orm.First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// Insert 创建SysConfig对象
func (e *SysConfig) Insert(c *dto.SysConfigControl) error {
	if c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	now := time.Now()
	data := models.SysConfig{}
	data.ConfigName = c.ConfigName
	data.ConfigKey = c.ConfigKey
	data.ConfigValue = c.ConfigValue
	data.ConfigType = c.ConfigType
	data.IsFrontend = c.IsFrontend
	data.Remark = c.Remark
	data.CreateBy = c.CurrAdminId
	data.UpdateBy = c.CurrAdminId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysConfigService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改SysConfig对象
func (e *SysConfig) Update(c *dto.SysConfigControl) (bool, error) {
	if c.Id <= 0 || c.CurrAdminId <= 0 {
		return false, errors.New("参数错误")
	}
	var model = models.SysConfig{}
	err := e.Orm.Debug().Model(&models.SysConfig{}).First(&model, c.Id).Error
	if err != nil {
		return false, errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	updates := map[string]interface{}{}

	if model.ConfigName != c.ConfigName {
		updates["config_name"] = c.ConfigName
	}
	/*if model.ConfigKey != c.ConfigKey {
		updates["config_key"] = c.ConfigKey
	}*/
	if model.ConfigValue != c.ConfigValue {
		updates["config_value"] = c.ConfigValue
	}
	if model.ConfigType != c.ConfigType {
		updates["config_type"] = c.ConfigType
	}
	if model.IsFrontend != c.IsFrontend {
		updates["is_frontend"] = c.IsFrontend
	}
	if model.Remark != c.Remark {
		updates["remark"] = c.Remark
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysConfig{}).Where("id=?", c.Id).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysConfigService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除SysConfig
func (e *SysConfig) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.SysConfig

	err = e.Orm.Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// SetSysConfig 修改SysConfig对象
func (e *SysConfig) SetSysConfig(c *[]dto.GetSetSysConfigReq) error {
	var err error
	for _, req := range *c {
		var model = models.SysConfig{}
		err = e.Orm.Debug().Where("config_key = ?", req.ConfigKey).First(&model).Error
		if err != nil {
			return errors.New(fmt.Sprintf("无权更新该数据%s", err))
		}

		if model.Id > 0 {
			updates := map[string]interface{}{}
			if model.ConfigKey != req.ConfigKey {
				updates["config_key"] = req.ConfigKey
			}
			if model.ConfigValue != req.ConfigValue {
				updates["config_value"] = req.ConfigValue
			}

			if len(updates) > 0 {
				updates["update_by"] = req.CurrAdminId
				updates["updated_at"] = time.Now()
				err = e.Orm.Model(&models.SysConfig{}).Where("id=?", model.Id).Updates(updates).Error
				if err != nil {
					e.Log.Errorf("SysConfigService Update error:%s", err)
					return err
				}
			}
		}
	}
	return nil
}

// GetWithKey 根据Key获取SysConfig
func (e *SysConfig) GetWithKey(c *dto.SysConfigByKeyReq) (*dto.GetSysConfigByKEYForServiceResp, error) {
	var err error
	var data models.SysConfig
	resp := &dto.GetSysConfigByKEYForServiceResp{}
	err = e.Orm.Scopes().Table(data.TableName()).Where("config_key = ?", c.ConfigKey).First(resp).Error
	if err != nil {
		e.Log.Errorf("SysConfigService GetWithKey Error:%s", err)
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return resp, nil
}

// GetWithKeyStr 使用字符串类型快速获取配置结果
func (e *SysConfig) GetWithKeyStr(key, defaultVal string) string {
	query := dto.SysConfigByKeyReq{}
	query.ConfigKey = key

	resp, err := e.GetWithKey(&query)
	if err != nil || resp.ConfigValue == "" {
		return defaultVal
	}
	return resp.ConfigValue
}

func (e *SysConfig) GetWithKeyList(c *dto.SysConfigSearch) ([]models.SysConfig, error) {
	var list []models.SysConfig
	var err error
	err = e.Orm.Scopes(
		cDto.MakeCondition(c.GetNeedSearch()),
	).Find(&list).Error
	if err != nil {
		e.Log.Errorf("SysConfigService GetWithKeyList Error:%s", err)
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return list, nil
}

// GetWithKeyDecimal 返回
func (e *SysConfig) GetWithKeyDecimal(key string, defaultVal decimal.Decimal) decimal.Decimal {
	resultValue := e.GetWithKeyStr(key, "")
	if resultValue == "" {
		return defaultVal
	}
	result, _ := decimal.NewFromString(resultValue)
	return result
}

// GetExcel 导出配置
func (e *SysConfig) GetExcel(list []models.SysConfig) ([]byte, error) {
	//sheet名称
	sheetName := "config"
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
