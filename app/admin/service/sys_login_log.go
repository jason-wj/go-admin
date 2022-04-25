package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/core/sdk/service"
	cDto "go-admin/common/dto"
)

type SysLoginLog struct {
	service.Service
}

func NewSysLoginLogService(s *service.Service) *SysLoginLog {
	var srv = new(SysLoginLog)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetSysLoginLogPage 获取SysLoginLog列表
func (e *SysLoginLog) GetPage(c *dto.SysLoginLogSearch) ([]models.SysLoginLog, int64, error) {
	var list []models.SysLoginLog
	var data models.SysLoginLog
	var count int64

	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysLoginLogService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取SysLoginLog对象
func (e *SysLoginLog) Get(id int64) (*models.SysLoginLog, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysLoginLog{}
	err := e.Orm.First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// Remove 删除SysLoginLog
func (e *SysLoginLog) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.SysLoginLog

	err = e.Orm.Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetExcel 导出Category
func (e *SysLoginLog) GetExcel(list []models.SysLoginLog) ([]byte, error) {
	//sheet名称
	sheetName := "LoginLog"
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
