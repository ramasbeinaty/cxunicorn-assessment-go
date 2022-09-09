package handler

import (
	"clinicapp/pkg/booking"
	"clinicapp/pkg/canceling"
	"clinicapp/pkg/deleting"
	"clinicapp/pkg/editing"
	"clinicapp/pkg/listing"

	"github.com/gin-gonic/gin"
)

func Handler(ls listing.Service, bs booking.Service, cs canceling.Service,
	ds deleting.Service, es editing.Service) {

	router := gin.Default()

	superRoute := router.Group("/api")
	{
		doctorRoute := superRoute.Group("/doctors")
		{
			doctorRoute.GET("/:id", GetDoctor(ls))
			doctorRoute.GET("/", GetAllDoctors(ls))
		}

		appointmentRoute := superRoute.Group("/appointments")
		{
			appointmentRoute.POST("/", CreateAppointment(bs))
			// appointmentRoute.GET("/", GetAllAppointmentsOfDoctor(ls))
			appointmentRoute.PATCH("/:id", CancelAppointment(cs))
			appointmentRoute.DELETE("/:id", DeleteAppointment(ds))
			appointmentRoute.PUT("/:id", EditAppointment(es))
		}
	}

	router.Run()
}
