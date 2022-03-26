package service

import (
	"errors"
	"fmt"
	"time"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/core/sdk/service"
	cDto "go-admin/common/dto"
)

type SysPost struct {
	service.Service
}

func NewSysPostService(s *service.Service) *SysPost {
	var srv = new(SysPost)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取SysPost列表
func (e *SysPost) GetPage(c *dto.SysPostPageReq) ([]models.SysPost, int64, error) {
	var list []models.SysPost
	var data models.SysPost
	var count int64

	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysPostService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取SysPost对象
func (e *SysPost) Get(id int64) (*models.SysPost, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysPost{}
	err := e.Orm.First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// Insert 创建SysPost对象
func (e *SysPost) Insert(c *dto.SysPostInsertReq) error {
	if c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	now := time.Now()
	data := models.SysPost{}
	data.PostName = c.PostName
	data.PostCode = c.PostCode
	data.Sort = c.Sort
	data.Status = c.Status
	data.Remark = c.Remark
	data.CreateBy = c.CurrAdminId
	data.UpdateBy = c.CurrAdminId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysPostService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改SysPost对象
func (e *SysPost) Update(c *dto.SysPostUpdateReq) (bool, error) {
	if c.PostId <= 0 || c.CurrAdminId <= 0 {
		return false, errors.New("参数错误")
	}
	var model = models.SysPost{}
	err := e.Orm.Debug().Model(&models.SysPost{}).First(&model, c.PostId).Error
	if err != nil {
		return false, errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	model.Remark = c.Remark

	updates := map[string]interface{}{}

	if model.PostName != c.PostName {
		updates["post_name"] = c.PostName
	}
	if model.PostCode != c.PostCode {
		updates["post_code"] = c.PostCode
	}
	if model.Sort != c.Sort {
		updates["sort"] = c.Sort
	}
	if model.Status != c.Status {
		updates["status"] = c.Status
	}
	if model.Remark != c.Remark {
		updates["remark"] = c.Remark
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysConfig{}).Where("post_id=?", c.PostId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysPostService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除SysPost
func (e *SysPost) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.SysPost

	err = e.Orm.Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}
