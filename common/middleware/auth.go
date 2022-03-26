package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin/common/middleware/auth/jwtauth"
)

func Auth() gin.HandlerFunc {
	return jwtauth.JwtAuthMiddleware.MiddlewareFunc()
}
