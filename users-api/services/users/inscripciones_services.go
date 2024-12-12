package users

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
	"users-api/dao"
)

type RepositoryIns interface {
	InsertInscripcion(newIns dao.Inscripciones) (int64, error)
	IsSubscribed(idUsuario int64, idCurso string) (bool, error)
	GetInscripcionesByUserID(idUsuario int64) ([]string, error)
	GetInscripcionesByCursoID(idCurso string) ([]int64, error)
}

type ServiceIns struct {
	mainRepositoryIns RepositoryIns
}

func NewServiceIns(mainRepositoryIns RepositoryIns) ServiceIns {

	return ServiceIns{
		mainRepositoryIns: mainRepositoryIns,
	}
}

func (service ServiceIns) IsSubscribed(idUsuario int64, idCurso string) (bool, error) {
	result, err := service.mainRepositoryIns.IsSubscribed(idUsuario, idCurso)
	if err != nil {
		return false, fmt.Errorf("error consulting database")
	}
	return result, nil
}

func (service ServiceIns) InsertInscripcion(idUsuario int64, idCurso string, newFechaInscripcion time.Time) (int64, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Variables compartidas
	var currentEnrollments int
	var capacity int
	var fetchError error

	// Llamadas concurrentes
	wg.Add(2)

	// 1. Obtener cantidad actual de inscripciones
	go func() {
		defer wg.Done()
		count, err := service.mainRepositoryIns.GetInscripcionesByCursoID(idCurso)
		if err != nil {
			fetchError = fmt.Errorf("error fetching enrollments count: %w", err)
			return
		}
		mu.Lock()
		currentEnrollments = len(count)
		mu.Unlock()
	}()

	// 2. Obtener capacidad desde API Cursos
	go func() {
		defer wg.Done()
		cap, err := service.fetchCourseCapacity(idCurso)
		if err != nil {
			fetchError = fmt.Errorf("error fetching course capacity: %w", err)
			return
		}
		mu.Lock()
		capacity = cap
		mu.Unlock()
	}()

	wg.Wait()

	if fetchError != nil {
		return 0, fetchError
	}

	// Verificar disponibilidad
	if currentEnrollments >= capacity {
		return 0, fmt.Errorf("cannot enroll: course is full")
	}

	newInscripcion := dao.Inscripciones{
		IdUsuario:        idUsuario,
		IdCurso:          idCurso,
		FechaInscripcion: newFechaInscripcion,
	}

	// Creamos en main repository
	id, err := service.mainRepositoryIns.InsertInscripcion(newInscripcion)

	if err != nil {
		return 0, fmt.Errorf("error creating user in repository layer: %w", err)
	}

	return id, err
}

func (service ServiceIns) fetchCourseCapacity(courseID string) (int, error) {
	url := fmt.Sprintf("http://cursos-api:8081/courses/%s", courseID)
	fmt.Println("Requesting URL:", url) // Log para verificar la URL final

	// Hacer la solicitud HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error calling courses API: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Verificar el código de respuesta
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("courses API returned status: %d", resp.StatusCode)
	}

	// Estructura para decodificar la respuesta completa
	var result struct {
		Capacidad int `json:"capacidad"` // Extraemos el campo "capacidad"
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	fmt.Println("Course capacity received:", result.Capacidad) // Log del valor obtenido
	return result.Capacidad, nil
}

func (service ServiceIns) GetInscripcionesByUserID(userID int64) ([]string, error) {
	return service.mainRepositoryIns.GetInscripcionesByUserID(userID)
}

func (service ServiceIns) GetInscripcionesByCursoID(idCurso string) ([]int64, error) {
	return service.mainRepositoryIns.GetInscripcionesByCursoID(idCurso)
}
