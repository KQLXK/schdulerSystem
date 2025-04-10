package class

import (
	"errors"
	"schedule/models"
)

var (
	NotFoundError = errors.New("班级未找到")
)

func DeleteClass(id string) error {
	// 检查班级是否存在
	if _, err := models.NewClassDaoInstance().GetClassByID(id); err != nil {
		return NotFoundError
	}

	if err := models.NewClassDaoInstance().DeleteClass(id); err != nil {
		return err
	}
	return nil
}
