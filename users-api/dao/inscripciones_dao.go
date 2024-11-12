package dao

import "time"

type Inscripciones struct {
	IdInscripcion    int64     `gorm:"primaryKey;autoIncrement;column:inscripcion_id"`
	IdUsuario        int64     `gorm:"not null" binding:"required"`                   // Foreign Key de `users`
	IdCurso          string    `gorm:"type:varchar(100);not null" binding:"required"` // ID del curso (proveniente de MongoDB)
	FechaInscripcion time.Time `gorm:"not null"`

	User Users `gorm:"foreignKey:IdUsuario"`
}

type Inscripcioness []Inscripciones
