package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/common/core/sdk/pkg"
	"go-admin/common/core/tools/language"
	"go-admin/common/middleware/auth/jwtauth"
	"go-admin/config/lang"
	"strconv"
)

//
//  InitAuth
//  @Description: 初始化
//
func InitAuth() {
	jwtauth.Init()
}

//
//  Login
//  @Description: 登录
//  @param c
//
func Login(c *gin.Context) {
	jwtauth.JwtAuthMiddleware.LoginHandler(c)
	return
}

//
//  Logout
//  @Description: 退出登录
//  @param c
//
func Logout(c *gin.Context) {

}

func Get(c *gin.Context, key string) interface{} {
	data := ExtractClaims(c)
	if data[key] != nil {
		return data[key]
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " Get 缺少 " + key)
	return nil
}

func GetUserName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["nice"] != nil {
		return (data["nice"]).(string)
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserName 缺少 nice")
	return ""
}

func GetRoleName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["rolekey"] != nil {
		return (data["rolekey"]).(string)
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleName 缺少 rolekey")
	return ""
}

//
//  GetUserId
//  @Description: 获取用户编号，整型
//  @param c
//  @return int64
//  @return int
//  @return error
//
func GetUserId(c *gin.Context) (int64, int, error) {
	l := "zh-CN"
	languages := language.ParseAcceptLanguage(c.GetHeader("Accept-Language"), nil)
	if len(languages) > 0 {
		l = languages[0]
	}
	data := ExtractClaims(c)
	if data["identity"] != nil {
		id, err := strconv.ParseInt(data["identity"].(string), 10, 64)
		if err != nil {
			return 0, lang.AuthErr, lang.MsgErr(lang.AuthErr, l)
		}
		return id, lang.SuccessCode, nil
	}
	return 0, lang.AuthErr, lang.MsgErr(lang.AuthErr, l)
}

func GetRoleId(c *gin.Context) int64 {
	data := ExtractClaims(c)
	if data["roleid"] != nil {
		i := int64((data["roleid"]).(float64))
		return i
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleId 缺少 roleid")
	return 0
}

//
//  GetUserIdStr
//  @Description: 获取用户编号，字符串
//  @param c
//  @return string
//  @return int
//  @return error
//
func GetUserIdStr(c *gin.Context) (string, int, error) {
	data := ExtractClaims(c)
	if data["identity"] != nil {
		return (data["identity"]).(string), lang.SuccessCode, nil
	}
	l := "zh-CN"
	languages := language.ParseAcceptLanguage(c.GetHeader("Accept-Language"), nil)
	if len(languages) > 0 {
		l = languages[0]
	}
	return "", lang.AuthErr, lang.MsgErr(lang.AuthErr, l)
}

func GetDeptId(c *gin.Context) int {
	data := ExtractClaims(c)
	if data["deptid"] != nil {
		i := int((data["deptid"]).(float64))
		return i
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetDeptId 缺少 deptid")
	return 0
}

func GetDeptName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["deptkey"] != nil {
		return (data["deptkey"]).(string)
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetDeptName 缺少 deptkey")
	return ""
}

// -------------------------- jwt start --------------------------

func ExtractClaims(c *gin.Context) jwtauth.MapClaims {
	claims, exists := c.Get(jwtauth.JwtPayloadKey)
	if !exists {
		return make(jwtauth.MapClaims)
	}

	return claims.(jwtauth.MapClaims)
}

func JwtGetVal(c *gin.Context, key string) interface{} {
	data := ExtractClaims(c)
	if data[key] != nil {
		return data[key]
	}
	fmt.Println(pkg.GetCurrentTimeStr() + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetArticle 缺少 " + key)
	return nil
}

// -------------------------- jwt end --------------------------
