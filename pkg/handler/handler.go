package handler

import (
	"clinicapp/pkg/booking"
	"clinicapp/pkg/listing"

	"github.com/gin-gonic/gin"
)

func Handler(ls listing.Service, bs booking.Service) {

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
			appointmentRoute.GET("/", GetAllAppointmentsOfDoctor(ls))
		}
	}

	router.Run()
}
