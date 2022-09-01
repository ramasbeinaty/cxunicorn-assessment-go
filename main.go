package main

import (
	"clinicapp/pkg/config"
	"clinicapp/pkg/handler"
	"clinicapp/pkg/listing"
	"clinicapp/pkg/storage/postgres"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
}

func main() {
	// migration.MigrateUp()
	// migration.MigrateDown()

	s, _ := postgres.NewStorage()
	lister := listing.NewService(s)
	router := gin.Default()

	superRoute := router.Group("/api")
	{
		doctorRoute := superRoute.Group("/doctors")
		{
			doctorRoute.GET("/:id", handler.GetDoctor(lister))
		}
	}

	router.Run()
}
