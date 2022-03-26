package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/core/sdk/api"
	_ "go-admin/common/core/sdk/pkg/response"
	"go-admin/common/middleware/auth"
)

type SysPost struct {
	api.Api
}

// GetPage 岗位列表数据
func (e SysPost) GetPage(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	list := make([]models.SysPost, 0)
	var count int64

	list, count, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取岗位信息
func (e SysPost) Get(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostGetReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	result, err := s.Get(req.Id)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("岗位信息获取失败！错误详情：%s", err.Error()))
		return
	}

	e.OK(result, "查询成功")
}

// Insert 添加岗位
func (e SysPost) Insert(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostInsertReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err, err.Error())
		return
	}
	req.CurrAdminId = uid
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("新建岗位失败！错误详情：%s", err.Error()))
		return
	}
	e.OK(req.PostId, "创建成功")
}

// Update 修改岗位
func (e SysPost) Update(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostUpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	uid, rCode, err := auth.GetUserId(c)
	if err != nil {
		e.Error(rCode, err, err.Error())
		return
	}
	req.CurrAdminId = uid

	b, err := s.Update(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("%s", err.Error()))
		return
	}
	if !b {
		e.OK(req.PostId, "未修改任何信息")
		return
	}
	e.OK(req.PostId, "修改成功")
}

// Delete 删除岗位
func (e SysPost) Delete(c *gin.Context) {
	s := service.SysPost{}
	req := dto.SysPostDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	err = s.Remove(req.Ids)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("岗位删除失败！错误详情：%s", err.Error()))
		return
	}
	e.OK(req.Ids, "删除成功")
}
