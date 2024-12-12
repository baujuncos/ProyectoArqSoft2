package dao

import "time"

type Course struct {
	Course_id    string    `bson:"_id,omitempty"`
	Nombre       string    `bson:"nombre"`
	Profesor_id  string    `bson:"profesor_id"`
	Categoria    string    `bson:"categoria"`
	Descripcion  string    `bson:"descripcion"`
	Valoracion   float64   `gorm:"valoracion"`
	Duracion     int       `gorm:"duracion"`
	Requisitos   string    `gorm:"requisitos"`
	Capacidad    int       `bson:"capacidad"` //Capacidad máxima del curso
	Url_image    string    `gorm:"url_image"`
	Fecha_inicio time.Time `gorm:"fecha_inicio"`
}
type Courses []Course
