package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	//400
	CourseExsitsStatus         = NewStatus(http.StatusBadRequest, 40001, "课程名或课程号已存在")
	CourseDataInvalidStatus    = NewStatus(http.StatusBadRequest, 40002, "学时数据不合法")
	CourseIDEmptyStatus        = NewStatus(http.StatusBadRequest, 40003, "课程ID不能为空")
	CourseDataEmptyStatus      = NewStatus(http.StatusBadRequest, 40013, "课程ID或课程名为空")
	TeacherExistsStatus        = NewStatus(http.StatusBadRequest, 40004, "教师工号或姓名已存在")
	TeacherDataInvalidStatus   = NewStatus(http.StatusBadRequest, 40005, "教师数据不合法")
	TeacherIDEmptyStatus       = NewStatus(http.StatusBadRequest, 40006, "教师ID不能为空")
	ClassroomExistsStatus      = NewStatus(http.StatusBadRequest, 40007, "教室编号或名称已存在")
	ClassroomDataInvalidStatus = NewStatus(http.StatusBadRequest, 40008, "教室数据不合法")
	ClassroomIDEmptyStatus     = NewStatus(http.StatusBadRequest, 40009, "教室ID不能为空")
	FileNotReceiveStatus       = NewStatus(http.StatusBadRequest, 40010, "接收课程文件失败")
	FileFormatErrStatus        = NewStatus(http.StatusBadRequest, 40011, "文件格式错误")
	PageDataErrStatus          = NewStatus(http.StatusBadRequest, 40012, "页码超出范围")

	//401
	TokenExpiredStatus   = NewStatus(http.StatusUnauthorized, 40101, "登录已过期")
	TokenRequiredStatus  = NewStatus(http.StatusUnauthorized, 40102, "请先登录")
	TokenFormatErrStatus = NewStatus(http.StatusUnauthorized, 40103, "token格式有误")

	//404
	CourseNotFoundStatus    = NewStatus(http.StatusNotFound, 40401, "课程未找到")
	TeacherNotFoundStatus   = NewStatus(http.StatusNotFound, 40402, "教师未找到")
	ClassroomNotFoundStatus = NewStatus(http.StatusNotFound, 40403, "教室未找到")

	ServerInteralErrStatus = NewStatus(http.StatusInternalServerError, 50000, "服务器内部错误")
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
