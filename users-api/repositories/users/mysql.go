package users

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
	dao "users-api/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

import _ "github.com/go-sql-driver/mysql"

type MySQLConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type MySQL struct {
	db *gorm.DB
}

var sqlDB *sql.DB

func NewMySQL(config MySQLConfig) MySQL {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)

	// Open connection to MySQL using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %s", err.Error())
	}

	// Obtener la conexión SQL nativa de gorm
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB from gorm: ", err)
	}

	//Migramos

	userTableExists := db.Migrator().HasTable(&dao.Users{})

	inscripcionesTableExists := db.Migrator().HasTable(&dao.Inscripciones{})

	if !userTableExists && !inscripcionesTableExists {
		err := db.Migrator().CreateTable(&dao.Users{})
		if err != nil {
			log.Fatal("Failed to migrate User table: ", err)
		}

		err = db.Migrator().CreateTable(&dao.Inscripciones{})
		if err != nil {
			log.Fatal("Failed to migrate Inscripcion table: ", err)
		}

		// Sembrar datos iniciales
		SeedDB(db)
	}

	return MySQL{db: db}
}

func SeedDB(db *gorm.DB) MySQL {
	userss := []dao.Users{
		{Email: "pauliortiz@example.com", Password: "contraseña1", Nombre: "paulina", Apellido: "ortiz", Admin: true},
		{Email: "baujuncos@example.com", Password: "contraseña2", Nombre: "bautista", Apellido: "juncos", Admin: true},
		{Email: "belutreachi2@example.com", Password: "contraseña3", Nombre: "belen", Apellido: "treachi", Admin: false},
		{Email: "virchurodiguez@example.com", Password: "contraseña4", Nombre: "virginia", Apellido: "rodriguez", Admin: false},
		{Email: "johndoe@example.com", Password: "contraseña5", Nombre: "John", Apellido: "Doe", Admin: false},
		{Email: "alicesmith@example.com", Password: "contraseña6", Nombre: "Alice", Apellido: "Smith", Admin: true},
		{Email: "bobjohnson@example.com", Password: "contraseña7", Nombre: "Bob", Apellido: "Johnson", Admin: false},
		{Email: "janedoe@example.com", Password: "contraseña8", Nombre: "Jane", Apellido: "Doe", Admin: false},
		{Email: "emilywilliams@example.com", Password: "contraseña9", Nombre: "Emily", Apellido: "Williams", Admin: true},
	}

	for i, users := range userss {
		// Hashear la contraseña con MD5
		hasher := md5.New()
		hasher.Write([]byte(users.Password))
		hashedPassword := hex.EncodeToString(hasher.Sum(nil))
		userss[i].Password = hashedPassword
		db.FirstOrCreate(&userss[i], dao.Users{Email: users.Email})
	}

	inscripcioness := []dao.Inscripciones{
		{IdUsuario: 3, IdCurso: "1", FechaInscripcion: time.Now()},
		{IdUsuario: 5, IdCurso: "2", FechaInscripcion: time.Now()},
		{IdUsuario: 6, IdCurso: "3", FechaInscripcion: time.Now()},
	}

	for _, inscripcion := range inscripcioness {
		db.FirstOrCreate(&inscripcion, dao.Inscripciones{IdUsuario: inscripcion.IdUsuario, IdCurso: inscripcion.IdCurso})
	}

	return MySQL{
		db: db,
	}
}

func (repository MySQL) GetUserByID(id int64) (dao.Users, error) {
	var user dao.Users
	if err := repository.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error fetching user by id: %w", err)
	}
	return user, nil
}

func (repository MySQL) GetUserByEmail(email string) (dao.Users, error) {
	var user dao.Users
	if err := repository.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("error fetching user by email: %w", err)
	}
	return user, nil
}

func (repository MySQL) CreateUser(user dao.Users) (int64, error) {
	// Usamos GORM para crear el registro
	result := repository.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	// Retornamos el ID generado automáticamente
	return user.User_id, nil
}

func (repository MySQL) InsertInscripcion(inscripcion dao.Inscripciones) (int64, error) {
	if err := repository.db.Create(&inscripcion).Error; err != nil {
		log.Printf("Error al guardar la inscripción: %v\n", err)
		return 0, fmt.Errorf("error al guardar la inscripción: %w", err)
	}

	return inscripcion.IdInscripcion, nil

}

func (repository MySQL) IsSubscribed(userID int64, cursoID string) (bool, error) {
	var count int64
	result := repository.db.Model(&dao.Inscripciones{}). //con model nos aseguramos de buscar en tabla inscripciones
								Where("id_usuario = ? AND id_curso = ?", userID, cursoID). //filtramos donde coincidan los parametros dados
								Count(&count)                                              //contamos coincidencias

	if result.Error != nil {
		return false, result.Error
	}
	if count > 0 {
		return true, nil
	} //retornamos true en caso de coincidencia y false en caso de no...
	return false, nil
}

func (repository MySQL) GetInscripcionesByUserID(idUsuario int64) ([]string, error) {
	var inscripciones []dao.Inscripciones
	var cursos []string

	// Buscar inscripciones por el ID del usuario
	if err := repository.db.Where("id_usuario = ?", idUsuario).Find(&inscripciones).Error; err != nil {
		return nil, fmt.Errorf("error fetching inscripciones by user ID: %w", err)
	}

	// Extraer los IDs de los cursos de las inscripciones encontradas
	for _, inscripcion := range inscripciones {
		cursos = append(cursos, inscripcion.IdCurso)
	}

	return cursos, nil
}

func (repository MySQL) GetInscripcionesByCursoID(idCurso string) ([]int64, error) {
	var inscripciones []dao.Inscripciones
	var usuarios []int64

	// Buscar inscripciones por el ID del curso
	if err := repository.db.Where("id_curso = ?", idCurso).Find(&inscripciones).Error; err != nil {
		return nil, fmt.Errorf("error fetching inscripciones by course ID: %w", err)
	}

	// Si no se encontraron inscripciones, devolver un error personalizado
	/*if len(inscripciones) == 0 {
		return nil, fmt.Errorf("no se encontraron inscripciones para el curso con ID %s", idCurso)
	}*/

	// Extraer los IDs de los usuarios de las inscripciones encontradas
	for _, inscripcion := range inscripciones {
		usuarios = append(usuarios, inscripcion.IdUsuario)
	}

	return usuarios, nil
}
