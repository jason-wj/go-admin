package jwtauth

import (
	"github.com/gin-gonic/gin"
	"go-admin/common/core/config"
	"go-admin/common/middleware/auth/authdto"
	"log"
	"net/http"
	"time"
)

// jwt
var JwtAuthMiddleware = &GinJWTMiddleware{}

// AuthInit
func Init() {
	timeout := time.Hour
	if config.ApplicationConfig.Mode == "dev" {
		timeout = time.Duration(876010) * time.Hour
	} else {
		if config.JwtConfig.Timeout != 0 {
			timeout = time.Duration(config.JwtConfig.Timeout) * time.Second
		}
	}
	var err error
	JwtAuthMiddleware, err = New(&GinJWTMiddleware{
		Realm:           config.ApplicationConfig.Name,
		Key:             []byte(config.JwtConfig.Secret),
		Timeout:         timeout,
		MaxRefresh:      time.Hour,
		PayloadFunc:     PayloadFunc,
		IdentityHandler: IdentityHandler,
		Authenticator:   Authenticator,
		Authorizator:    Authorizator,
		Unauthorized:    Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})
	if err != nil {
		log.Fatalf("JWT Init Error, %s", err.Error())
	}
}

func PayloadFunc(data interface{}) MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		userId, _ := v[authdto.LoginUserId]
		roleName, _ := v[authdto.RoleName]
		roleId, _ := v[authdto.RoleId]
		roleKey, _ := v[authdto.RoleKey]
		userName, _ := v[authdto.UserName]
		dataScope, _ := v[authdto.DataScope]
		return MapClaims{
			IdentityKey:  userId,
			RoleIdKey:    roleId,
			RoleKey:      roleKey,
			NiceKey:      userName,
			DataScopeKey: dataScope,
			RoleNameKey:  roleName,
		}
	}
	return MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := ExtractClaims(c)
	return map[string]interface{}{
		IdentityKey:  claims[IdentityKey],
		"UserName":   claims[NiceKey],
		"RoleKey":    claims[RoleKey],
		"UserId":     claims[IdentityKey],
		"RoleIds":    claims[RoleIdKey],
		DataScopeKey: claims[DataScopeKey],
	}
}

func Authenticator(c *gin.Context) (interface{}, error) {
	userId, b := c.Get(authdto.LoginUserId)
	if !b || userId == "" {
		return nil, ErrFailedAuthentication
	}

	roleId, _ := c.Get(authdto.RoleId)
	roleName, _ := c.Get(authdto.RoleName)
	roleKey, _ := c.Get(authdto.RoleKey)
	userName, _ := c.Get(authdto.UserName)
	dataScope, _ := c.Get(authdto.DataScope)

	resp := map[string]interface{}{
		authdto.LoginUserId: userId,
		authdto.RoleName:    roleName,
		authdto.RoleKey:     roleKey,
		authdto.UserName:    userName,
		authdto.DataScope:   dataScope,
		authdto.RoleId:      roleId,
	}
	return resp, nil
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		userId, _ := v[authdto.LoginUserId]
		if userId != nil {
			c.Set(authdto.LoginUserId, userId.(string))
		}

		roleName, _ := v[authdto.RoleName]
		if roleName != nil {
			c.Set(authdto.RoleName, roleName)
		}

		roleId, _ := v[authdto.RoleIds]
		if roleId != nil {
			c.Set(authdto.RoleIds, roleId)
		}

		userName, _ := v[authdto.UserName]
		if userName != nil {
			c.Set(authdto.UserName, userName)
		}

		dataScope, _ := v[authdto.DataScope]
		if dataScope != nil {
			c.Set(authdto.DataScope, dataScope)
		}

		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	resp := &authdto.Resp{
		Msg:  message,
		Code: code,
	}
	c.JSON(http.StatusOK, resp)
}
