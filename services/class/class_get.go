package class

import (
	"schedule/models"
)

func GetAllClasses() ([]models.Class, error) {
	return models.NewClassDaoInstance().GetAllClasses()
}

func GetClassByID(id string) (*models.Class, error) {
	class, err := models.NewClassDaoInstance().GetClassByID(id)
	if err != nil {
		return nil, NotFoundError
	}
	return class, nil
}
