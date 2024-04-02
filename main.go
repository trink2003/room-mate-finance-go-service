package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"room-mate-finance-go-service/config"
	"room-mate-finance-go-service/service/user"
)

func main() {
	applicationPort := os.Getenv("APPLICATION_PORT")
	if applicationPort == "" {
		applicationPort = "8080"
	}

	db, err := config.InitDatabaseConnection()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	user.RegisterRoutes(router, db)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"port": applicationPort,
		})
	})

	ginErr := router.Run(":" + applicationPort)
	if ginErr != nil {
		panic(ginErr)
	}
}
