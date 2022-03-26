package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin/common/core/sdk"
)

func WithContextDb(c *gin.Context) {
	c.Set("db", sdk.Runtime.GetDbByKey(c.Request.Host).WithContext(c))
	c.Next()
}
