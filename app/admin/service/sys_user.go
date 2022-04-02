package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common"
	"go-admin/common/core/config"
	"go-admin/common/core/sdk"
	"go-admin/common/global"
	"go-admin/common/utils/strutils"
	"gorm.io/gorm"
	"time"

	"go-admin/common/actions"
	"go-admin/common/core/sdk/pkg"
	"go-admin/common/core/sdk/service"
	cDto "go-admin/common/dto"
)

type SysUser struct {
	service.Service
}

func NewSysUserService(s *service.Service) *SysUser {
	var srv = new(SysUser)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取SysUser列表
func (e *SysUser) GetPage(c *dto.SysUserQueryReq, p *actions.DataPermission) ([]models.SysUser, int64, error) {
	var list []models.SysUser
	var data models.SysUser
	var count int64

	err := e.Orm.Debug().Preload("Dept").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("SysUserService GetPage error:%s", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Count 获取条数
func (e *SysUser) Count(c *dto.SysUserQueryReq) (int64, error) {
	var err error
	var count int64
	err = e.Orm.Model(&models.SysUser{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Limit(-1).Offset(-1).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		e.Log.Errorf("SysUserService Count error:%s", err)
		return 0, err
	}
	return count, nil
}

// FindOne 获取一条
func (e *SysUser) FindOne(c *dto.SysUserQueryReq) (*models.SysUser, error) {
	var err error
	result := &models.SysUser{}
	err = e.Orm.Model(&models.SysUser{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).First(result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		e.Log.Errorf("SysUserService FindOne error:%s", err)
		return nil, err
	}
	return result, nil
}

// Get 获取SysUser对象
func (e *SysUser) Get(id int64, p *actions.DataPermission) (*models.SysUser, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.SysUser{}

	err := e.Orm.Model(&models.SysUser{}).Debug().
		Scopes(
			actions.Permission(model.TableName(), p),
		).First(model, id).Error
	return model, err
}

// Insert 创建SysUser对象
func (e *SysUser) Insert(c *dto.SysUserInsertReq) error {
	if c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var i int64
	err = e.Orm.Model(&models.SysUser{}).Where("username = ?", c.Username).Count(&i).Error
	if err != nil {
		e.Log.Errorf("SysUserService Insert error:%s", err)
		return err
	}
	if i > 0 {
		return errors.New("用户名已存在！")
	}

	if c.Username != "" {
		query := dto.SysUserQueryReq{}
		query.Username = c.Username
		count, err := e.Count(&query)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("用户名已存在")
		}
	}
	if c.NickName != "" {
		query := dto.SysUserQueryReq{}
		query.NickName = c.NickName
		count, err := e.Count(&query)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("昵称已存在")
		}
	}
	if c.Phone != "" {
		query := dto.SysUserQueryReq{}
		query.Phone = c.Phone
		count, err := e.Count(&query)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("手机号已存在")
		}
	}
	if c.Email != "" {
		query := dto.SysUserQueryReq{}
		query.Email = c.Email
		count, err := e.Count(&query)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("邮箱已存在")
		}
	}

	now := time.Now()
	data := models.SysUser{}
	data.Username = c.Username
	data.Password = c.Password
	data.NickName = c.NickName
	data.Phone = c.Phone
	data.RoleId = c.RoleId
	data.Avatar = c.Avatar
	data.Sex = c.Sex
	data.Email = c.Email
	data.DeptId = c.DeptId
	data.PostId = c.PostId
	data.Status = c.Status
	data.Remark = c.Remark
	data.CreateBy = c.CurrAdminId
	data.UpdateBy = c.CurrAdminId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysUserService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改SysUser对象
func (e *SysUser) Update(c *dto.SysUserUpdateReq, p *actions.DataPermission) (bool, error) {
	if c.UserId <= 0 || c.CurrAdminId <= 0 {
		return false, errors.New("参数错误")
	}
	var err error
	var model models.SysUser
	err = e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.UserId).Error
	if err != nil {
		return false, errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	updates := map[string]interface{}{}

	if c.Username != "" && model.Username != c.Username {
		query := dto.SysUserQueryReq{}
		query.Username = c.Username
		result, err := e.FindOne(&query)
		if err != nil {
			return false, err
		}
		if result != nil && result.UserId != c.UserId {
			return false, errors.New("用户名已存在")
		}
		updates["username"] = c.Username
	}
	if c.NickName != "" && model.NickName != c.NickName {
		query := dto.SysUserQueryReq{}
		query.NickName = c.NickName
		result, err := e.FindOne(&query)
		if err != nil {
			return false, err
		}
		if result != nil && result.UserId != c.UserId {
			return false, errors.New("昵称已存在")
		}
		updates["nick_name"] = c.NickName
	}
	if c.Phone != "" && model.Phone != c.Phone {
		if len(c.Phone) < 6 {
			return false, errors.New(fmt.Sprintf("手机号格式异常%s", err))
		}
		query := dto.SysUserQueryReq{}
		query.Phone = c.Phone
		result, err := e.FindOne(&query)
		if err != nil {
			return false, err
		}
		if result != nil && result.UserId != c.UserId {
			return false, errors.New("手机号已存在")
		}
		updates["phone"] = c.Phone
	}
	if c.RoleId > 0 && model.RoleId != c.RoleId {
		updates["role_id"] = c.RoleId
	}
	if c.Avatar != "" && model.Avatar != c.Avatar {
		updates["avatar"] = c.Avatar
	}
	if c.Sex != "" && model.Sex != c.Sex {
		updates["sex"] = c.Sex
	}
	if c.Email != "" && model.Email != c.Email {
		if !strutils.VerifyEmailFormat(c.Email) {
			return false, errors.New(fmt.Sprintf("邮箱格式异常%s", err))
		}
		query := dto.SysUserQueryReq{}
		query.Email = c.Email
		result, err := e.FindOne(&query)
		if err != nil {
			return false, err
		}
		if result != nil && result.UserId != c.UserId {
			return false, errors.New("邮箱已存在")
		}
		updates["email"] = c.Email
	}
	if c.DeptId > 0 && model.DeptId != c.DeptId {
		updates["dept_id"] = c.DeptId
	}
	if c.PostId > 0 && model.PostId != c.PostId {
		updates["post_id"] = c.PostId
	}
	if c.Status != "" && model.Status != c.Status {
		updates["status"] = c.Status
	}
	if c.Remark != "" && model.Remark != c.Remark {
		updates["remark"] = c.Remark
	}

	if c.Password != "" && model.Password != c.Password {
		updates["password"] = c.Remark
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysUser{}).Where("user_id=?", c.UserId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysUserService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// UpdateSelfPhone 修改手机号
func (e *SysUser) UpdateSelfPhone(c *dto.SysUserUpdatePhoneReq) (bool, error) {
	if c.CurrAdminId <= 0 || len(c.Phone) < 6 {
		return false, errors.New("参数错误")
	}
	var err error
	u, err := e.Get(c.CurrAdminId, nil)
	if err != nil {
		return false, err
	}

	updates := map[string]interface{}{}
	if c.Phone != "" && u.Phone != c.Phone {
		if len(c.Phone) < 6 {
			return false, errors.New(fmt.Sprintf("手机号格式异常%s", err))
		}
		query := dto.SysUserQueryReq{}
		query.Phone = c.Phone
		result, err := e.FindOne(&query)
		if err != nil {
			return false, err
		}
		if result != nil && result.UserId != c.CurrAdminId {
			return false, errors.New("手机号已存在")
		}
		updates["phone"] = c.Phone
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysUser{}).Where("user_id=?", c.CurrAdminId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysUserService UpdateSelfPhone error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// UpdateSelfNickName 更新昵称
func (e *SysUser) UpdateSelfNickName(c *dto.SysUserUpdateNickNameReq) (bool, error) {
	if c.CurrAdminId <= 0 || c.NickName == "" {
		return false, errors.New("参数错误")
	}
	var err error
	u, err := e.Get(c.CurrAdminId, nil)
	if err != nil {
		return false, err
	}

	updates := map[string]interface{}{}
	if c.NickName != "" && u.NickName != c.NickName {
		query := dto.SysUserQueryReq{}
		query.NickName = c.NickName
		result, err := e.FindOne(&query)
		if err != nil {
			return false, err
		}
		if result != nil && result.UserId != c.CurrAdminId {
			return false, errors.New("昵称已存在")
		}
		updates["nick_name"] = c.NickName
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysUser{}).Where("user_id=?", c.CurrAdminId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysUserService UpdateSelfNickName error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// UpdateSelfEmail 修改邮箱号
func (e *SysUser) UpdateSelfEmail(c *dto.SysUserUpdateEmailReq) (bool, error) {
	if c.CurrAdminId <= 0 || !strutils.VerifyEmailFormat(c.Email) {
		return false, errors.New("邮箱格式错误")
	}
	var err error
	u, err := e.Get(c.CurrAdminId, nil)
	if err != nil {
		return false, err
	}

	updates := map[string]interface{}{}
	if c.Email != "" && u.Email != c.Email {
		query := dto.SysUserQueryReq{}
		query.Email = c.Email
		result, err := e.FindOne(&query)
		if err != nil {
			return false, err
		}
		if result != nil && result.UserId != c.CurrAdminId {
			return false, errors.New("邮箱已存在")
		}
		updates["email"] = c.Email
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysUser{}).Where("user_id=?", c.CurrAdminId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysUserService UpdateSelfEmail error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// UpdateAvatar 更新用户头像
func (e *SysUser) UpdateAvatar(c *dto.UpdateSysUserAvatarReq, p *actions.DataPermission) error {
	if c.UserId <= 0 || c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var model models.SysUser
	err = e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.UserId).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	updates := map[string]interface{}{}
	if model.Avatar != c.Avatar {
		updates["avatar"] = c.Avatar
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysConfig{}).Where("user_id=?", c.UserId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysUserService UpdateAvatar error:%s", err)
			return err
		}
	}
	return nil
}

// UpdateStatus 更新用户状态
func (e *SysUser) UpdateStatus(c *dto.UpdateSysUserStatusReq, p *actions.DataPermission) error {
	if c.UserId <= 0 || c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var model models.SysUser
	err = e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.UserId).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	updates := map[string]interface{}{}
	if model.Status != c.Status {
		updates["status"] = c.Status
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysConfig{}).Where("user_id=?", c.UserId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysUserService UpdateStatus error:%s", err)
			return err
		}
	}
	return nil
}

// ResetPwd 重置用户密码
func (e *SysUser) ResetPwd(c *dto.ResetSysUserPwdReq, p *actions.DataPermission) error {
	if c.UserId <= 0 || c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var model models.SysUser
	err = e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.UserId).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}

	updates := map[string]interface{}{}
	if model.Password != c.Password {
		updates["password"] = c.Password
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.SysConfig{}).Where("user_id=?", c.UserId).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysUserService ResetPwd error:%s", err)
			return err
		}
	}
	return nil
}

// Remove 删除SysUser
func (e *SysUser) Remove(ids []int64, p *actions.DataPermission) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	var err error
	var data models.SysUser

	err = e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

// UpdatePwd 修改SysUser对象密码
func (e *SysUser) UpdatePwd(id int64, oldPassword, newPassword string, p *actions.DataPermission) error {
	var err error

	if newPassword == "" {
		return nil
	}
	c := &models.SysUser{}

	err = e.Orm.Model(c).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Select("UserId", "Password", "Salt").
		First(c, id).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权更新该数据%s", err))
	}
	if !pkg.CompareHashAndPassword(c.Password, oldPassword) {
		return errors.New("密码输入错误")
	}
	c.Password = newPassword
	err = e.Orm.Model(c).Where("user_id = ?", id).Select("Password", "Salt").Updates(c).Error
	if err != nil {
		e.Log.Errorf("SysUserService UpdatePwd error:%s", err)
		return err
	}
	return nil
}

func (e *SysUser) GetProfile(userId int64) (*models.SysUser, error) {
	user := &models.SysUser{}
	err := e.Orm.Preload("Dept").Preload("Post").Preload("Role").First(user, userId).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (e *SysUser) GetUser(login *dto.LoginReq) (*models.SysUser, error) {
	user := &models.SysUser{}
	err := e.Orm.Preload("Dept").Preload("Post").Preload("Role").Where("username = ?  and status = 2", login.Username).First(user).Error
	if err != nil {
		e.Log.Errorf("SysUserService GetUser error:%s", err)
		return nil, err
	}
	if !pkg.CompareHashAndPassword(user.Password, login.Password) {
		return nil, errors.New("密码输入错误")
	}
	return user, nil
}

// LoginLogToDB Write log to database
func (e *SysUser) LoginLogToDB(c *gin.Context, status string, msg string, userId int64) {
	if !config.LoggerConfig.EnabledDB {
		return
	}
	l := make(map[string]interface{})

	ua := user_agent.New(c.Request.UserAgent())
	l["ipaddr"] = common.GetClientIP(c)
	//用于定位ip所在城市
	//fmt.Println("gaConfig.ExtConfig.AMap.Key", config.ApplicationConfig.AmpKey)
	l["loginLocation"] = pkg.GetLocation(common.GetClientIP(c), config.ApplicationConfig.AmpKey)
	l["loginTime"] = pkg.GetCurrentTime()
	l["status"] = status
	l["remark"] = c.Request.UserAgent()
	browserName, browserVersion := ua.Browser()
	l["browser"] = browserName + " " + browserVersion
	l["os"] = ua.OS()
	l["platform"] = ua.Platform()
	l["username"] = userId
	l["msg"] = msg

	q := sdk.Runtime.GetMemoryQueue(c.Request.Host)
	message, err := sdk.Runtime.GetStreamMessage("", global.LoginLog, l)
	if err != nil {
		e.Log.Errorf("SysUserService LoginLogToDB error:%s", err)
		//日志报错错误，不中断请求
	} else {
		err = q.Append(message)
		if e.Log != nil {
			e.Log.Errorf("SysUserService LoginLogToDB error:%s", err)
		}
	}
}
