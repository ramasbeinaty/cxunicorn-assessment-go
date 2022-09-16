package main

import (
	"clinicapp/pkg/auth"
	"clinicapp/pkg/booking"
	"clinicapp/pkg/canceling"
	"clinicapp/pkg/config"
	"clinicapp/pkg/deleting"
	"clinicapp/pkg/editing"
	"clinicapp/pkg/handler"
	"clinicapp/pkg/listing"
	"clinicapp/pkg/storage/postgres"

	_ "github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
}

func main() {
	// migration.MigrateDown()
	// migration.MigrateUp()

	// loc, err := time.LoadLocation("UTC")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// time.Local = loc

	s, _ := postgres.NewStorage()
	lister := listing.NewService(s)
	booker := booking.NewService(s)
	canceler := canceling.NewService(s)
	deleter := deleting.NewService(s)
	editer := editing.NewService(s)
	authenticator := auth.NewService(s)

	handler.Handler(lister, booker, canceler, deleter, editer, authenticator)
}
