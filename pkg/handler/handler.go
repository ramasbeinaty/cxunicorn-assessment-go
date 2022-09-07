package handler

import (
	"clinicapp/pkg/listing"

	"github.com/gin-gonic/gin"
)

func Handler(ls listing.Service) {

	router := gin.Default()

	superRoute := router.Group("/api")
	{
		doctorRoute := superRoute.Group("/doctors")
		{
			doctorRoute.GET("/:id", GetDoctor(ls))
			doctorRoute.GET("/", GetAllDoctors(ls))
		}

		// appointmentRoute := superRoute.Group("/appointments")
		// {
		// 	appointmentRoute.POST("/", CreateAppointment(bs))
		// }
	}

	router.Run()
}
