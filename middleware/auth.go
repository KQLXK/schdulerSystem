package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"schedule/commen/utils"
	"strings"
)

// JWTKey 是用于签名和验证 JWT 的密钥
var JWTKey = []byte("your-secret-key") // 请替换为实际的密钥

// Claims 定义了 JWT 的声明结构
type Claims struct {
	UserID string `json:"user_id"` // 用户ID
	Role   string `json:"role"`    // 用户角色
	jwt.StandardClaims
}

// AuthMiddleware 是 JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// 检查 Authorization 头的格式是否为 "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		// 解析 JWT Token
		tokenString := tokenParts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		})

		// 检查 Token 是否有效
		if err != nil || !token.Valid {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		// 继续处理请求
		c.Next()
	}
}
