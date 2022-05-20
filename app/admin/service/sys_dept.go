package service

import (
	"errors"
	"fmt"
	"go-admin/app/admin/models"
	"go-admin/common/actions"
	"time"

	log "go-admin/common/core/logger"
	"go-admin/common/core/sdk/pkg"

	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"

	"go-admin/common/core/sdk/service"
)

type SysDept struct {
	service.Service
}

func NewSysDeptService(s *service.Service) *SysDept {
	var srv = new(SysDept)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取SysDept列表
func (e *SysDept) GetPage(c *dto.SysDeptSearch, p *actions.DataPermission) ([]models.SysDept, int64, error) {
	var list []models.SysDept
	var err error
	var data models.SysDept
	var count int64

	err = e.Orm.Model(&data).Scopes(
		cDto.MakeCondition(c.GetNeedSearch()),
		cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		actions.Permission(data.TableName(), p),
	).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysDeptService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Get 获取SysDept对象
func (e *SysDept) Get(id int64, p *actions.DataPermission) (*models.SysDept, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysDept{}
	err := e.Orm.Model(&models.SysDept{}).
		Scopes(
			actions.Permission(model.TableName(), p),
		).First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// Insert 创建SysDept对象
func (e *SysDept) Insert(c *dto.SysDeptControl) error {
	if c.CurrUserId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	now := time.Now()
	data := models.SysDept{}
	data.DeptName = c.DeptName
	data.ParentId = c.ParentId
	data.DeptPath = c.DeptPath
	data.Sort = c.Sort
	data.Leader = c.Leader
	data.Phone = c.Phone
	data.Email = c.Email
	data.Status = c.Status
	data.CreateBy = c.CurrUserId
	data.UpdateBy = c.CurrUserId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysDeptService Insert error:%s", err)
		return err
	}
	deptPath := pkg.Int64ToString(data.DeptId) + "/"
	if data.ParentId > 0 {
		var deptP models.SysDept
		if err = tx.First(&deptP, data.ParentId).Error; err != nil {
			e.Log.Errorf("SysDeptService Insert error:%s", err)
			return err
		}
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0/" + deptPath
	}
	var mp = map[string]string{}
	mp["dept_path"] = deptPath
	if err = tx.Model(&data).Update("dept_path", deptPath).Error; err != nil {
		e.Log.Errorf("SysDeptService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改SysDept对象
func (e *SysDept) Update(c *dto.SysDeptControl) error {
	var err error
	if c.DeptId <= 0 || c.CurrUserId <= 0 {
		return errors.New("参数错误")
	}
	var model = models.SysDept{}
	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.First(&model, c.DeptId).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	deptPath := pkg.Int64ToString(model.DeptId) + "/"
	if model.ParentId != 0 {
		var deptP models.SysDept
		err = tx.First(&deptP, model.ParentId).Error
		if err != nil {
			e.Log.Errorf("SysDeptService Update error:%s", err)
			return err
		}
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0/" + deptPath
	}
	model.DeptPath = deptPath

	if c.DeptId > 0 {
		model.DeptId = c.DeptId
	}
	now := time.Now()
	model.DeptName = c.DeptName
	model.ParentId = c.ParentId
	model.Sort = c.Sort
	model.Leader = c.Leader
	model.Phone = c.Phone
	model.Email = c.Email
	model.Status = c.Status
	model.UpdatedAt = &now
	model.UpdateBy = c.CurrUserId

	err = tx.Save(&model).Error
	if err != nil {
		e.Log.Errorf("SysDeptService Update error:%s", err)
		return err
	}
	return nil
}

// Remove 删除SysDept
func (e *SysDept) Remove(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}

	var data models.SysDept

	err := e.Orm.Model(&data).Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// GetSysDeptList 获取组织数据
func (e *SysDept) getList(c *dto.SysDeptSearch) ([]models.SysDept, error) {
	var list []models.SysDept

	err := e.Orm.Model(&models.SysDept{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Find(&list).Error
	if err != nil {
		e.Log.Errorf("SysDeptService getList error:%s", err)
		return nil, err
	}
	return list, nil
}

// SetDeptTree 设置组织数据
func (e *SysDept) SetDeptTree(c *dto.SysDeptSearch) ([]dto.DeptLabel, error) {
	list, err := e.getList(c)

	m := make([]dto.DeptLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		e := dto.DeptLabel{}
		e.Id = list[i].DeptId
		e.Label = list[i].DeptName
		deptsInfo := deptTreeCall(&list, e)

		m = append(m, deptsInfo)
	}
	return m, err
}

// Call 递归构造组织数据
func deptTreeCall(deptList *[]models.SysDept, dept dto.DeptLabel) dto.DeptLabel {
	list := *deptList
	min := make([]dto.DeptLabel, 0)
	for j := 0; j < len(list); j++ {
		if dept.Id != list[j].ParentId {
			continue
		}
		mi := dto.DeptLabel{Id: list[j].DeptId, Label: list[j].DeptName, Children: []dto.DeptLabel{}}
		ms := deptTreeCall(deptList, mi)
		min = append(min, ms)
	}
	dept.Children = min
	return dept
}

// SetDeptPage 设置dept页面数据
func (e *SysDept) SetDeptPage(c *dto.SysDeptSearch) (m []models.SysDept, err error) {
	list, err := e.getList(c)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := e.deptPageCall(&list, list[i])
		m = append(m, info)
	}
	return
}

func (e *SysDept) deptPageCall(deptlist *[]models.SysDept, menu models.SysDept) models.SysDept {
	list := *deptlist
	min := make([]models.SysDept, 0)
	for j := 0; j < len(list); j++ {
		if menu.DeptId != list[j].ParentId {
			continue
		}
		mi := models.SysDept{}
		mi.DeptId = list[j].DeptId
		mi.ParentId = list[j].ParentId
		mi.DeptPath = list[j].DeptPath
		mi.DeptName = list[j].DeptName
		mi.Sort = list[j].Sort
		mi.Leader = list[j].Leader
		mi.Phone = list[j].Phone
		mi.Email = list[j].Email
		mi.Status = list[j].Status
		mi.CreatedAt = list[j].CreatedAt
		mi.Children = []models.SysDept{}
		ms := e.deptPageCall(deptlist, mi)
		min = append(min, ms)
	}
	menu.Children = min
	return menu
}

// GetRoleDeptId 获取角色的部门ID集合
func (e *SysDept) GetWithRoleId(roleId int64) ([]int64, error) {
	deptIds := make([]int64, 0)
	deptList := make([]dto.DeptIdList, 0)
	if err := e.Orm.Table("sys_role_dept").
		Select("sys_role_dept.dept_id").
		Joins("LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id").
		Where("role_id = ? ", roleId).
		Where(" sys_role_dept.dept_id not in(select sys_dept.parent_id from sys_role_dept LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id where role_id =? )", roleId).
		Find(&deptList).Error; err != nil {
		return nil, err
	}
	for i := 0; i < len(deptList); i++ {
		deptIds = append(deptIds, deptList[i].DeptId)
	}
	return deptIds, nil
}

func (e *SysDept) SetDeptLabel() (m []dto.DeptLabel, err error) {
	list := make([]models.SysDept, 0)
	err = e.Orm.Find(&list).Error
	if err != nil {
		log.Error("find dept list error, %s", err.Error())
		return
	}
	m = make([]dto.DeptLabel, 0)
	var item dto.DeptLabel
	for i := range list {
		if list[i].ParentId != 0 {
			continue
		}
		item = dto.DeptLabel{}
		item.Id = list[i].DeptId
		item.Label = list[i].DeptName
		deptInfo := deptLabelCall(&list, item)
		m = append(m, deptInfo)
	}
	return
}

// deptLabelCall
func deptLabelCall(deptList *[]models.SysDept, dept dto.DeptLabel) dto.DeptLabel {
	list := *deptList
	var mi dto.DeptLabel
	min := make([]dto.DeptLabel, 0)
	for j := 0; j < len(list); j++ {
		if dept.Id != list[j].ParentId {
			continue
		}
		mi = dto.DeptLabel{Id: list[j].DeptId, Label: list[j].DeptName, Children: []dto.DeptLabel{}}
		ms := deptLabelCall(deptList, mi)
		min = append(min, ms)
	}
	dept.Children = min
	return dept
}
