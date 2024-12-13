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
	// Crear una instancia de WaitGroup para sincronizar las Go Routines
	var wg sync.WaitGroup
	// Mutex para proteger el acceso a las variables compartidas
	var mu sync.Mutex

	// Variables compartidas entre las Go Routines
	var currentEnrollments int // Numero actual de inscripciones
	var capacity int           // Capacidad maxima del curso
	var fetchError error       // Variable para almacenar cualquier error durante las operaciones concurrentes

	// Añadimos 2 tareas al WaitGroup (una por cada Go Routine)
	wg.Add(2)

	// -------------------- Tarea 1: Obtener inscripciones actuales --------------------
	go func() {
		defer wg.Done() // Marca la tarea como completada al finalizar la rutina

		// Llamada al repositorio para obtener la cantidad de inscripciones del curso
		count, err := service.mainRepositoryIns.GetInscripcionesByCursoID(idCurso)
		if err != nil {
			// Si ocurre un error, lo guardamos en fetchError de forma segura usando Mutex
			fetchError = fmt.Errorf("error fetching enrollments count: %w", err)
			return
		}

		// Guardar la cantidad de inscripciones actuales de forma segura usando Mutex
		mu.Lock()
		currentEnrollments = len(count) // La cantidad de inscripciones actuales
		mu.Unlock()
	}()

	// -------------------- Tarea 2: Obtener capacidad del curso desde API Cursos --------------------
	go func() {
		defer wg.Done() // Marca la tarea como completada al finalizar la rutina

		// Llamada a la función que obtiene la capacidad del curso desde la API de Cursos
		cap, err := service.fetchCourseCapacity(idCurso)
		if err != nil {
			// Si ocurre un error, lo guardamos en fetchError de forma segura usando Mutex
			mu.Lock()
			fetchError = fmt.Errorf("error fetching course capacity: %w", err)
			mu.Unlock()
			return
		}

		// Guardar la capacidad obtenida de forma segura usando Mutex
		mu.Lock()
		capacity = cap // La capacidad máxima del curso
		mu.Unlock()
	}()

	// -------------------- Esperar a que todas las Go Routines finalicen --------------------
	wg.Wait() // Bloquea la ejecución hasta que ambas Go Routines hayan terminado

	// -------------------- Verificar si hubo algún error en las Go Routines --------------------
	if fetchError != nil {
		// Si se encontró un error en alguna rutina, se devuelve el error
		return 0, fetchError
	}

	// -------------------- Verificar la disponibilidad del curso --------------------
	if currentEnrollments >= capacity {
		// Si las inscripciones actuales son iguales o mayores a la capacidad, devolver ID = 0
		return 0, fmt.Errorf("cannot enroll: course is full") // El curso está lleno
	}

	// -------------------- Insertar la inscripción en la base de datos --------------------
	newInscripcion := dao.Inscripciones{
		IdUsuario:        idUsuario,
		IdCurso:          idCurso,
		FechaInscripcion: newFechaInscripcion,
	}

	// Llamada al repositorio para insertar la nueva inscripción
	id, err := service.mainRepositoryIns.InsertInscripcion(newInscripcion)

	if err != nil {
		return 0, fmt.Errorf("error creating user in repository layer: %w", err)
	}

	// Devolver el ID de la inscripción creada
	return id, err
}

// fetchCourseCapacity obtiene la capacidad del curso llamando a la API de Cursos.
// Recibe el ID del curso como parámetro y devuelve la capacidad del curso (int) o un error.
func (service ServiceIns) fetchCourseCapacity(courseID string) (int, error) {
	// Construir la URL para hacer la llamada a la API de Cursos.
	// La URL incluye el ID del curso como parámetro dinámico.
	url := fmt.Sprintf("http://cursos-api:8081/courses/%s", courseID)
	fmt.Println("Requesting URL:", url) // Log para verificar la URL final

	// Realizar una solicitud HTTP GET a la URL construida.
	resp, err := http.Get(url)
	if err != nil {
		// Si ocurre un error al hacer la solicitud (problemas de red, servidor caído, etc.), devolver el error.
		return 0, fmt.Errorf("error calling courses API: %w", err)
	}
	defer func(Body io.ReadCloser) { // Cerrar el cuerpo de la respuesta para liberar recursos.
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Verificar el código de estado HTTP en la respuesta.
	if resp.StatusCode != http.StatusOK {
		// Si el servidor devuelve un estado diferente a 200 OK, devolver un error con el código de estado.
		return 0, fmt.Errorf("courses API returned status: %d", resp.StatusCode)
	}

	// Definir una estructura temporal para decodificar la respuesta JSON.
	// Solo nos interesa extraer el campo "capacidad" del JSON de la respuesta.
	var result struct {
		Capacidad int `json:"capacidad"` // Extraemos el campo "capacidad"
	}

	// Decodificar el cuerpo de la respuesta HTTP (JSON) en la estructura temporal.
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Si ocurre un error al decodificar el JSON, devolver el error.
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	// Log para verificar que la capacidad fue recibida correctamente.
	fmt.Println("Course capacity received:", result.Capacidad) // Log del valor obtenido

	// Devolver la capacidad obtenida del curso.
	return result.Capacidad, nil
}

func (service ServiceIns) GetInscripcionesByUserID(userID int64) ([]string, error) {
	return service.mainRepositoryIns.GetInscripcionesByUserID(userID)
}

func (service ServiceIns) GetInscripcionesByCursoID(idCurso string) ([]int64, error) {
	return service.mainRepositoryIns.GetInscripcionesByCursoID(idCurso)
}
