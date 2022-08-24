package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/test", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "hi",
		})
	})

	server.Run(":8080")
}
