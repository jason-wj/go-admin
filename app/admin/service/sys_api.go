package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"time"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/sdk/service"
	cDto "go-admin/common/dto"
)

type SysApi struct {
	service.Service
}

func NewSysApiService(s *service.Service) *SysApi {
	var srv = new(SysApi)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取SysApi列表
func (e *SysApi) GetPage(c *dto.SysApiQueryReq, p *actions.DataPermission) ([]models.SysApi, int64, error) {
	var err error
	var list []models.SysApi
	var data models.SysApi
	var count int64

	err = e.Orm.Debug().Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysApiService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取SysApi对象with id
func (e *SysApi) Get(id int64, p *actions.DataPermission) (*models.SysApi, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysApi{}
	err := e.Orm.Model(models.SysApi{}).
		Scopes(
			actions.Permission(model.TableName(), p),
		).First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// Update 修改SysApi对象
func (e *SysApi) Update(c *dto.SysApiUpdateReq, p *actions.DataPermission) (bool, error) {
	if c.Id <= 0 || c.CurrUserId <= 0 {
		return false, errors.New("参数错误")
	}
	var model = models.SysApi{}
	err := e.Orm.Debug().Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.Id).Error
	if err != nil {
		return false, errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	updates := map[string]interface{}{}

	if model.Handle != c.Handle {
		updates["handle"] = c.Handle
	}

	if model.Title != c.Title {
		updates["title"] = c.Title
	}

	if model.Path != c.Path {
		updates["path"] = c.Path
	}

	if model.Type != c.Type {
		updates["type"] = c.Type
	}

	if model.Action != c.Action {
		updates["action"] = c.Action
	}
	if len(updates) > 0 {
		updates["update_by"] = c.CurrUserId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysApi{}).Where("id=?", c.Id).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysApiService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除SysApi
func (e *SysApi) Remove(ids []int64, p *actions.DataPermission) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var data models.SysApi

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}

	return nil
}

// GetExcel 导出SysApi
func (e *SysApi) GetExcel(list []models.SysApi) ([]byte, error) {
	//sheet名称
	sheetName := "Api"
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
