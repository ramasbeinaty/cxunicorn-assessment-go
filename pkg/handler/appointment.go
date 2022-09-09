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
		// var hi = ctx.PostForm("data")
		// var bye = ctx.PostForm("patient_id")

		var appointment booking.Appointment

		if err := ctx.BindJSON(&appointment); err != nil {
			fmt.Println("ERROR: Create Appointment - ", err)
		}

		// _appointment := booking.Appointment{}
		// _appointment.PatientID, _ = strconv.Atoi(ctx.PostForm("patient_id"))
		// _appointment.DoctorID, _ = strconv.Atoi(ctx.PostForm("doctor_id"))
		// _appointment.CreatedBy, _ = strconv.Atoi(ctx.PostForm("created_by"))
		// _appointment.StartDatetime, _ = time.Parse(time.RFC822, ctx.PostForm(("start_datetime")))
		// _appointment.EndDatetime, _ = time.Parse(time.RFC822, ctx.PostForm(("end_datetime")))

		bs.CreateAppointment(appointment)

	}

}

func GetAllAppointmentsOfDoctor(ls listing.Service) gin.HandlerFunc {
	var appointments []listing.Appointment

	return func(ctx *gin.Context) {
		// _doctor_id := ctx.Query("doctor_id")
		// date:= ctx.Query("date")

		// if _doctor_id != "" && date != "" {
		// 	doctor_id, _ := strconv.Atoi(_doctor_id)
		// 	appointments = ls.GetAllAppointmentsOfDoctor(doctor_id, date)
		// }
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
			fmt.Println("ERROR: Edit Appointment - ", err)
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
