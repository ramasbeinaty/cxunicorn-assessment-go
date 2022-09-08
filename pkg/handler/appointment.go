package handler

import (
	"clinicapp/pkg/booking"
	"clinicapp/pkg/listing"
	"fmt"
	"net/http"

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
			"data": appointments,
		})
	}
}
