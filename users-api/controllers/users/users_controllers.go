<<<<<<< Updated upstream
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

func (controller Controller) GetUsersByID(c *gin.Context) {
	userID := c.Param("id")                     //Extraemos el ID de la URL
	id, err := strconv.ParseInt(userID, 10, 64) //convertimos ID a tipo entero
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ //si hay error en la conversion devolvemos 400
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
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

func (controller Controller) Create(c *gin.Context) {
	var user dto.UserDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{ //si el json no es valido devuelve 502
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
	id, err := controller.service.Create(user) //llamamos al servicio que crea el usuario
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ //si hay error devuelve 500
			"error": fmt.Sprintf("error creating user: %s", err.Error()),
		})
		return
	}

	// Send ID
	c.JSON(http.StatusCreated, gin.H{ //exito en registro devuelve 201 Created
		"id": id,
	})
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
=======
package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	domain "users-api/users_dto"
)

type Service interface {
	GetAll() ([]domain.User, error)
	GetUSerByID(id int64) (domain.User, error)
	CreateUser(user domain.User) (int64, error)
	Login(email string, password string) (domain.LoginResponse, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) GetAll(c *gin.Context) {
	// Invoke service
	users, err := controller.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting all users: %s", err.Error()),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, users)
}

func (controller Controller) GetUserByID(c *gin.Context) {
	// Parse user ID from HTTP request
	userID := c.Param("id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
	user, err := controller.service.GetUSerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("user not found: %s", err.Error()),
		})
		return
	}

	// Send user
	c.JSON(http.StatusOK, user)
}

func (controller Controller) CreateUser(c *gin.Context) {
	// Parse user from HTTP Request
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
	id, err := controller.service.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error creating user: %s", err.Error()),
		})
		return
	}

	// Send ID
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (controller Controller) Update(c *gin.Context) {
	// Parse user ID from HTTP request
	userID := c.Param("id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Parse updated user data from HTTP request
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Set the ID of the user to be updated
	user.ID = id

	// Invoke service
	if err := controller.service.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error updating user: %s", err.Error()),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, user)
}

func (controller Controller) Delete(c *gin.Context) {
	// Parse user ID from HTTP request
	userID := c.Param("id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
	if err := controller.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error deleting user: %s", err.Error()),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (controller Controller) Login(c *gin.Context) {
	// Parse user from HTTP request
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
	response, err := controller.service.Login(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("unauthorized: %s", err.Error()),
		})
		return
	}

	// Send login with token
	c.JSON(http.StatusOK, response)
}
>>>>>>> Stashed changes
