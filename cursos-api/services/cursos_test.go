package services

import (
	"context"
	_ "errors"
	"testing"

	"cursos-api/dao"
	"cursos-api/dto"

	"github.com/stretchr/testify/assert"
)

func TestService_GetCourseByID(t *testing.T) {
	mockRepo := new(MockRepository)
	mockQueue := new(MockQueue)
	service := NewService(mockRepo, mockQueue)
	ctx := context.TODO()

	expectedCourse := dao.Course{
		Course_id: "1", Nombre: "Go Programming", Capacidad: 30,
	}

	mockRepo.On("GetCourseByID", ctx, "1").Return(expectedCourse, nil)

	// Ejecutar función
	result, err := service.GetCourseByID(ctx, "1")

	// Validaciones
	assert.NoError(t, err)
	assert.Equal(t, "Go Programming", result.Nombre)
	mockRepo.AssertExpectations(t)
}

func TestService_Create(t *testing.T) {
	mockRepo := new(MockRepository)
	mockQueue := new(MockQueue)
	service := NewService(mockRepo, mockQueue)
	ctx := context.TODO()

	inputCourse := dto.CourseDto{Nombre: "Go Programming"}
	createdCourse := dao.Course{Course_id: "1", Nombre: "Go Programming"}

	mockRepo.On("Create", ctx, createdCourse).Return("1", nil)
	mockQueue.On("Publish", dto.CourseNew{Operation: "CREATE", Course_id: "1"}).Return(nil)

	// Ejecutar función
	id, err := service.Create(ctx, inputCourse)

	// Validaciones
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
	mockRepo.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

func TestService_Update(t *testing.T) {
	mockRepo := new(MockRepository)
	mockQueue := new(MockQueue)
	service := NewService(mockRepo, mockQueue)
	ctx := context.TODO()

	inputCourse := dto.CourseDto{Course_id: "1", Nombre: "Updated Go"}
	updatedCourse := dao.Course{Course_id: "1", Nombre: "Updated Go"}

	mockRepo.On("Update", ctx, updatedCourse).Return(nil)
	mockQueue.On("Publish", dto.CourseNew{Operation: "UPDATE", Course_id: "1"}).Return(nil)

	// Ejecutar función
	err := service.Update(ctx, inputCourse)

	// Validaciones
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

func TestService_Delete(t *testing.T) {
	mockRepo := new(MockRepository)
	mockQueue := new(MockQueue)
	service := NewService(mockRepo, mockQueue)
	ctx := context.TODO()

	mockRepo.On("Delete", ctx, "1").Return(nil)
	mockQueue.On("Publish", dto.CourseNew{Operation: "DELETE", Course_id: "1"}).Return(nil)

	// Ejecutar función
	err := service.Delete(ctx, "1")

	// Validaciones
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

func TestService_GetCourses(t *testing.T) {
	mockRepo := new(MockRepository)
	mockQueue := new(MockQueue)
	service := NewService(mockRepo, mockQueue)
	ctx := context.TODO()

	expectedCourses := dao.Courses{
		{Course_id: "1", Nombre: "Go Programming"},
		{Course_id: "2", Nombre: "Python Basics"},
	}

	mockRepo.On("GetCourses", ctx).Return(expectedCourses, nil)

	// Ejecutar función
	result, err := service.GetCourses(ctx)

	// Validaciones
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Go Programming", result[0].Nombre)
	mockRepo.AssertExpectations(t)
}
