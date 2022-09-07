package handler

import (
	"clinicapp/pkg/listing"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// func DoctorsRoutes(ls listing.Service) *gin.Engine {
// 	router := gin.Default()

// 	router.Group("/doctors")
// 	{
// 		// router.GET("/", getAllDoctors(ls))
// 		router.GET("/:id", GetDoctor(ls))
// 	}

// 	return router
// }

// a handler for GET /doctors/:id requests
func GetDoctor(ls listing.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		doctor_id, _ := strconv.Atoi(ctx.Params.ByName("id"))

		doctor, err := ls.GetDoctor(doctor_id)
		if err == listing.ErrIdNotFound {
			ctx.Error(errors.New("Get Doctor - " + listing.ErrIdNotFound.Error()))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": doctor,
		})
	}

}

func GetAllDoctors(ls listing.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		doctors := ls.GetAllDoctors()
		ctx.JSON(http.StatusOK, gin.H{
			"data": doctors,
		})
	}
}
