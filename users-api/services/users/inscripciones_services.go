package users

import (
	"fmt"
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

func (service ServiceIns) GetInscripcionesByUserID(userID int64) ([]string, error) {
	return service.mainRepositoryIns.GetInscripcionesByUserID(userID)
}

func (service ServiceIns) GetInscripcionesByCursoID(idCurso string) ([]int64, error) {
	return service.mainRepositoryIns.GetInscripcionesByCursoID(idCurso)
}
