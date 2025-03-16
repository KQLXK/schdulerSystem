package classroom

import (
	"schedule/models"
)

// DeleteClassroom 删除教室
func DeleteClassroom(id string) error {
	// 检查教室是否存在
	if _, err := models.NewClassroomDao().GetClassroomByID(id); err != nil {
		return NotFoundError
	}

	if err := models.NewClassroomDao().DeleteClassroom(id); err != nil {
		return err
	}
	return nil
}
