// Package app
// @Description: 多语言消息管理，根据业务需要新增消息内容，标识内容勿操作
package lang

import (
	"errors"
	"fmt"
	"go-admin/common/core/logger"
	"go-admin/common/utils/i18n"
	"strings"
)

const (
	//为避免多语言异常，如无特殊必要，该部分勿修改
	//600（包含600）以下必须符合http规则，否则接口会异常
	SuccessCode = 200
	RequestErr  = 400
	AuthErr     = 401
	OpErrCode   = 500

	GetLogErrCode               = 601
	DataDecodeLogErrCode        = 602
	DataGetLogErrCode           = 603
	IdGenErrCode                = 604
	CreateLogErrCode            = 605
	UploadLogErrCode            = 606
	UploadErrCode               = 607
	IdEmptyErrCode              = 608
	DirCreateLogErrCode         = 609
	ParamErrCode                = 610
	DataInsertLogErrCode        = 611
	DataInsertErrCode           = 612
	UpdateLogErrCode            = 613
	IdGenerateErrCode           = 614
	AnnouncementNotAllowErrCode = 615
	SendSuccessCode             = 700
	EmailSendCacheLogErrCode    = 942
	EmailSendLogErrCode         = 943
	EmailSendErrCode            = 944
	EmailGenerateErrCode        = 945
	EmailPortLogErrCode         = 946
	EmailConfigErrCode          = 947
	EmailCodeCmpWrongErrCode    = 948
	EmailCodeCacheLogErrCode    = 949
	EmailFormatErrCode          = 950
	SmsSendCacheLogErrCode      = 986
	SmsSendErrCode              = 987
	SmsSendLogErrCode           = 988
	SmsSendFastErrCode          = 989
	SmsSendFastLogErrCode       = 990
	SmsMobileErrCode            = 991
	SmsMobileLogErrCode         = 992
	SmsInitLogErrCode           = 993
	SmsGenerateErrCode          = 994
	SmsCodeCmpWrongErrCode      = 995
	SmsCodeCacheLogErrCode      = 996
	SmsCodeEmptyCode            = 997
	SmsTemplateEmptyCode        = 998
	SmsConfigErrCode            = 999
	UserMobileTitleEmptyCode    = 1000
	UserMobileEmptyCode         = 1001
	UserPasswordEmptyCode       = 1002
	UserPhoneInputErrCode       = 1003
	UserPasswordErrCode         = 1005
	UserLoginTypeErrCode        = 1006
	UserAccountErrCode          = 1007
	UserMobileCodeEmptyCode     = 1008
	UserPasswordDoubleErrCode   = 1009

	UserPasswordHashLogErrCode    = 1011
	UserCardEmptyCode             = 1012
	UserCardTypeEmptyCode         = 1013
	UserTrueNameEmptyCode         = 1014
	UserHasIdentityErrCode        = 1015
	UserNameSupportCnOrEnErrCode  = 1016
	UserNameLengthErrCode         = 1017
	UserCardSupportNumOrEnErrCode = 1018
	UserCardLengthErrCode         = 1019
	UserCardHasIdentityErrCode    = 1020
	UserCardTypeErrCode           = 1021
	UserUploadFileOneErrCode      = 1022
	UserUploadFileLogErrCode      = 1023
	UserEmailHasBindErrCode       = 1024
	UserEmailHasUsedErrCode       = 1025
	UserEmailCodeEmptyErrCode     = 1026
	UserPayPwdOpErrCode           = 1027
	UserNoIdentityErrCode         = 1028
	ConfigGetLogErrCode           = 2000
	AppPlatformEmptyErrCode       = 3000
	AppTypeEmptyErrCode           = 3001
	AppVersionEmptyErrCode        = 3002
	AppNoCheckErrCode             = 3003
	AppInputVersionErrCode        = 3004
	AppServerVersionErrCode       = 3005
	AppLatestVersionErrCode       = 3006
	AppDownloadUrlErrCode         = 3007
	AppDownloadUrlLogErrCode      = 3008
	AppHasNewVersionErrCode       = 3009

	CountryChina       = 10000
	CountryChinaTaiwan = 10001
	CountryChinaMacao  = 10002
	CountryHongkong    = 10003
	CountrySingapore   = 10004
	CountryCanada      = 10005
	CountryKorea       = 10006
	CountryJapan       = 10007
	CountryThailand    = 10008
	CountryBurma       = 10009
	CountryLaos        = 10010
	CountryAustralia   = 10011
	CountryRussia      = 10012

	//业务扩展-以下可修改

)

var (
	MsgInfo = map[int]string{
		//为避免多语言异常，如无特殊必要，该部分勿修改
		//1-基础通用
		SuccessCode:                   "操作成功",
		RequestErr:                    "请求失败",
		AuthErr:                       "状态失效，请重新登录",
		OpErrCode:                     "操作异常，请检查",
		IdEmptyErrCode:                "编号不得为空",
		GetLogErrCode:                 "获取失败：%s",
		DataDecodeLogErrCode:          "数据解析异常：%s",
		DataGetLogErrCode:             "数据获取异常：%s",
		IdGenErrCode:                  "主键生成异常",
		CreateLogErrCode:              "创建失败：%s",
		UploadLogErrCode:              "上传失败：%s",
		UploadErrCode:                 "上传失败",
		DirCreateLogErrCode:           "文件目录创建异常：%s",
		ParamErrCode:                  "参数错误",
		DataInsertLogErrCode:          "数据新增失败：%s",
		DataInsertErrCode:             "数据新增失败",
		SendSuccessCode:               "发送成功",
		EmailSendCacheLogErrCode:      "邮箱缓存数据保存异常：%s",
		EmailSendErrCode:              "邮箱发送异常",
		EmailSendLogErrCode:           "邮箱发送异常：%s",
		EmailGenerateErrCode:          "邮箱验证码生成异常",
		EmailPortLogErrCode:           "邮箱端口号异常，请检查：%s",
		EmailConfigErrCode:            "邮箱配置异常，请检查",
		EmailCodeCmpWrongErrCode:      "邮箱验证码错误",
		EmailCodeCacheLogErrCode:      "邮箱缓存数据获取异常：%s",
		EmailFormatErrCode:            "邮箱格式不正确",
		SmsSendCacheLogErrCode:        "短信缓存异常：%s",
		SmsSendErrCode:                "短信发送异常",
		SmsSendLogErrCode:             "短信发送异常：%s",
		SmsSendFastErrCode:            "短信发送频率过高：%s",
		SmsSendFastLogErrCode:         "短信发送频率过高：%s",
		SmsMobileErrCode:              "手机号格式异常",
		SmsMobileLogErrCode:           "手机号格式异常：%s",
		SmsInitLogErrCode:             "短信初始化异常：%s",
		SmsGenerateErrCode:            "短信验证码生成异常",
		SmsCodeCmpWrongErrCode:        "手机验证码错误",
		SmsCodeCacheLogErrCode:        "短信缓存数据获取异常：%s",
		SmsCodeEmptyCode:              "短信验证码不得为空",
		SmsTemplateEmptyCode:          "短信模板为空",
		SmsConfigErrCode:              "短信配置异常，请检查",
		UserMobileTitleEmptyCode:      "国际区号不得为空",
		UserMobileEmptyCode:           "手机号不得为空",
		UserMobileCodeEmptyCode:       "手机验证码不得为空",
		UserPasswordEmptyCode:         "密码不得为空",
		UserPhoneInputErrCode:         "手机号输出错误",
		UserPasswordErrCode:           "密码错误，请检查",
		UserLoginTypeErrCode:          "登录类型异常",
		UserAccountErrCode:            "用户账号异常",
		UserPasswordDoubleErrCode:     "新密码两次输入不一致",
		UpdateLogErrCode:              "更新异常：%s",
		IdGenerateErrCode:             "编号生成异常",
		AnnouncementNotAllowErrCode:   "公告内容为空",
		UserPasswordHashLogErrCode:    "密码处理异常：%s",
		UserCardEmptyCode:             "证件号码不得为空",
		UserCardTypeEmptyCode:         "证件类型不得为空",
		UserTrueNameEmptyCode:         "真实姓名不得为空",
		UserHasIdentityErrCode:        "当前用户已实名认证",
		UserNameSupportCnOrEnErrCode:  "姓名仅可输入中文或英文字符",
		UserNameLengthErrCode:         "姓名长度在%s之间",
		UserCardSupportNumOrEnErrCode: "证件号仅可输入数字或英文字符",
		UserCardLengthErrCode:         "证件号码长度在%s之间",
		UserCardHasIdentityErrCode:    "当前证件已被使用，请更换",
		UserCardTypeErrCode:           "证件类型错误",
		UserUploadFileOneErrCode:      "单次仅可上传一个",
		UserUploadFileLogErrCode:      "上传文件异常：%s",
		UserEmailHasBindErrCode:       "邮箱已绑定，请勿重复操作",
		UserEmailHasUsedErrCode:       "邮箱已被占用，请使用新的邮箱",
		UserEmailCodeEmptyErrCode:     "邮箱验证码不得为空",
		UserPayPwdOpErrCode:           "支付密码不安全，请先重置",
		UserNoIdentityErrCode:         "请先实名认证",
		ConfigGetLogErrCode:           "配置信息获取异常：%s",
		AppPlatformEmptyErrCode:       "App平台不得为空",
		AppTypeEmptyErrCode:           "App类型不得为空",
		AppVersionEmptyErrCode:        "App版本号不得为空",
		AppNoCheckErrCode:             "未检测到最新版本",
		AppInputVersionErrCode:        "版本号输入异常",
		AppServerVersionErrCode:       "后台发布App版本异常",
		AppLatestVersionErrCode:       "已是最新版本",
		AppDownloadUrlErrCode:         "获取下载链接异常",
		AppDownloadUrlLogErrCode:      "获取下载链接异常：%s",
		AppHasNewVersionErrCode:       "检测到新版本",

		//2-国家
		CountryChina:       "中国大陆",
		CountryChinaTaiwan: "中国台湾",
		CountryChinaMacao:  "中国澳门",
		CountryHongkong:    "中国香港",
		CountrySingapore:   "新加坡",
		CountryCanada:      "加拿大",
		CountryKorea:       "韩国",
		CountryJapan:       "日本",
		CountryThailand:    "泰国",
		CountryBurma:       "缅甸",
		CountryLaos:        "老挝",
		CountryAustralia:   "澳大利亚",
		CountryRussia:      "俄罗斯",

		//业务扩展-以下可修改
	}
)

//
//  MsgByCode
//  @Description: i18n
//  @param errCode
//  @param lang
//  @return string
//
func MsgByCode(errCode int, lang string) string {
	switch lang {
	case "en":
		return i18n.EnI18nClient.T(MsgInfo[errCode])
	default:
		return MsgInfo[errCode]
	}
}

//
//  MsgByValue
//  @Description: 直接根据值返回对应语言
//  @param value
//  @param lang
//  @return string
//
func MsgByValue(value string, lang string) string {
	switch lang {
	case "en":
		return i18n.EnI18nClient.T(value)
	default:
		return value
	}
}

//
//  MsgErr
//  @Description: 获取error
//  @param errCode
//  @param lang
//  @return error
//
func MsgErr(errCode int, lang string) error {
	return errors.New(MsgByCode(errCode, lang))
}

//
//  MsgErrf
//  @Description:
//  @param errCode
//  @param lang
//  @param f
//  @return error
//
func MsgErrf(errCode int, lang string, f ...interface{}) error {
	return errors.New(fmt.Sprintf(MsgByCode(errCode, lang), f))
}

//
//  MsgLogErrf
//  @Description: 带有参数，有些底层消息不应当被使用者感知，该类消息记录在日志中，并返回应用层可理解的消息
//  @param log  用于记录日志，
//  @param errCodeReplace 最终需要给应用层返回的消息，这里传入消息码。若errCodeReplace=errCode或errCodeReplace<=0，则返回真实消息
//  @param errCode  真实消息码
//  @param lang  语言
//  @param f
//  @return error
//
func MsgLogErrf(log *logger.Helper, lang string, errCodeReplace, errCode int, f ...interface{}) error {
	err := MsgErrf(errCode, lang, f)
	log.Error(err)
	if errCodeReplace <= 0 || errCodeReplace == errCode {
		return err
	}
	return MsgErr(errCodeReplace, lang)
}

//
//  MsgLogErr
//  @Description: 无参数，有些底层消息不应当被使用者感知，该类消息记录在日志中，并返回应用层可理解的消息
//  @param log
//  @param lang
//  @param errCodeReplace
//  @param errCode
//  @return error
//
func MsgLogErr(log *logger.Helper, lang string, errCodeReplace, errCode int) error {
	err := MsgErr(errCode, lang)
	log.Error(err)
	if errCodeReplace <= 0 || errCodeReplace == errCode {
		return err
	}
	return MsgErr(errCodeReplace, lang)
}

//
//  TranslationText
//  @Description: 仅支持 - 分隔符
//  @param l
//  @param name
//  @return string
//
func TranslationText(l string, text string) string {
	values := strings.Split(text, "-")
	if len(values) <= 0 {
		return text
	}
	newValue := MsgByValue(values[0], l)
	return strings.Replace(text, values[0], newValue, 1)
}
