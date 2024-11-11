package users

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
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

	if !userTableExists {
		err := db.Migrator().CreateTable(&dao.Users{})
		if err != nil {
			log.Fatal("Failed to migrate User table: ", err)
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
