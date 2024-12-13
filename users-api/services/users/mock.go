package users

import (
	"github.com/stretchr/testify/mock" // Biblioteca de mocks.
	"users-api/dao"
)

// MockRepository simula el comportamiento de la interfaz Repository.
type MockRepository struct {
	mock.Mock // Hereda de testify/mock para usar sus métodos.
}

// GetUserByID simula la búsqueda de un usuario por ID.
func (m *MockRepository) GetUserByID(id int64) (dao.Users, error) {
	// Llamamos a "Called" con el parámetro `id` y devolvemos el resultado.
	args := m.Called(id)
	return args.Get(0).(dao.Users), args.Error(1)
}

// GetUserByEmail simula la búsqueda de un usuario por Email.
func (m *MockRepository) GetUserByEmail(email string) (dao.Users, error) {
	args := m.Called(email)
	return args.Get(0).(dao.Users), args.Error(1)
}

// CreateUser simula la creación de un usuario en la base de datos.
func (m *MockRepository) CreateUser(user dao.Users) (int64, error) {
	args := m.Called(user)
	return args.Get(0).(int64), args.Error(1)
}

// MockTokenizer simula la interfaz Tokenizer.
type MockTokenizer struct {
	mock.Mock
}

// GenerateToken simula la generación de un token JWT.
func (m *MockTokenizer) GenerateToken(username string, userID int64) (string, error) {
	args := m.Called(username, userID)
	return args.String(0), args.Error(1)
}
