package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/google/uuid"
	adminService "go-admin/app/admin/service"
	"go-admin/app/plugins/filemgr/models"
	"go-admin/app/plugins/filemgr/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/config"
	"go-admin/common/core/sdk/service"
	cDto "go-admin/common/dto"
	"go-admin/common/utils/fileUtils/ossUtils"
	"gorm.io/gorm"
	"mime/multipart"
	"path"
	"time"
)

const (
	//app下载方式
	AppDownloadTypeOss = "0" //OSS
	AppDownloadTypeUrl = "1" //外链

	//app状态
	AppStatusPublish     = "0" //已发布
	AppStatusWaitPublish = "1" //待发布

	AppPlatformAndroid = "0" //安卓
	AppPlatformIOS     = "1" //IOS
)

type App struct {
	service.Service
}

func NewAppService(s *service.Service) *App {
	var srv = new(App)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取App列表
func (e *App) GetPage(c *dto.AppQueryReq, p *actions.DataPermission) ([]models.App, int64, error) {

	var list []models.App
	var data models.App
	var count int64

	err := e.Orm.Order("created_at desc").Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("AppService GetPage error:%s", err)
		return nil, 0, err
	}

	for index, item := range list {
		if item.DownloadUrl != "" {
			continue
		}
		url, err := e.generateAppOssUrl(&item)
		if err != nil {
			return nil, 0, err
		}
		item.DownloadUrl = url
		list[index] = item
	}
	return list, count, nil
}

// Get 获取App对象
func (e *App) Get(id int64, p *actions.DataPermission) (*models.App, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.App{}
	err := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// QueryOne 通过自定义条件获取一条记录
func (e *App) QueryOne(queryCondition *dto.AppQueryReq) (*models.App, error) {
	model := &models.App{}
	err := e.Orm.Model(&models.App{}).
		Scopes(
			cDto.MakeCondition(queryCondition.GetNeedSearch()),
		).First(model).Error
	if err != nil {
		e.Log.Errorf("AppService QueryOne error:%s", err)
		return nil, err
	}
	return model, nil
}

// Count 获取条数
func (e *App) Count(c *dto.AppQueryReq) (int64, error) {
	var err error
	var count int64
	err = e.Orm.Model(&models.App{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Limit(-1).Offset(-1).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		e.Log.Errorf("AppService Count error:%s", err)
		return 0, err
	}
	return count, nil
}

// Insert 创建App对象
func (e *App) Insert(c *dto.AppInsertReq) (int64, error) {
	if c.CurrUserId <= 0 {
		return 0, errors.New("参数错误")
	}
	if c.Platform == "" {
		return 0, errors.New("请选择一个平台")
	}
	if c.Version == "" {
		return 0, errors.New("请输入版本号")
	}
	if c.DownloadType == "" {
		return 0, errors.New("请选择下载类型")
	}
	if c.DownloadType == AppDownloadTypeOss {
		if c.Type == "" {
			return 0, errors.New("请选择App类型")
		}
		if c.LocalAddress == "" {
			return 0, errors.New("请上传App文件")
		}
	}
	if c.Remark == "" {
		return 0, errors.New("更新内容不得为空")
	}
	query := dto.AppQueryReq{}
	query.Platform = c.Platform
	query.Type = c.Type
	query.Version = c.Version
	count, err := e.Count(&query)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("版本已存在,请检查后重新录入")
	}

	ossKey := ""
	buckName := ""

	//oss上传
	if c.DownloadType == AppDownloadTypeOss {
		result, err := e.uploadOssFile(c.Type, c.Version, c.Platform, c.LocalAddress)
		if err != nil {
			return 0, err
		}
		ossKey = result.OssKey
		buckName = result.BucketName
	}

	now := time.Now()
	data := models.App{}
	data.Version = c.Version
	data.Platform = c.Platform
	data.Type = c.Type
	data.LocalAddress = c.LocalAddress
	data.BucketName = buckName
	data.OssKey = ossKey
	data.DownloadNum = 0
	data.Status = "0"
	data.DownloadType = c.DownloadType
	data.DownloadUrl = c.DownloadUrl
	data.Remark = c.Remark
	data.CreateBy = c.CurrUserId
	data.UpdateBy = c.CurrUserId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysConfigService Insert error:%s", err)
		return 0, err
	}
	return 0, nil
}

// Update 修改App对象
func (e *App) Update(c *dto.AppUpdateReq, p *actions.DataPermission) (bool, error) {
	if c.Id <= 0 || c.CurrUserId <= 0 {
		return false, errors.New("参数错误")
	}
	if c.Status == "" {
		return false, errors.New("请设置app状态")
	}
	data, err := e.Get(c.Id, p)
	if err != nil {
		return false, err
	}

	//变更改动过的数据
	updates := map[string]interface{}{}

	if data.Status != c.Status {
		updates["status"] = c.Status
	}

	if len(updates) > 0 {
		updates["update_by"] = c.CurrUserId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.App{}).Where("id=?", c.Id).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysConfigService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除App
func (e *App) Remove(ids []int64, p *actions.DataPermission) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	for _, id := range ids {
		result, err := e.Get(id, p)
		if err != nil {
			return err
		}

		//同一个完全相同的版本，可能因为网路有多条记录，但这些记录都指向一个oss资源，此时只有最后一条记录，才能删除oss资源
		query := dto.AppQueryReq{}
		query.Platform = result.Platform
		query.Type = result.Type
		query.Version = result.Version
		var count int64
		count, err = e.Count(&query)
		if err != nil {
			return err
		}
		if count <= 1 {
			//oss删除对应资源,无论删除成功与否
			objectKey, _ := e.generateAppOssObjectKey(result)
			oss, _ := e.getOssClient()
			if oss != nil && objectKey != "" {
				_ = oss.Bucket.DeleteObject(objectKey, nil)
			}
		}
	}

	//删除记录
	var data models.App
	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, ids).Error
	if err != nil {
		return errors.New(fmt.Sprintf("无权删除该数据%s", err))
	}
	return nil
}

//
//  GetSingleUploadFileInfo
//  @Description: 获取单个上传文件信息
//  @receiver e
//  @param form
//  @param file
//  @param dst
//  @return error
//
func (e *App) GetSingleUploadFileInfo(form *multipart.Form, file *multipart.FileHeader, dst *string) error {
	if len(form.File) != 1 {
		return errors.New("每次仅可上传一个文件")

	}
	for _, files := range form.File {
		if len(files) != 1 {
			return errors.New("每次仅可上传一个文件")
		}
		for _, item := range files {
			*dst = config.ApplicationConfig.FileRootPath + "app/" + uuid.New().String() + path.Ext(item.Filename)
			*file = *item
			return nil
		}
	}
	return nil
}

//
//  uploadOssFile
//  @Description:  上传App文件
//  @receiver e
//  @param appType
//  @param version
//  @param platform
//  @param localAddress
//  @return *models.App
//  @return error
//
func (e *App) uploadOssFile(appType, version, platform, localAddress string) (*models.App, error) {
	App := models.App{}
	App.Type = appType
	App.Version = version
	App.Platform = platform
	key, err := e.generateAppOssObjectKey(&App)
	if err != nil {
		return nil, err
	}
	client, err := e.getOssClient()
	if err != nil {
		return nil, err
	}
	err = client.UploadWithSpace(key, localAddress)
	if err != nil {
		e.Log.Errorf("上传失败，失败原因:%s \r\n", err)
		return nil, err
	}
	App.BucketName = client.BucketName
	App.OssKey = key
	return &App, nil
}

//
//  generateAppOssUrl
//  @Description: 获取app下载链接
//  @receiver e
//  @param App
//  @return string
//  @return error
//
func (e *App) generateAppOssUrl(App *models.App) (string, error) {
	appPath, err := e.generateAppOssObjectKey(App)
	if err != nil {
		return "", err
	}
	oss, err := e.getOssClient()
	if err != nil {
		return "", err
	}
	return oss.GeneratePresignedUrl(appPath)
}

//
//  getOssClient
//  @Description: 获取oss客户端
//  @receiver e
//  @param App
//  @return *ossUtils.ALiYunOSS
//  @return error
//
func (e *App) getOssClient() (*ossUtils.ALiYunOSS, error) {
	var sysConfService = adminService.NewSysConfigService(&e.Service)
	endPoint := sysConfService.GetWithKeyStr("app_oss_endpoint", "")
	key := sysConfService.GetWithKeyStr("app_oss_access_key_id", "")
	secret := sysConfService.GetWithKeyStr("app_oss_access_key_secret", "")
	bucketName := sysConfService.GetWithKeyStr("app_oss_bucket", "")
	oss := ossUtils.ALiYunOSS{}
	err := oss.InitOssClient(key, secret, endPoint, bucketName)
	if err != nil {
		return nil, err
	}
	return &oss, nil
}

//
//  generateAppOssObjectKey
//  @Description: 生成oss key
//  @receiver e
//  @param App
//  @return string
//  @return error
//
func (e *App) generateAppOssObjectKey(App *models.App) (string, error) {
	var sysConfService = adminService.NewSysConfigService(&e.Service)

	//app目录
	appPath := sysConfService.GetWithKeyStr("app_oss_root_path", "")
	appPath += App.Type
	appPath += "_" + App.Version

	switch App.Platform {
	case AppPlatformAndroid:
		appPath += ".apk"
	case AppPlatformIOS:
		appPath += ".ipa"
	default:
		return "", errors.New("app平台异常")
	}

	return appPath, nil
}

// GetExcel 导出App
func (e *App) GetExcel(list []models.App) ([]byte, error) {
	//sheet名称
	sheetName := "App"
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
