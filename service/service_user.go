package service

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/log"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h Handler) GetAllActiveUser(c *gin.Context) {

	currentUser, isCurrentUserExist := utils.GetCurrentUsername(c)

	if isCurrentUserExist != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			utils.ReturnResponse(
				c,
				constant.Unauthorized,
				nil,
				isCurrentUserExist.Error(),
			),
		)
		return
	}

	ctx := context.Background()

	ctx = context.WithValue(ctx, constant.UsernameLogKey, *currentUser)
	ctx = context.WithValue(ctx, constant.TraceIdLogKey, utils.GetTraceId(c))

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
	var user []model.Users
	var total int64 = 0

	transactionResult := h.DB.Transaction(func(tx *gorm.DB) error {

		tx.WithContext(ctx).Model(&model.Users{}).
			Where(
				tx.Where(
					model.Users{
						BaseEntity: model.BaseEntity{
							Active: utils.GetPointerOfAnyValue(true),
						},
					},
				).Where("active is not null"),
			).Count(&total)

		tx.WithContext(ctx).Preload("Rooms").Limit(limit).
			Offset(offset).
			Order(utils.SortMapToString(requestPayload.Request.Sort)).
			Where(
				tx.Where(
					model.Users{
						BaseEntity: model.BaseEntity{
							Active: utils.GetPointerOfAnyValue(true),
						},
					},
				).Where("active is not null"),
			).Find(&user)
		return nil
	})
	if transactionResult != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.QueryError,
				nil,
				transactionResult.Error(),
			),
		)
		return
	}

	totalRecordsSelected := len(user)

	if totalRecordsSelected == 0 {
		c.JSON(
			http.StatusOK,
			utils.ReturnPageResponse(
				c,
				constant.Success,
				0,
				0,
				user,
			),
		)
		return
	}

	/*
		if result := h.DB.Where("active is not null AND active is true ORDER BY id DESC").Find(&user); result.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &payload.Response{
				Trace:        utils.GetTraceId(c),
				ErrorCode:    constant.QueryError.ErrorCode,
				ErrorMessage: constant.QueryError.ErrorMessage + result.Error.Error(),
			})
			return
		}
	*/

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
			user,
		),
	)
}

func (h Handler) GetMemberInRoom(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)

	if !isSuccess {
		return
	}

	var currentUserModel = model.Users{}
	var currentUsername = utils.GetCurrentUsernameFromContextForInsertOrUpdateDataInDb(ctx)

	h.DB.WithContext(ctx).Preload("Rooms").Where(
		model.Users{
			Username: currentUsername,
			BaseEntity: model.BaseEntity{
				Active: utils.GetPointerOfAnyValue(true),
			},
		},
	).Find(&currentUserModel)

	if currentUserModel.BaseEntity.Id == 0 {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			utils.ReturnResponse(
				c,
				constant.UserNotExisted,
				nil,
				"We can not determine who are you in the current session",
			),
		)
		return
	}

	var allActiveMemberInRoom = make([]model.Users, 0)

	h.DB.WithContext(ctx).Preload("Rooms").Where(
		h.DB.Where(
			model.Users{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
				RoomsID: currentUserModel.RoomsID,
			},
		).Where("id not in (?)", currentUserModel.BaseEntity.Id),
	).Find(&allActiveMemberInRoom)

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			allActiveMemberInRoom,
		),
	)
}

func (h Handler) GetMemberInASpecificRoomCode(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)

	if !isSuccess {
		return
	}

	requestPayload := payload.GetMemberInARoomRequestBody{}
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
	var currentUsername = utils.GetCurrentUsernameFromContextForInsertOrUpdateDataInDb(ctx)

	var currentUserModel = model.Users{}

	h.DB.WithContext(ctx).Where(
		model.Users{
			Username: currentUsername,
			BaseEntity: model.BaseEntity{
				Active: utils.GetPointerOfAnyValue(true),
			},
		},
	).Find(&currentUserModel)

	if currentUserModel.BaseEntity.Id == 0 {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			utils.ReturnResponse(
				c,
				constant.UserNotExisted,
				nil,
				"We can not determine who are you in the current session",
			),
		)
		return
	}

	roomModelObject := model.Rooms{}

	h.DB.WithContext(ctx).Where(model.Rooms{RoomCode: requestPayload.Request.RoomCode}).Find(&roomModelObject)

	if roomModelObject.BaseEntity.Id == 0 {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			utils.ReturnResponse(
				c,
				constant.RoomDoesNotExist,
				nil,
			),
		)
		return
	}

	var allActiveMemberInRoom = make([]model.Users, 0)

	h.DB.WithContext(ctx).Preload("Rooms").Where(
		h.DB.Where(
			model.Users{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
				RoomsID: roomModelObject.BaseEntity.Id,
			},
		).Where("id not in (?)", currentUserModel.BaseEntity.Id),
	).Find(&allActiveMemberInRoom)

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			allActiveMemberInRoom,
		),
	)
}

func (h Handler) RemoveMemberInARoom(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)

	if !isSuccess {
		return
	}

	requestPayload := payload.RemoveMemberInARoomBody{}
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

	transactionResultError := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		log.WithLevel(
			constant.Info,
			ctx,
			fmt.Sprintf(
				"find room with code %v",
				requestPayload.Request.RoomCode,
			),
		)
		roomObjectFoundFromRoomCode := model.Rooms{}

		var queryRoomError = tx.Where(
			model.Rooms{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
				RoomCode: requestPayload.Request.RoomCode,
			},
		).Find(&roomObjectFoundFromRoomCode)

		if roomObjectFoundFromRoomCode.BaseEntity.Id == 0 || queryRoomError.Error != nil {
			errorEnum = constant.RoomDoesNotExist
			return queryRoomError.Error
		}

		log.WithLevel(
			constant.Info,
			ctx,
			fmt.Sprintf(
				"find user with id %v",
				requestPayload.Request.UserId,
			),
		)

		userObjectFoundFromUserIdInRequest := model.Users{}

		var queryUserInDatabase = tx.Where(
			model.Users{
				BaseEntity: model.BaseEntity{
					Id:     requestPayload.Request.UserId,
					Active: utils.GetPointerOfAnyValue(true),
				},
				RoomsID: roomObjectFoundFromRoomCode.BaseEntity.Id,
			},
		).Find(&userObjectFoundFromUserIdInRequest)

		if userObjectFoundFromUserIdInRequest.BaseEntity.Id == 0 || queryUserInDatabase.Error != nil {
			errorEnum = constant.UserNotExisted
			return queryUserInDatabase.Error
		}

		if userObjectFoundFromUserIdInRequest.Username == ctx.Value(constant.UsernameLogKey).(string) {
			errorEnum = constant.ActionCannotPerformOnYourself
			return nil
		}

		userObjectFoundFromUserIdInRequest.BaseEntity.Active = utils.GetPointerOfAnyValue(false)
		utils.ChangeDataOfBaseEntityForUpdate(ctx, &userObjectFoundFromUserIdInRequest.BaseEntity)
		updateUserError := tx.Save(&userObjectFoundFromUserIdInRequest)
		if updateUserError != nil {
			return updateUserError.Error
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

	c.JSON(http.StatusOK, utils.ReturnResponse(c, errorEnum, nil))
}

func (h Handler) AddMemberToARoom(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)

	if !isSuccess {
		return
	}

	requestPayload := payload.AddMemberToARoomBody{}
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

	transactionResultError := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		log.WithLevel(
			constant.Info,
			ctx,
			fmt.Sprintf(
				"find room with code %v",
				requestPayload.Request.RoomCode,
			),
		)
		roomObjectFoundFromRoomCode := model.Rooms{}

		var queryRoomError = tx.Where(
			model.Rooms{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
				RoomCode: requestPayload.Request.RoomCode,
			},
		).Find(&roomObjectFoundFromRoomCode)

		if roomObjectFoundFromRoomCode.BaseEntity.Id == 0 || queryRoomError.Error != nil {
			errorEnum = constant.RoomDoesNotExist
			return queryRoomError.Error
		}

		log.WithLevel(
			constant.Info,
			ctx,
			fmt.Sprintf(
				"find user with id %v",
				requestPayload.Request.UserId,
			),
		)

		userObjectFoundFromUserIdInRequest := model.Users{}

		var queryUserInDatabase = tx.Where(
			model.Users{
				BaseEntity: model.BaseEntity{
					Id:     requestPayload.Request.UserId,
					Active: utils.GetPointerOfAnyValue(false),
				},
				RoomsID: roomObjectFoundFromRoomCode.BaseEntity.Id,
			},
		).Find(&userObjectFoundFromUserIdInRequest)

		if userObjectFoundFromUserIdInRequest.BaseEntity.Id == 0 || queryUserInDatabase.Error != nil {
			errorEnum = constant.UserNotExisted
			return queryUserInDatabase.Error
		}

		if userObjectFoundFromUserIdInRequest.Username == ctx.Value(constant.UsernameLogKey).(string) {
			errorEnum = constant.ActionCannotPerformOnYourself
			return nil
		}

		userObjectFoundFromUserIdInRequest.BaseEntity.Active = utils.GetPointerOfAnyValue(true)
		utils.ChangeDataOfBaseEntityForUpdate(ctx, &userObjectFoundFromUserIdInRequest.BaseEntity)
		updateUserError := tx.Save(&userObjectFoundFromUserIdInRequest)
		if updateUserError != nil {
			return updateUserError.Error
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

	c.JSON(http.StatusOK, utils.ReturnResponse(c, errorEnum, nil))
}

func (h Handler) MoveAllMemberToANewRoom(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)

	if !isSuccess {
		return
	}

	requestPayload := payload.MoveAllMemberToANewRoomBody{}
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

	transactionResultError := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		log.WithLevel(
			constant.Info,
			ctx,
			fmt.Sprintf(
				"find room with old code %v",
				requestPayload.Request.OldRoomCode,
			),
		)
		roomObjectFoundFromOldRoomCode := model.Rooms{}

		var queryOldRoomError = tx.Where(
			model.Rooms{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
				RoomCode: requestPayload.Request.OldRoomCode,
			},
		).Find(&roomObjectFoundFromOldRoomCode)

		if roomObjectFoundFromOldRoomCode.BaseEntity.Id == 0 || queryOldRoomError.Error != nil {
			errorEnum = constant.RoomDoesNotExist
			errorEnum.ErrorMessage = "Invalid old room code"
			return queryOldRoomError.Error
		}
		log.WithLevel(
			constant.Info,
			ctx,
			fmt.Sprintf(
				"find room with new code %v",
				requestPayload.Request.NewRoomCode,
			),
		)
		roomObjectFoundFromNewRoomCode := model.Rooms{}

		var queryNewRoomError = tx.Where(
			model.Rooms{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
				RoomCode: requestPayload.Request.NewRoomCode,
			},
		).Find(&roomObjectFoundFromNewRoomCode)

		if roomObjectFoundFromNewRoomCode.BaseEntity.Id == 0 || queryNewRoomError.Error != nil {
			errorEnum = constant.RoomDoesNotExist
			errorEnum.ErrorMessage = "Invalid new room code"
			return queryNewRoomError.Error
		}

		log.WithLevel(
			constant.Info,
			ctx,
			"find all user from old room code",
		)

		var listUsersInRoom = []model.Users{}

		var listUserQueryError = tx.Where(
			model.Users{
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
				RoomsID: roomObjectFoundFromOldRoomCode.BaseEntity.Id,
			},
		).Find(&listUsersInRoom)

		if len(listUsersInRoom) < 1 || listUserQueryError.Error != nil {
			errorEnum = constant.EmptyRoomError
			return listUserQueryError.Error
		}

		for _, user := range listUsersInRoom {
			utils.ChangeDataOfBaseEntityForUpdate(ctx, &user.BaseEntity)
			user.RoomsID = roomObjectFoundFromNewRoomCode.BaseEntity.Id
		}

		var saveListOfUsersToNewRoomError = tx.Save(&listUsersInRoom)

		if saveListOfUsersToNewRoomError.Error != nil {
			errorEnum = constant.EmptyRoomError
			return saveListOfUsersToNewRoomError.Error
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

	c.JSON(http.StatusOK, utils.ReturnResponse(c, errorEnum, nil))
}
