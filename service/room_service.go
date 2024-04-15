package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
)

func (h RoomHandler) AddNewRoom(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)

	if !isSuccess {
		return
	}

	requestPayload := payload.AddNewRoomRequestBody{}
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.JsonBindingError,
				nil,
				err.Error(),
			),
		)
		return
	}

	if requestPayload.Request.RoomName == "" {
		requestPayload.Request.RoomName = fmt.Sprintf("Room %v", requestPayload.Request.RoomCode)
	}
	var errorEnum = constant.ErrorEnums{}
	transactionResultError := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var listOfRoomThatHaveTheSameRoomCode []model.Rooms
		tx.WithContext(ctx).Where(
			model.Rooms{
				RoomCode: requestPayload.Request.RoomCode,
			},
		).Find(&listOfRoomThatHaveTheSameRoomCode)
		if len(listOfRoomThatHaveTheSameRoomCode) > 0 {
			errorEnum = constant.RoomHasBeenExisted
			return nil
		}
		return nil
	})

	if errorEnum.ErrorCode != 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				errorEnum,
				nil,
			),
		)
		return
	}

	if transactionResultError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.QueryError,
				nil,
				transactionResultError.Error(),
			),
		)
		return
	}

}
