package teacher

import (
	"schedule/models"
)

// DeleteTeacher 删除教师
func DeleteTeacher(id string) error {
	// 检查教师是否存在
	if _, err := models.GetTeacherByID(id); err != nil {
		return NotFoundError
	}

	if err := models.DeleteTeacher(id); err != nil {
		return err
	}
	return nil
}
