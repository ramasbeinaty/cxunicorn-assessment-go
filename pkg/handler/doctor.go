package handler

import (
	"clinicapp/pkg/listing"
	"errors"
	"fmt"
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
		doctorID, _ := strconv.Atoi(ctx.Params.ByName("id"))

		doctor, err := ls.GetDoctor(doctorID)
		if err == listing.ErrIdNotFound {
			ctx.Error(errors.New("Get Doctor - " + listing.ErrIdNotFound.Error()))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"response": doctor,
		})
	}

}

// a handler for GET /doctors requests
func GetAllDoctors(ls listing.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		doctors := ls.GetAllDoctors()
		ctx.JSON(http.StatusOK, gin.H{
			"response": doctors,
		})
	}
}

func GetAvailableSlotsOfDoctor(ls listing.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		doctorID, _ := strconv.Atoi(ctx.Params.ByName("id"))

		var doctorSlots listing.DoctorSlots

		if err := ctx.BindJSON(&doctorSlots); err != nil {
			fmt.Println("ERROR: GetAvailableSlotsOfDoctor - ", err.Error())
			return
		}

		availableSlots := ls.GetAvailableSlotsPerDay(doctorID, doctorSlots.SlotsDate)
		ctx.JSON(http.StatusOK, gin.H{
			"response": availableSlots,
		})
	}
}
