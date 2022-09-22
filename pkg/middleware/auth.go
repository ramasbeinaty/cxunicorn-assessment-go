package middleware

import (
	"clinicapp/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func ReturnUnauthorized(ctx *gin.Context) {
// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 		"response": "Unauthorized access",
// 	})
// }

func AuthorizeUser(as auth.Service, authorizedRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userRole := ctx.Keys["Role"].(string)

		isAuthorized, _ := as.AuthorizeUser(userRole, authorizedRole)

		if !isAuthorized {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"response": "Unauthorized access",
			})
			return
		}

		// ctx.JSON(http.StatusAccepted, gin.H{
		// 	"response": "Authorized access",
		// })

	}
}

func AuthenticateUser(as auth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenStr string = ""

		tokenStr = ctx.Request.Header.Get("Token")

		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"response": "AuthenticateUser - Token header is not found",
			})
			return
		}

		isValid, claims := as.VerifyJWT(tokenStr)

		if !isValid {
			ctx.AbortWithStatusJSON(http.StatusAccepted, gin.H{
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
