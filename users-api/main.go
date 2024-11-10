package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"users-api/config"
	controllers "users-api/controllers/users"
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

	// Services
	service := services.NewService(mySQLRepo, cacheRepo, memcachedRepo, jwtTokenizer)

	// Handlers
	controller := controllers.NewController(service)

	// Create router
	router := gin.Default()

	// URL mappings
	router.GET("/users/:id", controller.GetUsersByID)
	router.POST("/users", controller.Create)
	router.POST("/login", controller.Login)

	// Run application
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running application: %v", err)
	}
}
