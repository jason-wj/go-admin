package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/app/plugins/appmanager/models"
	"go-admin/app/plugins/appmanager/service"
	"go-admin/app/plugins/appmanager/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/sdk/api"
	"go-admin/common/middleware/auth"
	"mime/multipart"
)

type AppManager struct {
	api.Api
}

// GetPage 获取App管理列表
// @Summary 获取App管理列表
// @Description 获取App管理列表
// @Tags App管理
// @Param version query string false "版本号"
// @Param platform query string false "平台 ( 0:安卓 1:苹果 )"
// @Param type query string false "App类型"
// @Param downloadType query string false "下载类型(0-oss 1-外链)"
// @Param downloadUrl query string false "下载地址(download_type=1使用)"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.AppManager}} "{"code": 200, "data": [...]}"
// @Router /api/v1/app-manager [get]
// @Security Bearer
func (e AppManager) GetPage(c *gin.Context) {
	req := dto.AppManagerQueryReq{}
	s := service.AppManager{}
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
	list := make([]models.AppManager, 0)
	var count int64

	list, count, err = s.GetPage(&req, p)
	if err != nil {
		e.Error(500, fmt.Sprintf("获取App管理 失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取App管理
// @Summary 获取App管理
// @Description 获取App管理
// @Tags App管理
// @Param id path string false "id"
// @Success 200 {object} response.Response{data=models.AppManager} "{"code": 200, "data": [...]}"
// @Router /api/v1/app-manager/{id} [get]
// @Security Bearer
func (e AppManager) Get(c *gin.Context) {
	req := dto.AppManagerGetReq{}
	s := service.AppManager{}
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
		e.Error(500, fmt.Sprintf("获取App管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(result, "查询成功")
}

// Insert 创建App管理
// @Summary 创建App管理
// @Description 创建App管理
// @Tags App管理
// @Accept application/json
// @Product application/json
// @Param data body dto.AppManagerInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/app-manager [post]
// @Security Bearer
func (e AppManager) Insert(c *gin.Context) {
	req := dto.AppManagerInsertReq{}
	s := service.AppManager{}
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
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}

	e.OK(nil, "创建成功")
}

// Delete 删除App管理
// @Summary 删除App管理
// @Description 删除App管理
// @Tags App管理
// @Param id body int false "id"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/app-manager [delete]
// @Security Bearer
func (e AppManager) Delete(c *gin.Context) {
	s := service.AppManager{}
	req := dto.AppManagerDeleteReq{}
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

// Upload 上传App文件
// @Summary 上传App文件
// @Description 上传App文件
// @Tags App管理
// @Router /api/v1/app-manager/upload [post]
// @Security Bearer
func (e AppManager) Upload(c *gin.Context) {
	//初始化
	s := service.AppManager{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	form, err := e.Context.MultipartForm()
	if err != nil {
		e.Error(500, fmt.Sprintf("上传文件异常：%s", err.Error()))
		return
	}

	//获取上传文件信息
	var filePath string
	file := &multipart.FileHeader{}

	err = s.GetSingleUploadFileInfo(form, file, &filePath)
	if err != nil {
		e.Error(500, fmt.Sprintf("获取文件信息异常：%s", err.Error()))
		return
	}

	//保存上传文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		e.Error(500, "上传失败")
		return
	}
	e.OK(filePath, "上传成功")
}

// Update 修改App信息
// @Summary 修改App信息
// @Description 修改App信息
// @Tags App管理
// @Accept application/json
// @Product application/json
// @Param data body dto.FullSlProductUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /admin-api/v1/app-manager/{id} [put]
// @Security Bearer
func (e AppManager) Update(c *gin.Context) {
	req := dto.AppManagerUpdateReq{}
	s := service.AppManager{}
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
