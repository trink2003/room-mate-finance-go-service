package user

import (
	context2 "context"
	"net/http"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRegisterRequestBodyValue struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterRequestBody struct {
	Request UserRegisterRequestBodyValue `json:"request"`
}

func (h UserHandler) GetUsers(ginContext *gin.Context) {

	context := context2.Background()
	traceId := uuid.New().String()

	context = context2.WithValue(context, "traceId", traceId)

	var user []model.Users

	if result := h.DB.Find(&user); result.Error != nil {
		ginContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"trace": utils.GetTraceId(ginContext),
			"error": result.Error,
		})
		return
	}

	ginContext.JSON(http.StatusOK, &user)
}

func (h UserHandler) AddNewUser(ginContext *gin.Context) {

	context := context2.Background()
	traceId := uuid.New().String()

	context = context2.WithValue(context, "traceId", traceId)

	requestPayload := UserRegisterRequestBody{}

	if err := ginContext.ShouldBindJSON(&requestPayload); err != nil {
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"trace": utils.GetTraceId(ginContext),
			"error": err.Error(),
		})
		return
	}
	context = context2.WithValue(context, "username", requestPayload.Request.Username)

	var userInDatabase = model.Users{}

	userInDatabaseQueryResult := h.DB.Where("username = ? AND active is true OR active is false", requestPayload.Request.Username).Find(&userInDatabase)

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

	if result := SaveNewUser(h.DB, &user, context); result.Error != nil {
		ginContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"trace": utils.GetTraceId(ginContext),
			"error": result.Error,
		})
		return
	}
	ginContext.JSON(http.StatusOK, user)
}

func SaveNewUser(db *gorm.DB, user *model.Users, ctx context2.Context) *gorm.DB {
	user.BaseEntity.Active = true
	user.BaseEntity.CreatedAt = time.Now()
	user.BaseEntity.UpdatedAt = time.Now()
	user.BaseEntity.CreatedBy = ctx.Value("username").(string)
	user.BaseEntity.UpdatedBy = ctx.Value("username").(string)
	return db.Save(user)
}
