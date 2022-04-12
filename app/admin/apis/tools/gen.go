package tools

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/models/tools"
	"go-admin/app/admin/service"
	"go-admin/common/core/config"
	"go-admin/common/core/sdk/api"
	"go-admin/common/core/sdk/pkg"
	"strings"
	"text/template"
)

type Gen struct {
	api.Api
}

func (e Gen) Preview(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("tableId接收失败！错误详情：%s", err.Error()))
		return
	}
	table.TableId = id
	t1, err := template.ParseFiles("static/template/model.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("model模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t2, err := template.ParseFiles("static/template/apis.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("api模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t3, err := template.ParseFiles("static/template/js.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("js模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t4, err := template.ParseFiles("static/template/vue.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("vue模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t5, err := template.ParseFiles("static/template/router_check_role.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("路由模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t6, err := template.ParseFiles("static/template/dto.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("dto模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t7, err := template.ParseFiles("static/template/service.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("service模版读取失败！错误详情：%s", err.Error()))
		return
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(500, fmt.Sprintf("数据库链接获取失败！错误详情：%s", err.Error()))
		return
	}

	tab, _ := table.Get(db, false)
	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	var b5 bytes.Buffer
	err = t5.Execute(&b5, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)
	var b7 bytes.Buffer
	err = t7.Execute(&b7, tab)

	mp := make(map[string]interface{})
	mp["static/template/model.go.template"] = b1.String()
	mp["static/template/api.go.template"] = b2.String()
	mp["static/template/js.go.template"] = b3.String()
	mp["static/template/vue.go.template"] = b4.String()
	mp["static/template/router.go.template"] = b5.String()
	mp["static/template/dto.go.template"] = b6.String()
	mp["static/template/service.go.template"] = b7.String()
	e.OK(mp, "")
}

func (e Gen) GenCode(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	table := tools.SysTables{}
	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("tableId参数接收失败！错误详情：%s", err.Error()))
		return
	}

	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		e.Error(500, fmt.Sprintf("数据库链接获取失败！错误详情：%s", err.Error()))
		return
	}

	table.TableId = id
	tab, _ := table.Get(db, false)

	e.NOActionsGen(c, tab)

	e.OK("", "Code generated successfully！")
}

func (e Gen) NOActionsGen(c *gin.Context, tab tools.SysTables) {
	e.Context = c
	log := e.GetLogger()
	tab.MLTBName = strings.Replace(tab.TBName, "_", "-", -1)

	basePath := "static/template/"
	routerFile := basePath + "router_check_role.go.template"

	if tab.IsAuth == 2 {
		routerFile = basePath + "router_no_check_role.go.template"
	}

	t1, err := template.ParseFiles(basePath + "model.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("model模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t2, err := template.ParseFiles(basePath + "apis.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("api模版读取失败！错误详情：%s", err.Error()))
		return
	}
	t3, err := template.ParseFiles(routerFile)
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("路由模版失败！错误详情：%s", err.Error()))
		return
	}
	t4, err := template.ParseFiles(basePath + "js.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("js模版解析失败！错误详情：%s", err.Error()))
		return
	}
	t5, err := template.ParseFiles(basePath + "vue.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("vue模版解析失败！错误详情：%s", err.Error()))
		return
	}
	t6, err := template.ParseFiles(basePath + "dto.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("dto模版解析失败失败！错误详情：%s", err.Error()))
		return
	}
	t7, err := template.ParseFiles(basePath + "service.go.template")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("service模版失败！错误详情：%s", err.Error()))
		return
	}

	_ = pkg.PathCreate("./app/" + tab.PackageName + "/apis/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/models/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/router/")
	_ = pkg.PathCreate("./app/" + tab.PackageName + "/service/dto/")
	_ = pkg.PathCreate(config.GenConfig.FrontPath + "/api/" + tab.PackageName + "/")
	err = pkg.PathCreate(config.GenConfig.FrontPath + "/views/" + tab.PackageName + "/" + tab.MLTBName + "/")
	if err != nil {
		log.Error(err)
		e.Error(500, fmt.Sprintf("views目录创建失败！错误详情：%s", err.Error()))
		return
	}

	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	var b5 bytes.Buffer
	err = t5.Execute(&b5, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)
	var b7 bytes.Buffer
	err = t7.Execute(&b7, tab)
	pkg.FileCreate(b1, "./app/"+tab.PackageName+"/models/"+tab.TBName+".go")
	pkg.FileCreate(b2, "./app/"+tab.PackageName+"/apis/"+tab.TBName+".go")
	pkg.FileCreate(b3, "./app/"+tab.PackageName+"/router/"+tab.TBName+".go")
	pkg.FileCreate(b4, config.GenConfig.FrontPath+"/api/"+tab.PackageName+"/"+tab.MLTBName+".js")
	pkg.FileCreate(b5, config.GenConfig.FrontPath+"/views/"+tab.PackageName+"/"+tab.MLTBName+"/index.vue")
	pkg.FileCreate(b6, "./app/"+tab.PackageName+"/service/dto/"+tab.TBName+".go")
	pkg.FileCreate(b7, "./app/"+tab.PackageName+"/service/"+tab.TBName+".go")

}

func (e Gen) GenMenuAndApi(c *gin.Context) {
	s := service.SysMenu{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err.Error())
		return
	}

	id, err := pkg.StringToInt(c.Param("tableId"))
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, fmt.Sprintf("tableId参数解析失败！错误详情：%s", err.Error()))
		return
	}
	err = s.InsertConf(id)
	if err != nil {
		e.Error(500, fmt.Sprintf("初始化失败，错误详情：%s", err.Error()))
		return
	}

	e.OK("", "数据生成成功！")
}
