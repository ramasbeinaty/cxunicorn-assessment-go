package handler

import (
	"clinicapp/pkg/listing"
	"errors"
	"net/http"

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

// func GetAllDoctors(ls listing.Service) []listing.Doctor {
// 	var doctors []json.Doctor
// 	DB.Find(&doctors)

// 	return doctors
// }

// a handler for GET /doctors/:id requests
func GetDoctor(ls listing.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		doctor, err := ls.GetDoctor(ctx.Params.ByName("id"))
		if err == listing.ErrIdNotFound {
			ctx.Error(errors.New("Get Doctor - " + listing.ErrIdNotFound.Error()))
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": doctor,
		})
	}

}

// func GetDoctor(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"data": "getting the doc!",
// 	})

// }

// func GetDoctor(c *gin.Context, ls listing.Service) {
// 	doctor, err := ls.GetDoctor(c.Params.ByName("id"))
// 	if err == listing.ErrIdNotFound {
// 		log.Fatal("ERROR: Was not able to find doc with given id - ", listing.ErrIdNotFound)
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": doctor,
// 	})

// }
