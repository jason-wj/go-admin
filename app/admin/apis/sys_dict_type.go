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

type SysDictType struct {
	api.Api
}

// GetPage 字典类型列表数据
// @Summary 字典类型列表数据
// @Description 获取JSON
// @Tags 字典类型
// @Param dictName query string false "dictName"
// @Param dictId query string false "dictId"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/dict/type [get]
// @Security Bearer
func (e SysDictType) GetPage(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeSearch{}
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
	list := make([]models.SysDictType, 0)
	var count int64
	list, count, err = s.GetPage(&req)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 字典类型通过字典id获取
// @Summary 字典类型通过字典id获取
// @Description 获取JSON
// @Tags 字典类型
// @Param dictId path int true "字典类型编码"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/dict/type/{dictId} [get]
// @Security Bearer
func (e SysDictType) Get(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeById{}
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
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(result, "查询成功")
}

//Insert 字典类型创建
// @Summary 添加字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysDictTypeControl true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /admin-api/v1/dict/type [post]
// @Security Bearer
func (e SysDictType) Insert(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeControl{}
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
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf(" 创建字典类型失败，详情：%s", err.Error()))
		return
	}
	e.OK(req.Id, "创建成功")
}

// Update
// @Summary 修改字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysDictTypeControl true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /admin-api/v1/dict/type/{dictId} [put]
// @Security Bearer
func (e SysDictType) Update(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(500, err, err.Error())
		e.Logger.Error(err)
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
		e.OK(req.Id, "未修改任何信息")
		return
	}
	e.OK(req.Id, "修改成功")
}

// Delete
// @Summary 删除字典类型
// @Description 删除数据
// @Tags 字典类型
// @Param dictId path int true "dictId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /admin-api/v1/dict/type [delete]
// @Security Bearer
func (e SysDictType) Delete(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeById{}
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
	err = s.Remove(req.Ids)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(req.Id, "删除成功")
}

// GetAll
// @Summary 字典类型全部数据 代码生成使用接口
// @Description 获取JSON
// @Tags 字典类型
// @Param dictName query string false "dictName"
// @Param dictId query string false "dictId"
// @Param dictType query string false "dictType"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /admin-api/v1/dict/type-option-select [get]
// @Security Bearer
func (e SysDictType) GetAll(c *gin.Context) {
	s := service.SysDictType{}
	req := dto.SysDictTypeSearch{}
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
	list := make([]models.SysDictType, 0)
	list, err = s.GetAll(&req)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(list, "查询成功")
}
