package handler

import (
	"clinicapp/pkg/auth"
	"clinicapp/pkg/booking"
	"clinicapp/pkg/canceling"
	"clinicapp/pkg/deleting"
	"clinicapp/pkg/editing"
	"clinicapp/pkg/listing"
	"clinicapp/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func Handler(ls listing.Service, bs booking.Service, cs canceling.Service,
	ds deleting.Service, es editing.Service, as auth.Service, tc appinsights.TelemetryClient) {

	router := gin.Default()

	superRoute := router.Group("/api")
	{

		authRoute := superRoute.Group("/auth")
		{
			authRoute.POST("/register", CreateUser(as))
			authRoute.POST("/login", handleRequestWithLog(LoginUser(as), tc))
		}

		doctorRoute := superRoute.Group("/doctors")
		doctorRoute.Use(middleware.AuthenticateUser(as))
		{
			doctorRoute.GET("/:id", GetDoctor(ls))
			doctorRoute.GET("/:id/slots", GetAvailableSlotsOfDoctor(ls))
			doctorRoute.GET("/", GetAllDoctors(ls))
		}

		appointmentRoute := superRoute.Group("/appointments")
		appointmentRoute.Use(middleware.AuthenticateUser(as))
		{
			appointmentRoute.POST("/", middleware.AuthorizeUser(as, auth.Roles.Patient), CreateAppointment(bs))
			appointmentRoute.GET("/doctors/:id", GetAllAppointmentsOfDoctor(ls))
			appointmentRoute.PATCH("/:id", CancelAppointment(cs))
			appointmentRoute.DELETE("/:id", DeleteAppointment(ds))
			appointmentRoute.PUT("/:id", EditAppointment(es))
		}
	}

	router.Run()
}
