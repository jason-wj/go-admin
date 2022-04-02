package service

import (
	"errors"
	"fmt"
	"go-admin/app/plugins/content/models"
	"go-admin/app/plugins/content/service/dto"
	"go-admin/common/core/sdk/service"
	"gorm.io/gorm"
	"time"

	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Article struct {
	service.Service
}

func NewArticleService(s *service.Service) *Article {
	var srv = new(Article)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取Article列表
func (e *Article) GetPage(c *dto.ArticleQueryReq, p *actions.DataPermission) ([]models.Article, int64, error) {
	var list []models.Article
	var data models.Article
	var count int64

	err := e.Orm.Order("created_at desc").Model(&data).Preload("Category").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("ArticleService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取Article对象
func (e *Article) Get(id int64, p *actions.DataPermission) (*models.Article, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}

	model := &models.Article{}
	err := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// QueryOne 通过自定义条件获取一条记录
func (e *Article) QueryOne(queryCondition *dto.ArticleQueryReq) (*models.Article, error) {
	model := &models.Article{}
	err := e.Orm.Model(&models.Article{}).
		Scopes(
			cDto.MakeCondition(queryCondition.GetNeedSearch()),
		).First(model).Error
	if err != nil {
		e.Log.Errorf("ArticleService QueryOne error:%s", err)
		return nil, err
	}
	return model, nil
}

// Count 获取条数
func (e *Article) Count(c *dto.ArticleQueryReq) (int64, error) {
	var err error
	var count int64
	err = e.Orm.Model(&models.Article{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Limit(-1).Offset(-1).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		e.Log.Errorf("ArticleService Count error:%s", err)
		return 0, err
	}
	return count, nil
}

// Insert 创建Article对象
func (e *Article) Insert(c *dto.ArticleInsertReq) error {
	if c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}

	now := time.Now()
	data := models.Article{}
	data.CateId = c.CateId
	data.Content = c.Content
	data.Name = c.Name
	data.Status = "0"
	data.Remark = c.Remark
	data.CreateBy = c.CurrAdminId
	data.UpdateBy = c.CurrAdminId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err := e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ArticleService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改Article对象
func (e *Article) Update(c *dto.ArticleUpdateReq, p *actions.DataPermission) (bool, error) {
	if c.Id <= 0 || c.CurrAdminId <= 0 {
		return false, errors.New("参数错误")
	}

	data, err := e.Get(c.Id, p)
	if err != nil {
		return false, err
	}

	//最小化变更改动过的数据
	updates := map[string]interface{}{}

	if data.CateId != c.CateId {
		updates["cate_id"] = c.CateId
	}
	if data.Name != c.Name {
		updates["name"] = c.Name
	}
	if data.Status != c.Status {
		updates["status"] = c.Status
	}
	if data.Content != c.Content {
		updates["content"] = c.Content
	}
	if data.Remark != c.Remark {
		updates["remark"] = c.Remark
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		updates["update_by"] = c.CurrAdminId
		err = e.Orm.Model(&data).Where("id=?", data.Id).Updates(&updates).Error
		if err != nil {
			e.Log.Errorf("ArticleService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除Article
func (e *Article) Remove(ids []int64, p *actions.DataPermission) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.Article

	err = e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}
