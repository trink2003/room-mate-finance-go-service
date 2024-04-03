package service

import (
	context2 "context"
	"net/http"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/payload"
	"room-mate-finance-go-service/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h UserHandler) GetUsers(ginContext *gin.Context) {

	var user []model.Users

	if result := h.DB.Where("active is not null AND active is true ORDER BY id DESC").Find(&user); result.Error != nil {
		ginContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"trace": utils.GetTraceId(ginContext),
			"error": result.Error,
		})
		return
	}

	ginContext.JSON(http.StatusOK, &user)
}

func (h AuthHandler) AddNewUser(ginContext *gin.Context) {

	context := context2.Background()
	traceId := uuid.New().String()

	context = context2.WithValue(context, "traceId", traceId)

	requestPayload := payload.UserRegisterRequestBody{}

	if err := ginContext.ShouldBindJSON(&requestPayload); err != nil {
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"trace": utils.GetTraceId(ginContext),
			"error": err.Error(),
		})
		return
	}
	context = context2.WithValue(context, "username", requestPayload.Request.Username)

	var userInDatabase = model.Users{}

	userInDatabaseQueryResult := h.
		DB.
		Limit(1).
		Where(
			"username = ? AND active is true OR active is false", requestPayload.Request.Username,
		).
		Find(&userInDatabase)

	if userInDatabaseQueryResult.Error != nil {
		ginContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"trace": utils.GetTraceId(ginContext),
			"error": userInDatabaseQueryResult.Error,
		})
		return
	}

	if userInDatabase.BaseEntity.Id != 0 {
		ginContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"trace": utils.GetTraceId(ginContext),
			"error": "user already exist in db",
		})
		return
	}

	if encryptPasswordError := utils.EncryptPasswordPointer(&requestPayload.Request.Password); encryptPasswordError != nil {
		ginContext.AbortWithStatusJSON(http.StatusInternalServerError, encryptPasswordError)
		return
	}

	var user = model.Users{
		Username: requestPayload.Request.Username,
		Password: requestPayload.Request.Password,
		UserUid:  uuid.New().String(),
	}

	if result := SaveNewUser(h.DB, &user, context); result != nil {
		ginContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"trace": utils.GetTraceId(ginContext),
			"error": result.Error,
		})
		return
	}
	ginContext.JSON(http.StatusOK, user)
}

func (h AuthHandler) Login(ginContext *gin.Context) {

}

func SaveNewUser(db *gorm.DB, user *model.Users, ctx context2.Context) error {
	user.BaseEntity.Active = true
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
	return db.Transaction(saveFunc)
}
