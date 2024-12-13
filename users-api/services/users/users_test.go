package users

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing" // Biblioteca estándar para testing.
	"users-api/dao"
	"users-api/dto"

	"github.com/stretchr/testify/assert" // Biblioteca para aserciones.
	_ "github.com/stretchr/testify/mock"
)

func TestGetUserByID_Success(t *testing.T) {
	// Configuramos el mock del repositorio.
	mockRepo := new(MockRepository)
	mockTokenizer := new(MockTokenizer)
	service := NewService(mockRepo, mockRepo, mockRepo, mockTokenizer)

	// Definimos el usuario esperado.
	expectedUser := dao.Users{User_id: 1, Email: "test@example.com", Nombre: "John"}
	mockRepo.On("GetUserByID", int64(1)).Return(expectedUser, nil)

	// Ejecutamos la prueba.
	result, err := service.GetUserByID(1)

	// Validaciones
	assert.NoError(t, err)                            // No debe haber error.
	assert.Equal(t, "test@example.com", result.Email) // El email debe coincidir.
	mockRepo.AssertExpectations(t)                    // Verifica que se llamó al mock correctamente.
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockTokenizer := new(MockTokenizer)
	service := NewService(mockRepo, mockRepo, mockRepo, mockTokenizer)

	// Simulamos un error "not found".
	mockRepo.On("GetUserByID", int64(1)).Return(dao.Users{}, errors.New("not found"))

	// Ejecutamos la prueba.
	_, err := service.GetUserByID(1)

	// Validaciones
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockTokenizer := new(MockTokenizer)
	service := NewService(mockRepo, mockRepo, mockRepo, mockTokenizer)

	// Definimos el usuario DTO y DAO.
	inputUser := dto.UserDto{
		Email:    "test@example.com",
		Password: "password123",
		Nombre:   "John",
		Apellido: "Doe",
		Admin:    false,
	}

	mockRepo.On("CreateUser", mock.MatchedBy(func(user dao.Users) bool {
		return user.Email == "test@example.com" &&
			user.Password == Hash("password123") &&
			user.Nombre == "John" &&
			user.Apellido == "Doe" &&
			!user.Admin
	})).Return(int64(1), nil)

	// Ejecutamos la prueba.
	id, err := service.CreateUser(inputUser)

	// Validaciones
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockTokenizer := new(MockTokenizer)
	service := NewService(mockRepo, mockRepo, mockRepo, mockTokenizer)

	// Simulamos un usuario y token.
	inputEmail := "test@example.com"
	inputPassword := "password123"
	hashedPassword := Hash(inputPassword)

	user := dao.Users{
		User_id:  1,
		Email:    inputEmail,
		Password: hashedPassword,
		Admin:    true,
	}

	mockRepo.On("GetUserByEmail", inputEmail).Return(user, nil)
	mockTokenizer.On("GenerateToken", inputEmail, int64(1)).Return("mocked-token", nil)

	// Ejecutar
	response, err := service.Login(inputEmail, inputPassword)

	// Validaciones
	assert.NoError(t, err)
	assert.Equal(t, "mocked-token", response.Token)
	assert.True(t, response.Admin)
	mockRepo.AssertExpectations(t)
	mockTokenizer.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockRepository)
	mockTokenizer := new(MockTokenizer)
	service := NewService(mockRepo, mockRepo, mockRepo, mockTokenizer)

	inputEmail := "test@example.com"
	wrongPassword := "wrongpassword"
	hashedPassword := Hash("correctpassword")

	user := dao.Users{
		User_id:  1,
		Email:    inputEmail,
		Password: hashedPassword,
	}

	mockRepo.On("GetUserByEmail", inputEmail).Return(user, nil)

	// Ejecutar
	_, err := service.Login(inputEmail, wrongPassword)

	// Validaciones
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}
