package clients

import "cursos-api/dto"

type Mock struct{}

func NewMock() Mock {
	return Mock{}
}

func (Mock) Publish(cursesNew dto.CourseNew) error {
	return nil
}
