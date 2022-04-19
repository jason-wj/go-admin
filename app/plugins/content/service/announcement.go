package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"go-admin/app/plugins/content/models"
	"go-admin/app/plugins/content/service/dto"
	"go-admin/common/core/sdk/service"
	"gorm.io/gorm"
	"time"

	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Announcement struct {
	service.Service
}

func NewAnnouncementService(s *service.Service) *Announcement {
	var srv = new(Announcement)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取Announcement列表
func (e *Announcement) GetPage(c *dto.AnnouncementQueryReq, p *actions.DataPermission) ([]models.Announcement, int64, error) {
	var list []models.Announcement
	var data models.Announcement
	var count int64

	err := e.Orm.Order("created_at desc").Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("AnnouncementService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取Announcement对象
func (e *Announcement) Get(id int64, p *actions.DataPermission) (*models.Announcement, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}

	model := &models.Announcement{}
	err := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// QueryOne 通过自定义条件获取一条记录
func (e *Announcement) QueryOne(queryCondition *dto.AnnouncementQueryReq, p *actions.DataPermission) (*models.Announcement, error) {
	model := &models.Announcement{}
	err := e.Orm.Model(&models.Announcement{}).
		Scopes(
			cDto.MakeCondition(queryCondition.GetNeedSearch()),
			actions.Permission(model.TableName(), p),
		).First(model).Error
	if err != nil {
		e.Log.Errorf("AnnouncementService QueryOne error:%s", err)
		return nil, err
	}
	return model, nil
}

// Count 获取条数
func (e *Announcement) Count(c *dto.AnnouncementQueryReq) (int64, error) {
	var err error
	var count int64
	err = e.Orm.Model(&models.Announcement{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Limit(-1).Offset(-1).
		Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		e.Log.Errorf("AnnouncementService Count error:%s", err)
		return 0, err
	}
	return count, nil
}

// Insert 创建Announcement对象
func (e *Announcement) Insert(c *dto.AnnouncementInsertReq) error {
	if c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	if c.Num < 0 {
		return errors.New("阅读次数不得小于0！")
	}

	now := time.Now()
	data := models.Announcement{}
	data.Title = c.Title
	data.Content = c.Content
	data.Num = c.Num
	data.Status = "0"
	data.Remark = c.Remark
	data.CreateBy = c.CurrAdminId
	data.UpdateBy = c.CurrAdminId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err := e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("AnnouncementService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改Announcement对象
func (e *Announcement) Update(c *dto.AnnouncementUpdateReq, p *actions.DataPermission) (bool, error) {
	if c.Id <= 0 || c.CurrAdminId <= 0 {
		return false, errors.New("参数错误")
	}
	if c.Num < 0 {
		return false, errors.New("阅读次数不得小于0！")
	}

	data, err := e.Get(c.Id, p)
	if err != nil {
		return false, err
	}

	//最小化变更改动过的数据
	update := map[string]interface{}{}

	if data.Title != c.Title {
		update["title"] = c.Title
	}
	if data.Status != c.Status {
		update["status"] = c.Status
	}
	if data.Num != c.Num {
		update["num"] = c.Num
	}
	if data.Content != c.Content {
		update["content"] = c.Content
	}
	if data.Remark != c.Remark {
		update["remark"] = c.Remark
	}

	if len(update) > 0 {
		update["updated_at"] = time.Now()
		update["update_by"] = c.CurrAdminId
		err = e.Orm.Model(&data).Where("id=?", data.Id).Updates(&update).Error
		if err != nil {
			e.Log.Errorf("SysConfigService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除Announcement
func (e *Announcement) Remove(ids []int64, p *actions.DataPermission) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}

	var err error
	var data models.Announcement

	err = e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetExcel 导出Announcement
func (e *Announcement) GetExcel(list []models.Announcement) ([]byte, error) {
	//sheet名称
	sheetName := "Announcement"
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
