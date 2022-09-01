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
			doctorRoute.GET("/", GetDoctor(ls))
		}
	}

	router.Run()
}

// func Handler(l listing.Service) *gin.Engine {
// 	router := gin.Default()

// 	superGroup := router.Group("/api")
// 	{
// 		doctorsGroup := superGroup.Group("/doctors")
// 		{
// 			doctorsGroup.Handlers
// 		}
// 	}

// 	// router.GET("/doctors", getAllDoctors(l))
// 	router.GET("/doctors/:id", Get)

// 	return router
// }
