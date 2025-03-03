package utils

import "github.com/gin-gonic/gin"

// SuccessResponse 统一成功响应格式
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 0,         // 状态码，0 表示成功
		"msg":  "success", // 成功消息
		"data": data,      // 返回的数据
	})
}

// ErrorResponse 统一错误响应格式
func ErrorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"code": code, // 状态码
		"msg":  msg,  // 错误消息
	})
}
