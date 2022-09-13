package handler

import (
	"clinicapp/pkg/auth"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserRole struct {
	Role string `json:"role"`
}

func CreateUser(as auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// var _userRole UserRole

		// if err := ctx.BindJSON(&_userRole); err != nil {
		// 	fmt.Println("ERROR: CreateUser - ", err)
		// 	ctx.JSON(http.StatusBadRequest, gin.H{
		// 		"response": "Failed to get the user role",
		// 	})
		// }

		// if _userRole.Role == PatientRole {
		// 	var user auth.PatientRegister

		// 	as.CreatePatient(user)
		// }

		var user auth.UserRegister

		if err := ctx.BindJSON(&user); err != nil {
			fmt.Println("ERROR: CreateUser - ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to create user",
			})
			return
		}

		_pass, err := bcrypt.GenerateFromPassword([]byte(user.UserDetails.Password), bcrypt.DefaultCost)

		if err != nil {
			fmt.Println("ERROR: CreateUser - password encryption failed - ", err)

			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to create user - failed to encrypt password",
			})
			return
		}

		user.UserDetails.Password = string(_pass)

		if err = as.CreateUser(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to create user",
			})
			return
		}

		ctx.JSON(http.StatusAccepted, gin.H{
			"response": "Successfully created user",
		})

	}

}
