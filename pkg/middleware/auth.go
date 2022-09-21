package middleware

import (
	"clinicapp/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  http.StatusUnauthorized,
		"message": "Unauthorized access",
	})
}

func AuthorizeUser(as auth.Service, authorizedRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userRole := ctx.Keys["Roles"].(string)

		isAuthorized, _ := as.AuthorizeUser(userRole, authorizedRole)

		if !isAuthorized {
			// ctx.JSON(http.StatusUnauthorized, gin.H{
			// 	"response": "User is unauthorized - " + err.Error(),
			// })
			// return

			ReturnUnauthorized(ctx)
		}

		ctx.JSON(http.StatusAccepted, gin.H{
			"response": "User is authorized",
		})

	}
}

func AuthenticateUser(as auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenStr string = ""

		tokenStr = ctx.Request.Header.Get("Token")

		if tokenStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "AuthenticateUser - Token header is not found",
			})
			return
		}

		isValid, claims := as.VerifyJWT(tokenStr)

		if !isValid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"response": "AuthenticateUser - Token is invalid",
			})
			return
		}

		if len(ctx.Keys) == 0 {
			ctx.Keys = make(map[string]interface{})
		}

		ctx.Keys["Role"] = claims.Role
		ctx.Keys["Email"] = claims.Email
		ctx.Keys["UserID"] = claims.UserID
		ctx.Keys["UserName"] = claims.Name

	}
}
