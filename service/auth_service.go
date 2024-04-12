package service

import (
	context2 "context"
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

func (h AuthHandler) AddNewUser(ginContext *gin.Context) {

	context := context2.Background()

	requestPayload := payload.UserRegisterRequestBody{}

	if err := ginContext.ShouldBindJSON(&requestPayload); err != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["JSON_BINDING_ERROR"],
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
		Limit(1).
		Where(
			"username = ? AND active is true OR active is false", requestPayload.Request.Username,
		).
		Find(&userInDatabase)

	if userInDatabaseQueryResult.Error != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["QUERY_ERROR"],
				nil,
				userInDatabaseQueryResult.Error.Error(),
			),
		)
		return
	}

	if userInDatabase.BaseEntity.Id != 0 {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["CREATE_DUPLICATE_USER"],
				nil,
			),
		)
		return
	}

	if encryptPasswordError := utils.EncryptPasswordPointer(&requestPayload.Request.Password); encryptPasswordError != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["DATA_FORMAT_ERROR"],
				nil,
				encryptPasswordError.Error(),
			),
		)
		return
	}

	var user = model.Users{
		Username: requestPayload.Request.Username,
		Password: requestPayload.Request.Password,
		UserUid:  uuid.New().String(),
	}

	if result := SaveNewUser(h.DB, &user, context); result != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["QUERY_ERROR"],
				nil,
				result.Error(),
			),
		)
		return
	}

	ginContext.JSON(
		http.StatusOK,
		ReturnResponse(
			ginContext,
			constant.ErrorConstant["SUCCESS"],
			user,
		),
	)
}

func (h AuthHandler) Login(ginContext *gin.Context) {
	context := context2.Background()

	requestPayload := &payload.UserLoginRequestBody{}

	if err := ginContext.ShouldBindJSON(&requestPayload); err != nil {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["JSON_BINDING_ERROR"],
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
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["QUERY_ERROR"],
				nil,
				userInDatabaseQueryResult.Error.Error(),
			),
		)
		return
	}

	if userInDatabase.BaseEntity.Id == 0 {
		ginContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["AUTHENTICATE_FAILURE"],
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
			ReturnResponse(
				ginContext,
				constant.ErrorConstant["AUTHENTICATE_FAILURE"],
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
		ReturnResponse(
			ginContext,
			constant.ErrorConstant["SUCCESS"],
			response,
		),
	)

}

func SaveNewUser(db *gorm.DB, user *model.Users, ctx context2.Context) error {
	user.BaseEntity.Active = utils.GetPointerOfAnyValue(true)
	user.BaseEntity.UUID = uuid.New().String()
	user.BaseEntity.CreatedAt = time.Now()
	user.BaseEntity.UpdatedAt = time.Now()
	user.BaseEntity.CreatedBy = ctx.Value("username").(string)
	user.BaseEntity.UpdatedBy = ctx.Value("username").(string)

	saveFunc := func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			// return any error will rollback
			return err
		}
		return nil
	}
	return db.WithContext(ctx).Transaction(saveFunc)
}
