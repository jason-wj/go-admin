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

type Category struct {
	service.Service
}

//
// GetPage
// @Description: 获取{{.ClassName}}列表
// @receiver e
// @param c
// @param p
// @return []models.Category
// @return int64
// @return error
//
func (e *Category) GetPage(c *dto.CategoryQueryReq, p *actions.DataPermission) ([]models.Category, int64, error) {
	var list []models.Category
	var data models.Category
	var count int64

	err := e.Orm.Order("created_at desc").Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("ArticleService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

//
// Count
// @Description: 获取{{.ClassName}}对象
// @receiver e
// @param c
// @return int64
// @return error
//
func (e *Category) Count(c *dto.CategoryQueryReq) (int64, error) {
	var err error
	var count int64

	err = e.Orm.Model(&models.Category{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Limit(-1).Offset(-1).
		Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		e.Log.Errorf("CategoryService Count error:%s", err)
		return 0, err
	}
	return count, nil
}

//
// Get
// @Description: 获取{{.ClassName}}对象
// @receiver e
// @param id
// @param p
// @return *models.Category
// @return error
//
func (e *Category) Get(id int64, p *actions.DataPermission) (*models.Category, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.Category{}
	err := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

//
// QueryOne
// @Description: 通过自定义条件获取{{.ClassName}}一条记录
// @receiver e
// @param queryCondition
// @return *models.Category
// @return error
//
func (e *Category) QueryOne(queryCondition *dto.CategoryQueryReq) (*models.Category, error) {
	model := &models.Category{}
	err := e.Orm.Model(&models.Category{}).
		Scopes(
			cDto.MakeCondition(queryCondition.GetNeedSearch()),
		).First(model).Error
	if err != nil {
		e.Log.Errorf("CategoryService QueryOne error:%s", err)
		return nil, err
	}
	return model, nil
}

//
// Insert
// @Description: 创建{{.ClassName}}对象
// @receiver e
// @param c
// @return int64
// @return error
//
func (e *Category) Insert(c *dto.CategoryInsertReq) (int64, error) {
	if c.CurrUserId <= 0 {
		return 0, errors.New("参数错误")
	}

	now := time.Now()
	data := models.Category{}
	data.Name = c.Name
	data.Status = "0"
	data.Remark = c.Remark
	data.CreateBy = c.CurrUserId
	data.UpdateBy = c.CurrUserId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err := e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CategoryService Insert error:%s", err)
		return 0, err
	}
	return data.Id, nil
}

//
// Update
// @Description: 修改{{.ClassName}}对象
// @receiver e
// @param c
// @param p
// @return bool 是否有数据更新
// @return error
//
func (e *Category) Update(c *dto.CategoryUpdateReq, p *actions.DataPermission) (bool, error) {
	if c.Id <= 0 || c.CurrUserId <= 0 {
		return false, errors.New("参数错误")
	}
	data, err := e.Get(c.Id, p)
	if err != nil {
		return false, errors.New("当前数据不存在！")
	}

	//最小化变更改动过的数据
	updates := map[string]interface{}{}

	if data.Name != c.Name {
		updates["name"] = c.Name
	}

	if len(updates) > 0 {
		req := dto.CategoryQueryReq{}
		req.NameInner = c.Name
		var count int64
		count, err = e.Count(&req)
		if err != nil {
			return false, err
		}
		if count > 0 {
			return false, errors.New("该名称已被使用")
		}

		updates["updated_at"] = time.Now()
		updates["update_by"] = c.CurrUserId
		err = e.Orm.Model(&data).Where("id=?", data.Id).Updates(&updates).Error
		if err != nil {
			e.Log.Errorf("CategoryService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

//
// Remove
// @Description: 删除{{.ClassName}}
// @receiver e
// @param ids
// @param p
// @return error
//
func (e *Category) Remove(ids []int64, p *actions.DataPermission) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}

	for _, id := range ids {
		//若有文章，不得删除
		articleService := Article{}
		articleService.Orm = e.Orm
		articleService.Log = e.Log
		articleReq := dto.ArticleQueryReq{}
		articleReq.CateId = id
		var count int64
		count, err := articleService.Count(&articleReq)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New(fmt.Sprintf("分类：%d中有文章，不得删除", id))
		}
	}

	var err error
	var data models.Category

	err = e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetExcel 导出Category
func (e *Category) GetExcel(list []models.Category) ([]byte, error) {
	//sheet名称
	sheetName := "Category"
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
