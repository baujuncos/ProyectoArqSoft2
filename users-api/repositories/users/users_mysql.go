package users

import (
	"errors"
	"fmt"
	"log"
	dao "users-api/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

var (
	migrate = []interface{}{
		dao.Users{},
	}
)

func NewMySQL(config MySQLConfig) MySQL {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)

	// Open connection to MySQL using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %s", err.Error())
	}

	// Automigrate structs to Gorm
	for _, target := range migrate {
		if err := db.AutoMigrate(target); err != nil {
			log.Fatalf("error automigrating structs: %s", err.Error())
		}
	}

	return MySQL{
		db: db,
	}
}

func (repository MySQL) GetAll() ([]dao.Users, error) {
	var usersList []dao.Users
	if err := repository.db.Find(&usersList).Error; err != nil {
		return nil, fmt.Errorf("error fetching all users: %w", err)
	}
	return usersList, nil
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
	if err := repository.db.Create(&user).Error; err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	return user.User_id, nil
}
