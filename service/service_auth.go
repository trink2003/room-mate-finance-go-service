package service

import (
	context2 "context"
	"errors"
	"net/http"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h Handler) AddNewUser(ginContext *gin.Context) {

	context, isSuccess := utils.PrepareContext(ginContext)

	if !isSuccess {
		return
	}

	requestPayload := payload.UserRegisterRequestBody{}

	if err := ginContext.ShouldBindJSON(&requestPayload); err != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				ginContext,
				constant.JsonBindingError,
				nil,
				err.Error(),
			),
		)
		return
	}

	var userInDatabase = model.Users{}

	var roomObjectResult = model.Rooms{}

	h.DB.
		WithContext(context).
		Where(
			model.Rooms{
				RoomCode: requestPayload.Request.RoomCode,
			},
		).Find(&roomObjectResult)

	if roomObjectResult.BaseEntity.Id == 0 {
		ginContext.AbortWithStatusJSON(
			http.StatusNotFound,
			utils.ReturnResponse(
				ginContext,
				constant.RoomDoesNotExist,
				nil,
			),
		)
		return
	}

	userInDatabaseQueryResult := h.
		DB.
		WithContext(context).
		Limit(1).
		Where(
			"username = ? AND active is true OR active is false", requestPayload.Request.Username,
		).
		Find(&userInDatabase)

	if userInDatabaseQueryResult.Error != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				ginContext,
				constant.QueryError,
				nil,
				userInDatabaseQueryResult.Error.Error(),
			),
		)
		return
	}

	if userInDatabase.BaseEntity.Id != 0 {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				ginContext,
				constant.CreateDuplicateUser,
				nil,
			),
		)
		return
	}

	if encryptPasswordError := utils.EncryptPasswordPointer(&requestPayload.Request.Password); encryptPasswordError != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				ginContext,
				constant.DataFormatError,
				nil,
				encryptPasswordError.Error(),
			),
		)
		return
	}

	var user = model.Users{
		BaseEntity: utils.GenerateNewBaseEntity(context),
		Username:   requestPayload.Request.Username,
		Password:   requestPayload.Request.Password,
		UserUid:    uuid.New().String(),
		RoomsID:    roomObjectResult.BaseEntity.Id,
	}

	var errorEnum = constant.Success
	var transactionResultError = h.DB.WithContext(context).Transaction(func(tx *gorm.DB) error {

		var roleOfUser model.Roles

		tx.Where(
			model.Roles{
				RoleName: "USER",
			},
		).
			Find(&roleOfUser)

		if roleOfUser.BaseEntity.Id == 0 {
			errorEnum = constant.RoleNotExist
			return errors.New(constant.RoleNotExist.ErrorMessage)
		}

		saveNewUserQueryResult := tx.Save(&user)
		if saveNewUserQueryResult.Error != nil {
			return saveNewUserQueryResult.Error
		}

		var saveRoleForUserQueryResult = tx.Save(&model.UsersRoles{
			BaseEntity: utils.GenerateNewBaseEntity(context),
			UsersId:    user.BaseEntity.Id,
			RolesId:    roleOfUser.BaseEntity.Id,
		})

		if saveRoleForUserQueryResult.Error != nil {
			return saveRoleForUserQueryResult.Error
		}

		return nil
	})

	if transactionResultError != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				ginContext,
				constant.QueryError,
				nil,
				transactionResultError.Error(),
			),
		)
		return
	}

	if errorEnum.ErrorCode != 0 {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				ginContext,
				errorEnum,
				nil,
			),
		)
		return
	}

	ginContext.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			ginContext,
			constant.Success,
			user,
		),
	)
}

func (h Handler) Login(ginContext *gin.Context) {
	context := context2.Background()

	requestPayload := &payload.UserLoginRequestBody{}

	if err := ginContext.ShouldBindJSON(&requestPayload); err != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				ginContext,
				constant.JsonBindingError,
				nil,
				err.Error(),
			),
		)
		return
	}

	context = context2.WithValue(context, "username", requestPayload.Request.Username)
	context = context2.WithValue(context, "traceId", utils.GetTraceId(ginContext))

	var userInDatabase = model.Users{}

	userInDatabaseQueryResult := h.
		DB.
		WithContext(context).
		Where(
			"username = ? AND active is true", requestPayload.Request.Username,
		).
		Find(&userInDatabase)

	if userInDatabaseQueryResult.Error != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				ginContext,
				constant.QueryError,
				nil,
				userInDatabaseQueryResult.Error.Error(),
			),
		)
		return
	}

	if userInDatabase.BaseEntity.Id == 0 {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				ginContext,
				constant.AuthenticateFailure,
				nil,
				"username invalid",
			),
		)
		return
	}

	var role []string

	h.DB.WithContext(context).Raw(`
		select
			r.role_name
		from
			users u
			left join users_roles ur on ur.users_id = u.id
			left join roles r on r.id = ur.roles_id
		where
			u.user_uid = ?
    `, userInDatabase.UserUid).Scan(&role)

	if comparePasswordError := utils.ComparePassword(requestPayload.Request.Password, userInDatabase.Password); comparePasswordError != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				ginContext,
				constant.AuthenticateFailure,
				nil,
				"password invalid",
			),
		)
		return
	}

	token := utils.GenerateJwtToken(requestPayload.Request.Username, role...)

	response := &payload.UserLoginResponseBody{
		Token: token,
	}

	ginContext.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			ginContext,
			constant.Success,
			response,
		),
	)

}
