package class

import (
	"schedule/dto"
	"schedule/models"
)

type ClassSearchFlow struct {
	SearchString string
	Count        int
	Classes      []models.Class
}

func ClassSearch(s string) (*dto.ClassSearchResp, error) {
	return NewClassSearchFlow(s).Do()
}

func NewClassSearchFlow(s string) *ClassSearchFlow {
	return &ClassSearchFlow{
		SearchString: s,
	}
}

func (f *ClassSearchFlow) Do() (*dto.ClassSearchResp, error) {
	var resp dto.ClassSearchResp
	if err := f.Search(); err != nil {
		return nil, err
	}
	resp.Classes = f.Convert()
	resp.TotalCount = int64(f.Count)
	return &resp, nil
}

func (f *ClassSearchFlow) Search() error {
	classes, err := models.NewClassDaoInstance().SearchClass(f.SearchString)
	if err != nil {
		return err
	}
	f.Classes = classes
	f.Count = len(classes)
	return nil
}

func (f *ClassSearchFlow) Convert() []dto.ClassGetResp {
	classResp := make([]dto.ClassGetResp, len(f.Classes))
	for i, class := range f.Classes {
		classResp[i] = dto.ClassGetResp{
			ID:             class.ID,
			Name:           class.Name,
			Academic:       class.Academic,
			Cultivation:    class.Cultivation,
			Type:           class.Type,
			ExpectedYear:   class.ExpectedYear,
			IsGraduation:   class.IsGraduation,
			StudentCount:   class.StudentCount,
			MaxCount:       class.MaxCount,
			Year:           class.Year,
			Department:     class.Department,
			MajorID:        class.MajorID,
			Major:          class.Major,
			Campus:         class.Campus,
			FixedClassroom: class.FixedClassroom,
			IsFixed:        class.IsFixed,
		}
	}
	return classResp
}
