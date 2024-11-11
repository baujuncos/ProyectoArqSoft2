package main

import (
	controllers "cursos-api/controllers"
	queues "cursos-api/clients"
	repositories "cursos-api/repositories"
	services "cursos-api/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Mongo
	mainRepository := repositories.NewMongo(repositories.MongoConfig{
		Host:       "localhost",
		Port:       "27017",
		Username:   "root",
		Password:   "44898366",
		Database:   "cursos-api",
		Collection: "courses",
	})

	// Rabbit
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "localhost",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "some-rabbit",
	})

	// Services
	service := services.NewService(mainRepository, eventsQueue)

	// Controllers
	controller := controllers.NewController(service)

	// Router
	router := gin.Default()
	router.GET("/courses/:id", controller.GetCourseByID)
	router.POST("/courses", controller.Create)
	router.PUT("/courses/:id", controller.Update)
	router.DELETE("/courses/:id", controller.Delete)
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("error running application: %w", err)
	}
}
