package services

import (
	"context"
	cursosDAO "cursos-api/dao"
	cursosDTO "cursos-api/dto"
	"fmt"
)

type Repository interface {
	GetCourseByID(ctx context.Context, id string) (cursosDAO.Course, error)
	Create(ctx context.Context, curso cursosDAO.Course) (string, error)
	Update(ctx context.Context, curso cursosDAO.Course) error
	Delete(ctx context.Context, id string) error
	GetCourses(ctx context.Context) (cursosDAO.Courses, error)
}

type Queue interface {
	Publish(cursoNew cursosDTO.CourseNew) error
}

type Service struct {
	mainRepository Repository
	eventsQueue    Queue
}

func NewService(mainRepository Repository, eventsQueue Queue) Service {
	return Service{
		mainRepository: mainRepository,
		eventsQueue:    eventsQueue,
	}
}

func (service Service) GetCourseByID(ctx context.Context, id string) (cursosDTO.CourseDto, error) {
	courseDAO, err := service.mainRepository.GetCourseByID(ctx, id)
	if err != nil {
		// Get hotel from main repository
		courseDAO, err = service.mainRepository.GetCourseByID(ctx, id)
		if err != nil {
			return cursosDTO.CourseDto{}, fmt.Errorf("error getting hotel from repository: %v", err)
		}

	}
	// Convert DAO to DTO
	return cursosDTO.CourseDto{
		Course_id:    courseDAO.Course_id,
		Nombre:       courseDAO.Nombre,
		Profesor_id:  courseDAO.Profesor_id,
		Categoria:    courseDAO.Categoria,
		Descripcion:  courseDAO.Descripcion,
		Valoracion:   courseDAO.Valoracion,
		Duracion:     courseDAO.Duracion,
		Requisitos:   courseDAO.Requisitos,
		Url_image:    courseDAO.Url_image,
		Fecha_inicio: courseDAO.Fecha_inicio,
	}, nil
}

func (service Service) Create(ctx context.Context, course cursosDTO.CourseDto) (string, error) {
	record := cursosDAO.Course{
		Course_id:    course.Course_id,
		Nombre:       course.Nombre,
		Profesor_id:  course.Profesor_id,
		Categoria:    course.Categoria,
		Descripcion:  course.Descripcion,
		Valoracion:   course.Valoracion,
		Duracion:     course.Duracion,
		Requisitos:   course.Requisitos,
		Url_image:    course.Url_image,
		Fecha_inicio: course.Fecha_inicio,
	}
	id, err := service.mainRepository.Create(ctx, record)
	if err != nil {
		return "", fmt.Errorf("error creating course in main repository: %w", err)
	}
	// Set ID from main repository to use in the rest of the repositories
	if err := service.eventsQueue.Publish(cursosDTO.CourseNew{
		Operation: "CREATE",
		Course_id: id,
	}); err != nil {
		return "", fmt.Errorf("error publishing course new: %w", err)
	}
	return id, nil
}

func (service Service) Update(ctx context.Context, course cursosDTO.CourseDto) error {
	// Convert domain model to DAO model
	record := cursosDAO.Course{
		Course_id:    course.Course_id,
		Nombre:       course.Nombre,
		Profesor_id:  course.Profesor_id,
		Categoria:    course.Categoria,
		Descripcion:  course.Descripcion,
		Valoracion:   course.Valoracion,
		Duracion:     course.Duracion,
		Requisitos:   course.Requisitos,
		Url_image:    course.Url_image,
		Fecha_inicio: course.Fecha_inicio,
	}

	// Update the curso in the main repository
	err := service.mainRepository.Update(ctx, record)
	if err != nil {
		return fmt.Errorf("error updating curso in main repository: %w", err)
	}

	// Publish an event for the update operation
	if err := service.eventsQueue.Publish(cursosDTO.CourseNew{
		Operation: "UPDATE",
		Course_id: course.Course_id,
	}); err != nil {
		return fmt.Errorf("error publishing curso update: %w", err)
	}

	return nil
}

func (service Service) Delete(ctx context.Context, id string) error {
	// Delete the curso from the main repository
	err := service.mainRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting curso from main repository: %w", err)
	}

	// Publish an event for the delete operation
	if err := service.eventsQueue.Publish(cursosDTO.CourseNew{
		Operation: "DELETE",
		Course_id: id,
	}); err != nil {
		return fmt.Errorf("error publishing curso delete: %w", err)
	}

	return nil
}

func (service Service) GetCourses(ctx context.Context) (cursosDTO.CoursesDto, error) {
	// Llamar al metodo GetCourses del repositorio para obtener todos los cursos
	coursesDAO, err := service.mainRepository.GetCourses(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting courses from repository: %w", err)
	}

	// Convertir cada curso de DAO a DTO
	var coursesDto []cursosDTO.CourseDto

	for _, course := range coursesDAO {
		courseDto := cursosDTO.CourseDto{
			Course_id:    course.Course_id,
			Nombre:       course.Nombre,
			Profesor_id:  course.Profesor_id,
			Categoria:    course.Categoria,
			Descripcion:  course.Descripcion,
			Valoracion:   course.Valoracion,
			Duracion:     course.Duracion,
			Requisitos:   course.Requisitos,
			Url_image:    course.Url_image,
			Fecha_inicio: course.Fecha_inicio,
		}
		coursesDto = append(coursesDto, courseDto)
	}

	return coursesDto, nil
}
