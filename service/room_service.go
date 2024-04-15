package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
	"time"
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

		var roomQueryResult = tx.WithContext(ctx).Where(
			model.Rooms{
				RoomCode: requestPayload.Request.RoomCode,
			},
		).Find(&listOfRoomThatHaveTheSameRoomCode)

		if roomQueryResult.Error != nil {
			return roomQueryResult.Error
		}
		if len(listOfRoomThatHaveTheSameRoomCode) > 0 {
			errorEnum = constant.RoomHasBeenExisted
			return nil
		}
		var newRoom = model.Rooms{
			RoomCode: requestPayload.Request.RoomCode,
			RoomName: requestPayload.Request.RoomName,
		}

		var insertResult = saveNewRoom(tx, &newRoom, ctx)
		if insertResult.Error != nil {
			return insertResult.Error
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

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			"ok",
		),
	)

}

func saveNewRoom(db *gorm.DB, model *model.Rooms, ctx context.Context) *gorm.DB {
	var currentUsernameInsertOrUpdateData = ""
	var usernameFromContext = ctx.Value("username")
	if usernameFromContext != nil {
		currentUsernameInsertOrUpdateData = usernameFromContext.(string)
	}
	model.BaseEntity.Active = utils.GetPointerOfAnyValue(true)
	model.BaseEntity.UUID = uuid.New().String()
	model.BaseEntity.CreatedAt = time.Now()
	model.BaseEntity.UpdatedAt = time.Now()
	model.BaseEntity.CreatedBy = currentUsernameInsertOrUpdateData
	model.BaseEntity.UpdatedBy = currentUsernameInsertOrUpdateData

	return db.WithContext(ctx).Create(model)
}
