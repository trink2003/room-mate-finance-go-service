package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/big"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
	"time"
)

func (h Handler) AddNewRoom(c *gin.Context) {

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

func (h Handler) GetListOfRooms(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)

	if !isSuccess {
		return
	}

	requestPayload := payload.PageRequestBody{}
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

	if requestPayload.Request.PageSize == 0 || requestPayload.Request.PageNumber == 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Page number or page size can not be 0",
			),
		)
		return
	}

	limit := requestPayload.Request.PageSize
	offset := requestPayload.Request.PageSize * (requestPayload.Request.PageNumber - 1)
	var room []model.Rooms
	var total int64 = 0

	transactionResultError := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var countQueryResult = tx.WithContext(ctx).Model(&model.Rooms{}).Count(&total)

		if countQueryResult.Error != nil {
			return countQueryResult.Error
		}

		var getListOfRoomQueryResult = tx.WithContext(ctx).
			Preload("Users").
			Limit(limit).
			Offset(offset).
			Order(utils.SortMapToString(requestPayload.Request.Sort)).
			Find(&room)

		if getListOfRoomQueryResult.Error != nil {
			return getListOfRoomQueryResult.Error
		}

		return nil
	})

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

	var totalRecordsSelected = len(room)

	if totalRecordsSelected == 0 {
		c.JSON(
			http.StatusOK,
			utils.ReturnPageResponse(
				c,
				constant.Success,
				0,
				0,
				room,
			),
		)
		return
	}
	totalBigNumber := new(big.Float).SetInt64(total)
	totalRecordsSelectedBigNumber := new(big.Float).SetInt64(int64(totalRecordsSelected))
	totalPage := new(big.Float).Quo(totalRecordsSelectedBigNumber, totalBigNumber)
	utils.RoundHalfUpBigFloat(totalPage)
	totalPageInt, _ := totalPage.Int(nil)

	c.JSON(
		http.StatusOK,
		utils.ReturnPageResponse(
			c,
			constant.Success,
			total,
			totalPageInt.Int64(),
			room,
		),
	)
}

func (h Handler) DeleteRoom(c *gin.Context) {
	ctx, isSuccess := utils.PrepareContext(c)

	if !isSuccess {
		return
	}

	requestPayload := payload.DeleteRoomRequestBody{}
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

	if requestPayload.Request.RoomCode == "ADMIN_ROOM" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DeleteDefaultRoomError,
				nil,
			),
		)
		return
	}

	var errorEnum = constant.Success
	var transactionResultError = h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var numberOfMemberInTargetRoom int64

		tx.Raw(
			`
			select
				count(1)
			from
				users u
				left join rooms r on r.id = u.room_id
			where
				u.room_id is not null
				and r.room_code = ?
			`,
			requestPayload.Request.RoomCode,
		).
			Scan(&numberOfMemberInTargetRoom)

		if numberOfMemberInTargetRoom > 0 {
			errorEnum = constant.RoomStillHavePeople
			return nil
		}

		var roomObjectResult = model.Rooms{}

		var roomQueryResult = tx.WithContext(ctx).
			Clauses(
				clause.Locking{
					Strength: clause.LockingStrengthUpdate,
				},
			).
			Where(
				model.Rooms{
					RoomCode: requestPayload.Request.RoomCode,
				},
			).Find(&roomObjectResult)

		if roomQueryResult.Error != nil {
			return roomQueryResult.Error
		}

		if roomObjectResult.BaseEntity.Id == 0 {
			errorEnum = constant.RoomDoesNotExist
			return nil
		}

		var deleteRoomQueryResult = tx.Delete(roomObjectResult)

		if deleteRoomQueryResult.Error != nil {
			return deleteRoomQueryResult.Error
		}

		return nil
	})

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

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			"ok",
		),
	)

}

func (h Handler) EditRoomName(c *gin.Context) {
	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}

	requestPayload := payload.EditRoomNameRequestBody{}
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

	var errorEnum = constant.Success
	var transactionResultError = h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var roomObjectResult = model.Rooms{}

		var roomQueryResult = tx.WithContext(ctx).
			/*Clauses(
				clause.Locking{
					Strength: clause.LockingStrengthUpdate,
				},
			).*/
			Where(
				model.Rooms{
					RoomCode: requestPayload.Request.RoomCode,
				},
			).Find(&roomObjectResult)

		if roomQueryResult.Error != nil {
			return roomQueryResult.Error
		}

		if roomObjectResult.BaseEntity.Id == 0 {
			errorEnum = constant.RoomDoesNotExist
			return nil
		}

		roomObjectResult.RoomName = requestPayload.Request.RoomName

		var updateRoomQueryResult = updateRoom(tx, &roomObjectResult, ctx)

		if updateRoomQueryResult.Error != nil {
			return updateRoomQueryResult.Error
		}

		return nil
	})

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

func updateRoom(db *gorm.DB, model *model.Rooms, ctx context.Context) *gorm.DB {
	var currentUsernameInsertOrUpdateData = utils.GetCurrentUsernameFromContextForInsertOrUpdateDataInDb(ctx)
	model.BaseEntity.Active = utils.GetPointerOfAnyValue(true)
	model.BaseEntity.UpdatedAt = time.Now()
	model.BaseEntity.UpdatedBy = currentUsernameInsertOrUpdateData

	return db.WithContext(ctx).Save(model)
}
