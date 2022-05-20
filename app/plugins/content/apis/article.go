package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	adminService "go-admin/app/admin/service"
	"go-admin/app/plugins/content/models"
	"go-admin/app/plugins/content/service"
	"go-admin/app/plugins/content/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/sdk/api"
	_ "go-admin/common/core/sdk/pkg/response"
	"go-admin/common/middleware/auth"
	"go-admin/common/utils/dateUtils"
	"strconv"
	"time"
)

type Article struct {
	api.Api
}

// GetPage 获取富文本内容列表
func (e Article) GetPage(c *gin.Context) {
	req := dto.ArticleQueryReq{}
	s := service.Article{}
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
	list := make([]models.Article, 0)
	var count int64

	list, count, err = s.GetPage(&req, p)
	if err != nil {
		e.Error(500, fmt.Sprintf("获取富文本内容 失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, nil, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取富文本内容
func (e Article) Get(c *gin.Context) {
	req := dto.ArticleGetReq{}
	s := service.Article{}
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
		e.Error(500, fmt.Sprintf("获取富文本内容失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(result, "查询成功")
}

// Insert 创建富文本内容
func (e Article) Insert(c *gin.Context) {
	req := dto.ArticleInsertReq{}
	s := service.Article{}
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
	req.CurrUserId = uid
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err.Error())
		return
	}

	e.OK(nil, "创建成功")
}

// Update 修改富文本内容
func (e Article) Update(c *gin.Context) {
	req := dto.ArticleUpdateReq{}
	s := service.Article{}
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
	req.CurrUserId = uid
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

// Delete 删除富文本内容
func (e Article) Delete(c *gin.Context) {
	s := service.Article{}
	req := dto.ArticleDeleteReq{}
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

// Export 导出文章
func (e Article) Export(c *gin.Context) {
	req := dto.ArticleQueryReq{}
	s := service.Article{}
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

	var sysConfService = new(adminService.SysConfig)
	sysConfService.Orm = s.Orm
	sysConfService.Log = s.Log

	//最小导出数据量
	maxSize, err := strconv.Atoi(sysConfService.GetWithKeyStr("sys_max_export_size", "1000"))
	if err != nil {
		e.Error(500, fmt.Sprintf("配置读取异常：%s", err.Error()))
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.Article, 0)
	req.PageIndex = 1
	req.PageSize = maxSize
	list, _, err = s.GetPage(&req, p)
	if err != nil {
		e.Error(500, err.Error())
		return
	}
	data, _ := s.GetExcel(list)
	fileName := "article_" + dateUtils.ConvertToStr(time.Now(), 3) + ".xlsx"
	e.DownloadExcel(fileName, data)
}
