package class

import (
	"errors"
	"math"
	"schedule/dto"
	"schedule/models"
)

var (
	PageNumErr = errors.New("页码超出范围")
)

type ClassQueryByPageFlow struct {
	Page      int
	Pagesize  int
	Classes   []models.Class
	Total     int64
	TotalPage int64
}

func ClassQueryByPage(page int, pagesize int) (*dto.ClassQueryByPageResp, error) {
	return NewClassQueryByPageFlow(page, pagesize).Do()
}

func NewClassQueryByPageFlow(page int, pagesize int) *ClassQueryByPageFlow {
	return &ClassQueryByPageFlow{
		Page:     page,
		Pagesize: pagesize,
	}
}

func (f *ClassQueryByPageFlow) Do() (*dto.ClassQueryByPageResp, error) {
	var resp dto.ClassQueryByPageResp
	resp.Page = f.Page
	resp.PageSize = f.Pagesize

	if err := f.QueryByPage(); err != nil {
		return nil, err
	}

	resp.Classes = f.Convert()

	if err := f.CountTotal(); err != nil {
		return nil, err
	}

	resp.Total = f.Total
	resp.TotalPages = f.TotalPage

	if int64(f.Page) > f.TotalPage {
		return nil, PageNumErr
	}

	return &resp, nil
}

func (f *ClassQueryByPageFlow) QueryByPage() error {
	classes, err := models.NewClassDaoInstance().QueryByPage(f.Page, f.Pagesize)
	if err != nil {
		return err
	}
	f.Classes = classes
	return nil
}

func (f *ClassQueryByPageFlow) CountTotal() error {
	total, err := models.NewClassDaoInstance().CountTotal()
	if err != nil {
		return err
	}
	f.Total = total
	f.TotalPage = int64(math.Ceil(float64(total) / float64(f.Pagesize)))
	return nil
}

func (f *ClassQueryByPageFlow) Convert() []dto.ClassGetResp {
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
