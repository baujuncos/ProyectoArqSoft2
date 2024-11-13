package main

import (
	queues "cursos-api/clients"
	controllers "cursos-api/controllers"
	repositories "cursos-api/repositories"
	services "cursos-api/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	// Mongo
	mainRepository := repositories.NewMongo(repositories.MongoConfig{
		Host:       "mongo",
		Port:       "27017",
		Username:   "root",
		Password:   "belusql1",
		Database:   "cursos-api",
		Collection: "courses",
	})

	// Rabbit
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "courses-news",
	})

	// Services
	service := services.NewService(mainRepository, eventsQueue)

	// Controllers
	controller := controllers.NewController(service)

	// Router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/courses/:id", controller.GetCourseByID)
	router.GET("/courses", controller.GetCourses)
	router.POST("/courses", controller.Create)
	router.PUT("/courses/:id", controller.Update)
	router.DELETE("/courses/:id", controller.Delete)

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("error running application: %w", err)
	}
}
