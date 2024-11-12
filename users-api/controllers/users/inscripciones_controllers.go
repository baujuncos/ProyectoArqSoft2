package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
	"users-api/dto"
)

type ServiceIns interface { //definimos conjunto de metodos de users que tendran implementacion en services de users-api
	InsertInscripcion(idUsuario int64, idCurso string, newFechaInscripcion time.Time) (int64, error)
	IsSubscribed(idUsuario int64, idCurso string) (bool, error)
	GetInscripcionesByUserID(userId int64) ([]string, error)
	GetInscripcionesByCursoID(idCurso string) ([]int64, error)
}

type ControllerIns struct {
	serviceins ServiceIns //delegamos responsabilidad de lógica de negocio al Service
}

func NewControllerIns(service ServiceIns) ControllerIns { //creamos una instancia del controlador (Controller)
	return ControllerIns{ // con una dependencia inyectada (Servicio que implementa)
		serviceins: service,
	}
}

func (controller ControllerIns) InsertInscripcion(c *gin.Context) {
	var inscripcion dto.InscripcionesDto

	if err := c.ShouldBindJSON(&inscripcion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing the JSON": err.Error()})
		return
	}

	subscribed, err := controller.serviceins.IsSubscribed(inscripcion.IdUser, inscripcion.IdCourse) //verificamos si ya esta suscripto el usuario
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error checking subscription": err.Error()})
		return
	} //si cuando chequeamos que esta suscripto nos da error, devuelve error checking subscription

	if subscribed {
		cursoInt, _ := strconv.ParseInt(inscripcion.IdCourse, 10, 64)
		c.JSON(http.StatusConflict, fmt.Sprintf("Usuario %d ya está suscrito al curso %d", inscripcion.IdUser, cursoInt))
		return
	} //si ya esta inscripto devuelve mensaje de que ya esta inscripto

	FechaInscripcion := time.Now()

	id, err := controller.serviceins.InsertInscripcion(inscripcion.IdUser, inscripcion.IdCourse, FechaInscripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error making the subscription at service": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})

}

func (controller ControllerIns) GetInscripcionesByUserID(c *gin.Context) {
	// Obtener el ID del usuario de la URL
	userIdStr := c.Param("user_id")
	fmt.Printf("Valor recibido para id_usuario: %s\n", userIdStr) // Depuración

	//convertimos id de string a int
	userId, err := strconv.ParseInt(userIdStr, 10, 64) //convertimos ID a tipo entero
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario no válido"})
		return
	}

	// Consultar la base de datos para obtener los cursos del usuario
	inscripciones, err := controller.serviceins.GetInscripcionesByUserID(userId)

	if err != nil {
		log.Printf("Error al obtener los cursos del usuario: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los cursos del usuario"})
		return
	}

	// Devolver los cursos encontrados como respuesta
	c.JSON(http.StatusOK, inscripciones) //devolvemos slice de strings de los id cursos a los que esta inscripto el usuario, si es null no tiene inscripciones
}

func (controller ControllerIns) GetInscripcionesByCursoID(c *gin.Context) {
	idCursoParam := c.Param("course_id")

	usuarios, err := controller.serviceins.GetInscripcionesByCursoID(idCursoParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usuarios) //devolvemos slice de ints con los id usuarios inscriptos al curso, si el slice es null no tiene inscripciones
}
