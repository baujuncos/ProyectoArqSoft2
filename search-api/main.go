package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	queues "search-api/clients"
	controllers "search-api/controllers"
	repositories "search-api/repositories"
	services "search-api/services"
	"time"
)

func main() {
	// Solr
	solrRepo := repositories.NewSolr(repositories.SolrConfig{
		Host:       "solr",    // Solr host
		Port:       "8983",    // Solr port
		Collection: "courses", // Collection name
	})

	// Rabbit
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "courses-news",
	})

	// Hotels API
	cursosAPI := repositories.NewHTTP(repositories.HTTPConfig{
		Host: "cursos-api",
		Port: "8081",
	})

	// Services
	service := services.NewService(solrRepo, cursosAPI)

	// Controllers
	controller := controllers.NewController(service)

	// Launch rabbit consumer
	if err := eventsQueue.StartConsumer(service); err != nil {
		log.Fatalf("Error running consumer: %v", err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create router
	router.GET("/search", controller.Search)
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
