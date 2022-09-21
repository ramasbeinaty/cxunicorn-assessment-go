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

func LoginUser(as auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user auth.UserLogin

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to login user - " + err.Error(),
			})
			return
		}

		token, err := as.LoginUser(user)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to login user - " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusAccepted, gin.H{
			"response": "Successfully logged in user",
			"token":    token,
		})

	}
}

func CreateUser(as auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user auth.UserRegister

		if err := ctx.BindJSON(&user); err != nil {
			fmt.Println("ERROR: CreateUser - " + err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "Failed to create user",
			})
			return
		}

		_pass, err := bcrypt.GenerateFromPassword([]byte(user.UserDetails.Password), bcrypt.DefaultCost)

		if err != nil {
			fmt.Println("ERROR: CreateUser - password encryption failed - " + err.Error())

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

		ctx.JSON(http.StatusCreated, gin.H{
			"response": "Successfully created user",
		})

	}

}
