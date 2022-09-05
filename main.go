package main

import (
	"clinicapp/pkg/booking"
	"clinicapp/pkg/config"
	"clinicapp/pkg/handler"
	"clinicapp/pkg/listing"
	_ "clinicapp/pkg/migration"
	"clinicapp/pkg/storage/postgres"

	_ "github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
}

func main() {
	// migration.MigrateDown()
	// migration.MigrateUp()

	s, _ := postgres.NewStorage()
	lister := listing.NewService(s)
	booker := booking.NewService(s)
	handler.Handler(lister, booker)
}
