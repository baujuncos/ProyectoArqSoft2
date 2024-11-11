package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	dto "users-api/dto"
)

type Service interface { //definimos conjunto de metodos de users que tendran implementacion en services de users-api
	GetUserByID(id int64) (dto.UserDto, error)
	CreateUser(user dto.UserDto) (int64, error)
	Login(mail string, password string) (dto.LoginDtoResponse, error)
}

type Controller struct {
	service Service //delegamos responsabilidad de lógica de negocio al Service
}

func NewController(service Service) Controller { //creamos una instancia del controlador (Controller)
	return Controller{ // con una dependencia inyectada (Servicio que implementa)
		service: service,
	}
}

//El controlador está construido para manejar diferentes tipos de solicitudes HTTP relacionadas con usuarios:

func (controller Controller) CreateUser(c *gin.Context) {
	var user dto.UserDto
	err := c.BindJSON(&user)
	if err != nil {

		c.JSON(http.StatusBadGateway, gin.H{ //si el json no es valido devuelve 502
			"error": fmt.Sprintf("Error al parsear el JSON: %s", err.Error()),
		})
		return
	}

	//Validaciones explicitas
	if user.Email == "" || user.Nombre == "" || user.Apellido == "" {
		fmt.Printf("Received JSON: %v\n", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email, nombre, and apellido are required and cannot be empty",
		})
		return
	}

	id, err := controller.service.CreateUser(user) //llamamos al servicio que crea el usuario

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ //si hay error devuelve 500
			"error": fmt.Sprintf("error creating user: %s", err.Error()),
		})
		return
	}

	// Devolvemos ID del nuevo usuario
	c.JSON(http.StatusCreated, gin.H{ //exito en registro devuelve 201 Created
		"id": id,
	})
}

func (controller Controller) GetUsersByID(c *gin.Context) {
	userID := c.Param("id")                     //Extraemos el ID de la URL
	id, err := strconv.ParseInt(userID, 10, 64) //convertimos ID a tipo entero
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ //si hay error en la conversion devolvemos 400
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	user, err := controller.service.GetUserByID(id) // llamamos al servicio correspondiente
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{ //en caso de no encontrar el usuario devolvemos 404 not found
			"error": fmt.Sprintf("user not found: %s", err.Error()),
		})
		return
	}

	// Send user
	c.JSON(http.StatusOK, user) //si encuentra al usuario devolvemos 200 y el usuario
}

func (controller Controller) Login(c *gin.Context) {
	var user dto.UserDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ //error devuelve 400
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	response, err := controller.service.Login(user.Email, user.Password) //llamamos al servicio para autenticar el usuario
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{ //en caso de credenciales incorrectas devuelve 401 Unauthorized
			"error": fmt.Sprintf("unauthorized: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response) //en login exitoso devolvemos 200 y el token
}
