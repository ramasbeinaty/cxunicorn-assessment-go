package handler

import (
	"clinicapp/pkg/auth"
	"clinicapp/pkg/booking"
	"clinicapp/pkg/canceling"
	"clinicapp/pkg/deleting"
	"clinicapp/pkg/editing"
	"clinicapp/pkg/listing"
	"clinicapp/pkg/logging"
	"clinicapp/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Handler(ls listing.Service, bs booking.Service, cs canceling.Service,
	ds deleting.Service, es editing.Service, as auth.Service, tc *logging.Telemetry) {

	router := gin.Default()

	superRoute := router.Group("/api")
	{

		authRoute := superRoute.Group("/auth")
		{
			authRoute.Use(middleware.HandleRequestTelemetry(tc))
			authRoute.POST("/register", CreateUser(as))
			authRoute.POST("/login", LoginUser(as))
		}

		doctorRoute := superRoute.Group("/doctors")
		{
			doctorRoute.Use(middleware.HandleRequestTelemetry(tc), middleware.AuthenticateUser(as))
			doctorRoute.GET("/:id", GetDoctor(ls))
			doctorRoute.GET("/:id/slots", GetAvailableSlotsOfDoctor(ls))
			doctorRoute.GET("/", GetAllDoctors(ls))
		}

		appointmentRoute := superRoute.Group("/appointments")
		{
			appointmentRoute.Use(middleware.HandleRequestTelemetry(tc), middleware.AuthenticateUser(as))
			appointmentRoute.POST("/", middleware.AuthorizeUser(as, auth.Roles.Patient), CreateAppointment(bs))
			appointmentRoute.GET("/doctors/:id", GetAllAppointmentsOfDoctor(ls))
			appointmentRoute.PATCH("/:id", CancelAppointment(cs))
			appointmentRoute.DELETE("/:id", DeleteAppointment(ds))
			appointmentRoute.PUT("/:id", EditAppointment(es))
		}
	}

	router.Run()
}
