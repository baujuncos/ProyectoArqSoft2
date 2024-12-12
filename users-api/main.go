package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"users-api/config"
	controllersUser "users-api/controllers/users"
	"users-api/internal/tokenizers"
	repositories "users-api/repositories/users"
	services "users-api/services/users"
)

func main() {
	// MySQL
	mySQLRepo := repositories.NewMySQL(
		repositories.MySQLConfig{
			Host:     config.MySQLHost,
			Port:     config.MySQLPort,
			Database: config.MySQLDatabase,
			Username: config.MySQLUsername,
			Password: config.MySQLPassword,
		},
	)

	// Cache
	cacheRepo := repositories.NewCache(repositories.CacheConfig{
		TTL: config.CacheDuration,
	})

	// Memcached
	memcachedRepo := repositories.NewMemcached(repositories.MemcachedConfig{
		Host: config.MemcachedHost,
		Port: config.MemcachedPort,
	})

	// Tokenizer
	jwtTokenizer := tokenizers.NewTokenizer(
		tokenizers.JWTConfig{
			Key:      config.JWTKey,
			Duration: config.JWTDuration,
		},
	)

	// Services de users
	serviceUsers := services.NewService(mySQLRepo, cacheRepo, memcachedRepo, jwtTokenizer)

	// Service de inscripciones
	serviceInscripciones := services.NewServiceIns(mySQLRepo)

	// Handlers de user e inscripciones
	controllerUser := controllersUser.NewController(serviceUsers)

	inscripcionesController := controllersUser.NewControllerIns(serviceInscripciones)

	// Create router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// URL mappings
	router.GET("/users/:id", controllerUser.GetUsersByID)
	router.POST("/users", controllerUser.CreateUser)
	router.POST("/login", controllerUser.Login)

	router.POST("/inscripciones", inscripcionesController.InsertInscripcion)
	router.GET("/inscripciones/usuario/:user_id", inscripcionesController.GetInscripcionesByUserID)
	router.GET("/inscripciones/curso/:course_id", inscripcionesController.GetInscripcionesByCursoID)

	// Run application
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running application: %v", err)
	}
}
