package controllers

import (
	"context"
	cursosDTO "cursos-api/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Service interface {
	GetCourseByID(ctx context.Context, id string) (cursosDTO.CourseDto, error)
	Create(ctx context.Context, course cursosDTO.CourseDto) (string, error)
	Update(ctx context.Context, course cursosDTO.CourseDto) error
	Delete(ctx context.Context, id string) error
	GetCourses(ctx context.Context) (cursosDTO.CoursesDto, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) GetCourseByID(ctx *gin.Context) {
	// Validate ID param
	courseID := strings.TrimSpace(ctx.Param("id"))

	// Get hotel by ID using the service
	course, err := controller.service.GetCourseByID(ctx.Request.Context(), courseID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("error getting course: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, course)
}

func (controller Controller) Create(ctx *gin.Context) {
	// Parse course
	var course cursosDTO.CourseDto
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Create course
	id, err := controller.service.Create(ctx.Request.Context(), course)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error creating course: %s", err.Error()),
		})
		return
	}

	// Send ID
	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (controller Controller) Update(ctx *gin.Context) {
	// Validate ID param
	id := strings.TrimSpace(ctx.Param("id"))

	// Parse course
	var course cursosDTO.CourseDto
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Set the ID from the URL to the course object
	course.Course_id = id

	// Update curso
	if err := controller.service.Update(ctx.Request.Context(), course); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error updating curso: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, gin.H{
		"message": id,
	})
}

func (controller Controller) Delete(ctx *gin.Context) {
	// Validate ID param
	id := strings.TrimSpace(ctx.Param("id"))

	// Delete course
	if err := controller.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error deleting hotel: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, gin.H{
		"message": id,
	})
}

func (controller Controller) GetCourses(c *gin.Context) {
	// Llamar al servicio pasando el contexto
	courses, err := controller.service.GetCourses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting courses: %s", err.Error()),
		})
		return
	}

	// Enviar respuesta
	c.JSON(http.StatusOK, courses)
}
