package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	//400
	CourseExsitsStatus      = NewStatus(http.StatusBadRequest, 40001, "课程名或课程号已存在")
	CourseDataInvalidStatus = NewStatus(http.StatusBadRequest, 40002, "学时数据不合法")
	CourseIDEmptyStatus     = NewStatus(http.StatusBadRequest, 40003, "课程ID不能为空")

	//404
	CourseNotFoundStatus = NewStatus(http.StatusNotFound, 40401, "课程未找到")
)

type status struct {
	HTTPcode   int
	Statuscode int
	Message    string
}

func (s status) httpcode() int {
	return s.HTTPcode
}

func (s status) statuscode() int {
	return s.Statuscode
}

func (s status) message() string {
	return s.Message
}

func NewStatus(httpcode int, statuscode int, message string) status {
	return status{
		httpcode,
		statuscode,
		message,
	}
}

func Sucess(c *gin.Context, data interface{}) {
	H := gin.H{
		"status":  200,
		"message": "sucess",
	}
	H["data"] = data
	c.JSON(http.StatusOK, H)
}

func Error(c *gin.Context, s status) {
	c.JSON(s.httpcode(), gin.H{
		"status":  s.statuscode(),
		"message": s.Message,
	})
}

func Errors(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  400,
		"message": err.Error(),
	})
}
