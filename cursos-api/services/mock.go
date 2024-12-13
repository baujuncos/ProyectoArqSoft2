package services

import (
	"context"
	"cursos-api/dao"
	"cursos-api/dto"
	"github.com/stretchr/testify/mock"
)

// MockRepository simula el comportamiento de la interfaz Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetCourseByID(ctx context.Context, id string) (dao.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dao.Course), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, course dao.Course) (string, error) {
	args := m.Called(ctx, course)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, course dao.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetCourses(ctx context.Context) (dao.Courses, error) {
	args := m.Called(ctx)
	return args.Get(0).(dao.Courses), args.Error(1)
}

// MockQueue simula el comportamiento de la interfaz Queue
type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) Publish(courseNew dto.CourseNew) error {
	args := m.Called(courseNew)
	return args.Error(0)
}
