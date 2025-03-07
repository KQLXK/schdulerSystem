package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// ValidateRequest 校验请求参数
func ValidateRequest(c *gin.Context, request interface{}) bool {
	// 绑定请求参数
	if err := c.ShouldBindJSON(request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request parameters")
		return false
	}

	// 使用 validator 校验参数
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return false
	}

	return true
}
