package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	adminService "go-admin/app/admin/service"
	"go-admin/app/plugins/appmanager/models"
	"go-admin/app/plugins/appmanager/service/dto"
	"go-admin/common/actions"
	"go-admin/common/core/config"
	"go-admin/common/core/sdk/service"
	cDto "go-admin/common/dto"
	"go-admin/common/utils/fileUtils/ossUtils"
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

type AppManager struct {
	service.Service
}

func NewAppManagerService(s *service.Service) *AppManager {
	var srv = new(AppManager)
	srv.Orm = s.Orm
	srv.Log = s.Log
	return srv
}

// GetPage 获取AppManager列表
func (e *AppManager) GetPage(c *dto.AppManagerQueryReq, p *actions.DataPermission) ([]models.AppManager, int64, error) {

	var list []models.AppManager
	var data models.AppManager
	var count int64

	err := e.Orm.Order("created_at desc").Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Find(&list).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		e.Log.Errorf("AppManagerService GetPage error:%s", err)
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

// Get 获取AppManager对象
func (e *AppManager) Get(id int64, p *actions.DataPermission) (*models.AppManager, error) {
	if id <= 0 {
		return nil, errors.New("参数错误")
	}
	model := &models.AppManager{}
	err := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(model, id).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无权查看该数据%s", err))
	}
	return model, nil
}

// QueryOne 通过自定义条件获取一条记录
func (e *AppManager) QueryOne(queryCondition *dto.AppManagerQueryReq) (*models.AppManager, error) {
	model := &models.AppManager{}
	err := e.Orm.Model(&models.AppManager{}).
		Scopes(
			cDto.MakeCondition(queryCondition.GetNeedSearch()),
		).First(model).Error
	if err != nil {
		e.Log.Errorf("AppManagerService QueryOne error:%s", err)
		return nil, err
	}
	return model, nil
}

// Count 获取条数
func (e *AppManager) Count(c *dto.AppManagerQueryReq) (int64, error) {
	var err error
	var count int64
	err = e.Orm.Model(&models.AppManager{}).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).Limit(-1).Offset(-1).Count(&count).Error

	if err != nil {
		e.Log.Errorf("AppManagerService Count error:%s", err)
		return 0, err
	}
	return count, nil
}

// Insert 创建AppManager对象
func (e *AppManager) Insert(c *dto.AppManagerInsertReq) error {
	if c.CurrAdminId <= 0 {
		return errors.New("参数错误")
	}
	if c.Platform == "" {
		return errors.New("请选择一个平台")
	}
	if c.Version == "" {
		return errors.New("请输入版本号")
	}
	if c.DownloadType == "" {
		return errors.New("请选择下载类型")
	}
	if c.DownloadType == AppDownloadTypeOss {
		if c.Type == "" {
			return errors.New("请选择App类型")
		}
		if c.LocalAddress == "" {
			return errors.New("请上传App文件")
		}
	}
	if c.Remark == "" {
		return errors.New("更新内容不得为空")
	}
	query := dto.AppManagerQueryReq{}
	query.Platform = c.Platform
	query.Type = c.Type
	query.Version = c.Version
	count, err := e.Count(&query)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("版本已存在,请检查后重新录入")
	}

	ossKey := ""
	buckName := ""

	//oss上传
	if c.DownloadType == AppDownloadTypeOss {
		result, err := e.uploadOssFile(c.Type, c.Version, c.Platform, c.LocalAddress)
		if err != nil {
			return err
		}
		ossKey = result.OssKey
		buckName = result.BucketName
	}

	now := time.Now()
	data := models.AppManager{}
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
	data.CreateBy = c.CurrAdminId
	data.UpdateBy = c.CurrAdminId
	data.CreatedAt = &now
	data.UpdatedAt = &now
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysConfigService Insert error:%s", err)
		return err
	}
	return nil
}

// Update 修改AppManager对象
func (e *AppManager) Update(c *dto.AppManagerUpdateReq, p *actions.DataPermission) (bool, error) {
	if c.Id <= 0 || c.CurrAdminId <= 0 {
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
		updates["update_by"] = c.CurrAdminId
		updates["updated_at"] = time.Now()
		err = e.Orm.Model(&models.AppManager{}).Where("id=?", c.Id).Updates(updates).Error
		if err != nil {
			e.Log.Errorf("SysConfigService Update error:%s", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Remove 删除AppManager
func (e *AppManager) Remove(ids []int64, p *actions.DataPermission) error {
	if len(ids) <= 0 {
		return errors.New("参数错误")
	}
	for _, id := range ids {
		result, err := e.Get(id, p)
		if err != nil {
			return err
		}

		//同一个完全相同的版本，可能因为网路有多条记录，但这些记录都指向一个oss资源，此时只有最后一条记录，才能删除oss资源
		query := dto.AppManagerQueryReq{}
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
	var data models.AppManager
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
func (e *AppManager) GetSingleUploadFileInfo(form *multipart.Form, file *multipart.FileHeader, dst *string) error {
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
//  @return *models.AppManager
//  @return error
//
func (e *AppManager) uploadOssFile(appType, version, platform, localAddress string) (*models.AppManager, error) {
	appManager := models.AppManager{}
	appManager.Type = appType
	appManager.Version = version
	appManager.Platform = platform
	key, err := e.generateAppOssObjectKey(&appManager)
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
	appManager.BucketName = client.BucketName
	appManager.OssKey = key
	return &appManager, nil
}

//
//  generateAppOssUrl
//  @Description: 获取app下载链接
//  @receiver e
//  @param appManager
//  @return string
//  @return error
//
func (e *AppManager) generateAppOssUrl(appManager *models.AppManager) (string, error) {
	appPath, err := e.generateAppOssObjectKey(appManager)
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
//  @param appManager
//  @return *ossUtils.ALiYunOSS
//  @return error
//
func (e *AppManager) getOssClient() (*ossUtils.ALiYunOSS, error) {
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
//  @param appManager
//  @return string
//  @return error
//
func (e *AppManager) generateAppOssObjectKey(appManager *models.AppManager) (string, error) {
	var sysConfService = adminService.NewSysConfigService(&e.Service)

	//app目录
	appPath := sysConfService.GetWithKeyStr("app_oss_root_path", "")
	appPath += appManager.Type
	appPath += "_" + appManager.Version

	switch appManager.Platform {
	case AppPlatformAndroid:
		appPath += ".apk"
	case AppPlatformIOS:
		appPath += ".ipa"
	default:
		return "", errors.New("app平台异常")
	}

	return appPath, nil
}
