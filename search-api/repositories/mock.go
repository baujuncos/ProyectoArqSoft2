package repositories

import (
	"search-api/dao"
)

type Mock struct {
	data map[int64]dao.CourseDao
}

func NewMock() Mock {
	return Mock{
		data: make(map[int64]dao.CourseDao),
	}
}
