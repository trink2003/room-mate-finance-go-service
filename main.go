package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"room-mate-finance-go-service/config"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/log"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/service"
	"room-mate-finance-go-service/utils"
	"time"
)

func main() {
	var ctx = context.Background()
	router := gin.Default()

	router.NoRoute(
		func(context *gin.Context) {
			context.JSON(
				http.StatusNotFound, &payload.Response{
					Trace:        utils.GetTraceId(context),
					ErrorCode:    constant.PageNotFound.ErrorCode,
					ErrorMessage: constant.PageNotFound.ErrorMessage,
				},
			)
		},
	)

	router.NoMethod(
		func(context *gin.Context) {
			context.JSON(
				http.StatusNotFound, &payload.Response{
					Trace:        utils.GetTraceId(context),
					ErrorCode:    constant.MethodNotAllowed.ErrorCode,
					ErrorMessage: constant.MethodNotAllowed.ErrorMessage,
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

	router.GET(constant.BaseApiPath+"/", func(ctx *gin.Context) {
		ctx.Data(
			http.StatusOK,
			constant.ContentTypeHTML,
			[]byte(`
				<h1>Room mate finance service</h1><h3>This service is use for separate money for each daily expense for everybody in a room</h3>
			`),
		)
	})

	log.WithLevel(
		constant.Info,
		ctx,
		"Current directory is: "+utils.GetCurrentDirectory(),
	)

	router.GET(constant.BaseApiPath+"/favicon.ico", func(ctx *gin.Context) {
		ctx.Data(
			http.StatusOK,
			"image/x-icon",
			utils.ReadFileFromPath(utils.GetCurrentDirectory(), "icon", "favicon.ico"),
		)
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

	log.WithLevel(
		constant.Info,
		ctx,
		"Application starting",
	)

	ginErr := router.Run(":" + applicationPort)
	if ginErr != nil {
		panic(ginErr)
	}
}
