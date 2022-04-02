package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/common/core/sdk/api"
	"go-admin/common/core/sdk/pkg/captcha"
)

type System struct {
	api.Api
}

// GenerateCaptchaHandler 获取验证码
func (e System) GenerateCaptchaHandler(c *gin.Context) {
	err := e.MakeContext(c).Errors
	if err != nil {
		e.Error(500, "服务初始化失败！")
		return
	}
	id, b64s, err := captcha.DriverDigitFunc()
	if err != nil {
		e.Logger.Errorf("DriverDigitFunc error, %s", err.Error())
		e.Error(500, "验证码获取失败")
		return
	}
	e.Custom(gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}
