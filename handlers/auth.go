package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"schedule/middleware"
	"schedule/utils"
	"time"
)

// LoginRequest 定义了登录请求的结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 模拟用户验证（实际项目中应从数据库验证）
	if req.Username != "admin" || req.Password != "admin123" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// 生成 JWT Token
	token, err := generateToken("1", "admin") // 假设用户ID为1，角色为admin
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 返回 Token
	utils.SuccessResponse(c, gin.H{
		"token": token,
	})
}

// generateToken 生成 JWT Token
func generateToken(userID string, role string) (string, error) {
	claims := &middleware.Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token 有效期 24 小时
			Issuer:    "schedule-system",                     // 签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JWTKey)
}
