package dto

import "time"

type CursoDto struct {
	Course_id   string    `json:"course_id"`
	Nombre      string    `json:"nombre"`
	Profesor_id string    `json:"profesor_id"`
	Categoria   string    `json:"categoria"`
	Descripcion string    `json:"descripcion"`
	Valoracion  float64   `json:"valoracion"`
	Duracion    int       `json:"duracion"`
	Requisitos  string    `json:"requisitos"`
	UrlImage    string    `json:"url_image"`
	FechaInicio time.Time `json:"fecha_inicio"`
}

type InscripcionesDto struct {
	IdInscripcion    int       `json:"inscripcion_id,omitempty"`
	IdUser           int64     `json:"user_id"`
	IdCourse         string    `json:"course_id"`
	FechaInscripcion time.Time `json:"fecha_inscripcion"`
}
