package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/app/plugins/content/models"
	"go-admin/app/plugins/content/service"
	"go-admin/app/plugins/content/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/sdk/api"
	_ "go-admin/common/core/sdk/pkg/response"
	"go-admin/common/middleware/auth"
)

type Announcement struct {
	api.Api
}

// GetPage 获取公告管理列表
func (e Announcement) GetPage(c *gin.Context) {
	req := dto.AnnouncementQueryReq{}
	s := service.Announcement{}
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

	p := actions.GetPermissionFromContext(c)
	list := make([]models.Announcement, 0)
	var count int64

	list, count, err = s.GetPage(&req, p)
	if err != nil {
		e.Error(500, fmt.Sprintf("获取公告管理 失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取公告管理
func (e Announcement) Get(c *gin.Context) {
	req := dto.AnnouncementGetReq{}
	s := service.Announcement{}
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

	p := actions.GetPermissionFromContext(c)
	result, err := s.Get(req.Id, p)
	if err != nil {
		e.Error(500, fmt.Sprintf("获取公告管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(result, "查询成功")
}

// Insert 创建公告管理
func (e Announcement) Insert(c *gin.Context) {
	req := dto.AnnouncementInsertReq{}
	s := service.Announcement{}
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
	err = s.Insert(&req, uid)
	if err != nil {
		e.Error(500, err.Error())
		return
	}

	e.OK(nil, "创建成功")
}

// Update 修改公告管理
func (e Announcement) Update(c *gin.Context) {
	req := dto.AnnouncementUpdateReq{}
	s := service.Announcement{}
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
	p := actions.GetPermissionFromContext(c)
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err.Error())
		return
	}
	req.CurrAdminId = uid
	b, err := s.Update(&req, p)
	if err != nil {
		e.Error(500, fmt.Sprintf("%s", err.Error()))
		return
	}
	if !b {
		e.OK(req.Id, "未修改任何信息")
		return
	}
	e.OK(req.Id, "修改成功")
}

// Delete 删除公告管理
func (e Announcement) Delete(c *gin.Context) {
	s := service.Announcement{}
	req := dto.AnnouncementDeleteReq{}
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

	p := actions.GetPermissionFromContext(c)

	err = s.Remove(req.Ids, p)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	e.OK(req.Ids, "删除成功")
}
