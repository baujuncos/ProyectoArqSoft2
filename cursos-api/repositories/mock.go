package repositories

import (
	"context"
	"cursos-api/dao"
	"fmt"
	"github.com/google/uuid"
)

type Mock struct {
	docs map[string]dao.Course
}

func NewMock() Mock {
	return Mock{
		docs: make(map[string]dao.Course),
	}
}

func (repository Mock) GetCourseByID(ctx context.Context, Course_id string) (dao.Course, error) {
	return repository.docs[Course_id], nil
}

func (repository Mock) Create(ctx context.Context, curso dao.Course) (string, error) {
	Course_id := uuid.New().String()
	curso.Course_id = uuid.New().String()
	repository.docs[Course_id] = curso
	return Course_id, nil
}

func (repository Mock) Update(ctx context.Context, curso dao.Course) error {
	currentCurso, exists := repository.docs[curso.Course_id]
	if !exists {
		return fmt.Errorf("course with Course_id %s not found", curso.Course_id)
	}

	// Update only the fields that are non-zero or non-empty
	if curso.Nombre != "" {
		currentCurso.Nombre = curso.Nombre
	}
	if curso.Categoria != "" {
		currentCurso.Categoria = curso.Categoria
	}
	if curso.Descripcion != "" {
		currentCurso.Descripcion = curso.Descripcion
	}
	if curso.Requisitos != "" {
		currentCurso.Requisitos = curso.Requisitos
	}
	if curso.Duracion != 0 {
		currentCurso.Duracion = curso.Duracion
	}
	if curso.Capacidad != 0 {
		currentCurso.Capacidad = curso.Capacidad
	}

	// Save the updated curso back to the mock storage
	repository.docs[curso.Course_id] = currentCurso
	return nil
}
