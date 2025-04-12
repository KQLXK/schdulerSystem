package middleware

import (
	"github.com/gin-gonic/gin"
	"schedule/commen/result"
	"schedule/commen/utils"
	"strings"
)

// AuthMiddleware 是 JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			result.Error(c, result.TokenRequiredStatus)
			c.Abort()
			return
		}

		// 检查 Authorization 头的格式是否为 "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			result.Error(c, result.TokenFormatErrStatus)
			c.Abort()
			return
		}
		token := tokenParts[1]

		jwtinfo, err := utils.GetInfoFromToken(token)
		if err != nil {
			if err == utils.TokenExpiredErr {
				result.Error(c, result.TokenExpiredStatus)
			}
			result.Errors(c, err)
		}

		c.Set("identity", jwtinfo.Identity)
		c.Set("user_id", jwtinfo.UserID)

		// 继续处理请求
		c.Next()
	}
}
