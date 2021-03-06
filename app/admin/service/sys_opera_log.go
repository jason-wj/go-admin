package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"time"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"

	"go-admin/common/core/sdk/service"
)

type SysOperaLog struct {
	service.Service
}

func NewSysOperaLogService(s *service.Service) *SysOperaLog {
	var srv = new(SysOperaLog)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取SysOperaLog列表
func (e *SysOperaLog) GetPage(c *dto.SysOperaLogQueryReq) ([]models.SysOperaLog, int64, error) {
	var list []models.SysOperaLog
	var data models.SysOperaLog
	var count int64

	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysOperaLogService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取SysOperaLog对象
func (e *SysOperaLog) Get(id int64) (*models.SysOperaLog, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysOperaLog{}
	err := e.Orm.First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// Insert 创建SysOperaLog对象
func (e *SysOperaLog) Insert(c *dto.SysOperaLogInsertReq) error {
	if c.CurrUserId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	now := time.Now()
	data := models.SysOperaLog{}

	data.Title = c.Title
	data.BusinessType = c.BusinessType
	data.BusinessTypes = c.BusinessTypes
	data.Method = c.Method
	data.RequestMethod = c.RequestMethod
	data.OperatorType = c.OperatorType
	data.OperName = c.OperName
	data.DeptName = c.DeptName
	data.OperUrl = c.OperUrl
	data.OperIp = c.OperIp
	data.OperLocation = c.OperLocation
	data.OperParam = c.OperParam
	data.Status = c.Status
	data.OperTime = c.OperTime
	data.JsonResult = c.JsonResult
	data.LatencyTime = c.LatencyTime
	data.UserAgent = c.UserAgent
	data.Remark = c.Remark
	data.CreateBy = c.CurrUserId
	data.UpdateBy = c.CurrUserId
	data.CreatedAt = &now
	data.UpdatedAt = data.CreatedAt
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysOperaLogService Insert error:%s", err)
		return err
	}
	return nil
}

// Remove 删除SysOperaLog
func (e *SysOperaLog) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.SysOperaLog

	err = e.Orm.Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetExcel 导出OperaLog
func (e *SysOperaLog) GetExcel(list []models.SysOperaLog) ([]byte, error) {
	//sheet名称
	sheetName := "OperaLog"
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
