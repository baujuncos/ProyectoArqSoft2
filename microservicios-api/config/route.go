package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"microservicios-api/controllers"
	"os"
	"time"
)

var router *gin.Engine

func init() {
	// Configuración del logger
	log.SetOutput(os.Stdout) // Enviar los logs a la salida estándar
	log.SetLevel(log.DebugLevel)
	log.Info("Starting logger system")

	// Inicialización del router Gin
	router = gin.Default()

	// Configuración de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func StartRoute() {
	// Configuración de rutas
	router.GET("/microservices", controllers.GetServices)

	// Mensaje de inicio del servidor
	log.Info("Starting server on port 8004")
	err := router.Run(":8004")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
