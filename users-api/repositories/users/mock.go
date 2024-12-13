package users

import (
	"github.com/stretchr/testify/mock"
	dao "users-api/dao"
)

// MockMySQL estructura para simular el repositorio MySQL
type MockMySQL struct {
	mock.Mock
}

func NewMockMySQL() *MockMySQL {
	return &MockMySQL{}
}

func (m *MockMySQL) GetUserByID(id int64) (dao.Users, error) {
	args := m.Called(id)
	if err := args.Error(1); err != nil {
		return dao.Users{}, err
	}
	return args.Get(0).(dao.Users), nil
}

func (m *MockMySQL) GetUserByEmail(email string) (dao.Users, error) {
	args := m.Called(email)
	if err := args.Error(1); err != nil {
		return dao.Users{}, err
	}
	return args.Get(0).(dao.Users), nil
}

func (m *MockMySQL) CreateUser(user dao.Users) (int64, error) {
	args := m.Called(user)
	if err := args.Error(1); err != nil {
		return 0, err
	}
	return args.Get(0).(int64), nil
}

func (m *MockMySQL) InsertInscripcion(inscripcion dao.Inscripciones) (int64, error) {
	args := m.Called(inscripcion)
	if err := args.Error(1); err != nil {
		return 0, err
	}
	return args.Get(0).(int64), nil
}

func (m *MockMySQL) GetInscripcionesByCursoID(idCurso string) ([]int64, error) {
	args := m.Called(idCurso)
	if err := args.Error(1); err != nil {
		return nil, err
	}
	return args.Get(0).([]int64), nil
}
