package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"room-mate-finance-go-service/config"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/service"
	"room-mate-finance-go-service/utils"
	"time"
)

func main() {

	router := gin.Default()

	router.NoRoute(
		func(context *gin.Context) {
			context.JSON(
				http.StatusNotFound, &payload.Response{
					Trace:        utils.GetTraceId(context),
					ErrorCode:    constant.ErrorConstant["PAGE_NOT_FOUND"].ErrorCode,
					ErrorMessage: constant.ErrorConstant["PAGE_NOT_FOUND"].ErrorMessage,
				},
			)
		},
	)

	router.NoMethod(
		func(context *gin.Context) {
			context.JSON(
				http.StatusNotFound, &payload.Response{
					Trace:        utils.GetTraceId(context),
					ErrorCode:    constant.ErrorConstant["METHOD_NOT_ALLOWED"].ErrorCode,
					ErrorMessage: constant.ErrorConstant["METHOD_NOT_ALLOWED"].ErrorMessage,
				},
			)
		},
	)

	applicationPort := os.Getenv("APPLICATION_PORT")
	if applicationPort == "" {
		applicationPort = "8080"
	}

	db, err := config.InitDatabaseConnection()
	if err != nil {
		panic(err)
	}

	service.RegisterRoutes(router, db)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"port": applicationPort,
		})
	})

	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
				AllowHeaders:     []string{"Origin"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				AllowOriginFunc: func(origin string) bool {
					return origin == "*"
				},
				MaxAge: 12 * time.Hour,
			},
		),
	)

	ginErr := router.Run(":" + applicationPort)
	if ginErr != nil {
		panic(ginErr)
	}
}
