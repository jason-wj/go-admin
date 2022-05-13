package apis

import (
	"fmt"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/config"
	"go-admin/common/core/sdk/api"
	"go-admin/common/core/sdk/pkg/captcha"
	_ "go-admin/common/core/sdk/pkg/response"
	"go-admin/common/middleware/auth"
	"go-admin/common/middleware/auth/authdto"
	"net/http"
	"strconv"
)

type SysUser struct {
	api.Api
}

// GetPage
func (e SysUser) GetPage(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserQueryReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.SysUser, 0)
	var count int64

	list, count, err = s.GetPage(&req, p)
	if err != nil {
		e.Error(500, "查询失败")
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get
func (e SysUser) Get(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}
	//数据权限检查
	p := actions.GetPermissionFromContext(c)
	result, err := s.Get(req.Id, p)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err.Error())
		return
	}
	e.OK(result, "查询成功")
}

// Insert
func (e SysUser) Insert(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserInsertReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}
	// 设置创建人
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid
	err = s.Insert(&req)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	e.OK(req.UserId, "创建成功")
}

// Update
func (e SysUser) Update(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserUpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	b, err := s.Update(&req, p)
	if err != nil {
		e.Error(500, fmt.Sprintf("%s", err.Error()))
		return
	}
	if !b {
		e.OK(req.UserId, "未修改任何信息")
		return
	}
	e.OK(req.UserId, "修改成功")
}

//
//  UpdateSelfPhone
//  @Description: 更新手机号
//  @receiver e
//  @param c
//
func (e SysUser) UpdateSelfPhone(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserUpdatePhoneReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid

	b, err := s.UpdateSelfPhone(&req)
	if err != nil {
		e.Error(500, fmt.Sprintf("%s", err.Error()))
		return
	}
	if !b {
		e.OK(req.CurrAdminId, "未修改任何信息")
		return
	}
	e.OK(req.CurrAdminId, "修改成功")
}

func (e SysUser) UpdateSelfNickName(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserUpdateNickNameReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid

	b, err := s.UpdateSelfNickName(&req)
	if err != nil {
		e.Error(500, fmt.Sprintf("%s", err.Error()))
		return
	}
	if !b {
		e.OK(req.CurrAdminId, "未修改任何信息")
		return
	}
	e.OK(req.CurrAdminId, "修改成功")
}

//
//  UpdateSelfEmail
//  @Description: 更新邮箱
//  @receiver e
//  @param c
//
func (e SysUser) UpdateSelfEmail(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserUpdateEmailReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid

	b, err := s.UpdateSelfEmail(&req)
	if err != nil {
		e.Error(500, fmt.Sprintf("%s", err.Error()))
		return
	}
	if !b {
		e.OK(req.CurrAdminId, "未修改任何信息")
		return
	}
	e.OK(req.CurrAdminId, "修改成功")
}

// Delete
func (e SysUser) Delete(c *gin.Context) {
	s := service.SysUser{}
	req := dto.SysUserById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(req.Ids, p)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.Id, "删除成功")
}

// InsetAvatar
func (e SysUser) InsetAvatar(c *gin.Context) {
	s := service.SysUser{}
	req := dto.UpdateSysUserAvatarReq{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}
	// 数据权限检查
	p := actions.GetPermissionFromContext(c)
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	guid := uuid.New().String()
	path := "static/uploadfile/"
	isExist := fileutil.Exist(path)
	if !isExist {
		err = fileutil.CreateDirAll(path)
		if err != nil {
			e.Logger.Errorf("dir create error, %s", err.Error())
			e.Error(500, "")
			return
		}
	}
	filPath := path + guid + ".jpg"
	for _, file := range files {
		e.Logger.Debugf("upload avatar file: %s", file.Filename)
		// 上传文件至指定目录
		err = c.SaveUploadedFile(file, filPath)
		if err != nil {
			e.Logger.Errorf("save file error, %s", err.Error())
			e.Error(500, "")
			return
		}
	}
	req.UserId = p.UserId
	req.Avatar = "/" + filPath

	err = s.UpdateAvatar(&req, p)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(filPath, "修改成功")
}

// ResetPwd 重置用户密码
func (e SysUser) ResetPwd(c *gin.Context) {
	s := service.SysUser{}
	req := dto.ResetSysUserPwdReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	err = s.ResetPwd(&req, p)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.UserId, "更新成功")
}

// UpdatePwd
func (e SysUser) UpdatePwd(c *gin.Context) {
	s := service.SysUser{}
	req := dto.PassWord{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	err = s.UpdatePwd(uid, req.OldPassword, req.NewPassword, p)
	if err != nil {
		e.Logger.Error(err)
		e.Error(http.StatusForbidden, err.Error())
		return
	}
	e.OK(nil, "密码修改成功")
}

// GetProfile
func (e SysUser) GetProfile(c *gin.Context) {
	s := service.SysUser{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}

	user, err := s.GetProfile(uid)
	if err != nil {
		e.Logger.Errorf("get user profile error, %s", err.Error())
		e.Error(500, "获取用户信息失败")
		return
	}
	//resp := dto.UserInfoResp{}
	e.OK(user, "查询成功")
}

// GetInfo
func (e SysUser) GetInfo(c *gin.Context) {
	req := dto.SysUserById{}
	s := service.SysUser{}
	r := service.SysRole{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&r.Service).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	var roles = make([]string, 1)
	roles[0] = auth.GetRoleName(c)
	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"
	var buttons = make([]string, 1)
	buttons[0] = "*:*:*"

	var mp = make(map[string]interface{})
	mp["roles"] = roles
	if auth.GetRoleName(c) == "admin" || auth.GetRoleName(c) == "系统管理员" {
		mp["permissions"] = permissions
		mp["buttons"] = buttons
	} else {
		list, _ := r.GetById(auth.GetRoleId(c))
		mp["permissions"] = list
		mp["buttons"] = list
	}

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.Id = uid
	result, err := s.Get(req.Id, p)
	if err != nil {
		e.Error(http.StatusUnauthorized, "登录失败")
		return
	}
	mp["introduction"] = " am a super administrator"
	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if result.Avatar != "" {
		mp["avatar"] = result.Avatar
	}
	mp["userName"] = result.NickName
	mp["userId"] = result.UserId
	mp["deptId"] = result.DeptId
	mp["name"] = result.NickName
	mp["code"] = 200
	e.OK(mp, "")
}

func (e SysUser) Login(c *gin.Context) {
	req := dto.LoginReq{}
	s := service.SysUser{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(500, err.Error())
		return
	}

	if req.Code == "" || req.Password == "" || req.Username == "" {
		e.Error(401, "登录信息异常，请检查")
		return
	}

	if config.ApplicationConfig.Mode != "dev" {
		if !captcha.Verify(req.UUID, req.Code, true) {
			e.Error(500, "验证码错误")
			return
		}
	}

	userResp, err := s.GetUser(&req)
	if err != nil {
		e.Error(401, err.Error())
		return
	}

	c.Set(authdto.LoginUserId, strconv.FormatInt(userResp.UserId, 10))
	c.Set(authdto.RoleId, userResp.Role.RoleId)
	c.Set(authdto.RoleName, userResp.Role.RoleName)
	c.Set(authdto.RoleKey, userResp.Role.RoleKey)
	c.Set(authdto.UserName, userResp.Username)
	c.Set(authdto.DataScope, userResp.Role.DataScope)
	c.Set(authdto.UserInfo, userResp)
	auth.Login(c)
	s.LoginLogToDB(c, "0", "登录操作", userResp.UserId)
}

// LogOut
func (e *SysUser) LogOut(c *gin.Context) {
	s := service.SysUser{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(500, err.Error())
		return
	}

	s.LoginLogToDB(c, "2", "退出成功", 0)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})
}
