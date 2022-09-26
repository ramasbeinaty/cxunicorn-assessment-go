package handler

import (
	"clinicapp/pkg/booking"
	"clinicapp/pkg/canceling"
	"clinicapp/pkg/deleting"
	"clinicapp/pkg/editing"
	"clinicapp/pkg/listing"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateAppointment(bs booking.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var appointment booking.Appointment

		if err := ctx.BindJSON(&appointment); err != nil {
			fmt.Println("ERROR: CreateAppointment - ", err.Error())
			return
		}

		if err := bs.CreateAppointment(appointment); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"response": "Successfully created the appointment",
		})

	}

}

func GetAllAppointmentsOfDoctor(ls listing.Service) gin.HandlerFunc {
	var appointments []listing.Appointment

	return func(ctx *gin.Context) {
		_doctorID, _ := strconv.Atoi(ctx.Params.ByName("id"))

		var appointmentsRequest listing.AppointmentsRequest

		if err := ctx.BindJSON(&appointmentsRequest); err != nil {
			fmt.Println("ERROR: GetAvailableSlotsOfDoctor - ", err.Error())
			return

		}

		appointments = ls.GetAllAppointmentsOfDoctor(_doctorID, appointmentsRequest.Date)

		ctx.JSON(http.StatusOK, gin.H{
			"response": appointments,
		})
	}
}

func CancelAppointment(cs canceling.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appointment_id, _ := strconv.Atoi(ctx.Params.ByName("id"))

		err := cs.CancelAppointment(appointment_id)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to cancel appointment",
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Successfully canceled appointment",
			})
		}

	}
}

func DeleteAppointment(ds deleting.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appointment_id, _ := strconv.Atoi(ctx.Params.ByName("id"))

		err := ds.DeleteAppointment(appointment_id)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to delete appointment",
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Successfully delete appointment",
			})
		}

	}
}

func EditAppointment(es editing.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// bind the appointment data received
		var appointment editing.Appointment

		if err := ctx.BindJSON(&appointment); err != nil {
			fmt.Println("ERROR: Edit Appointment - ", err.Error())
		}

		// get the appointment id from the url
		appointment_id, _ := strconv.Atoi(ctx.Params.ByName("id"))

		err := es.EditAppointment(appointment_id, appointment)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to edited appointment",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"response": "Successfully edited appointment",
			})
		}

	}
}
