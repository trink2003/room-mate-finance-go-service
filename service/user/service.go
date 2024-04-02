package user

import (
	context2 "context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/utils"
	"time"
)

type UserRegisterRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h UserHandler) GetUser(ginContext *gin.Context) {

	context := context2.Background()
	traceId := uuid.New().String()

	context = context2.WithValue(context, "traceId", traceId)

	var user []model.User

	if result := h.DB.Find(&user); result.Error != nil {
		ginContext.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	ginContext.JSON(http.StatusOK, &user)
}

func (h UserHandler) AddNewUser(ginContext *gin.Context) {

	context := context2.Background()
	traceId := uuid.New().String()

	context = context2.WithValue(context, "traceId", traceId)

	requestPayload := UserRegisterRequestBody{}

	if err := ginContext.BindJSON(&requestPayload); err != nil {
		ginContext.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if encryptPasswordError := utils.EncryptPasswordPointer(&requestPayload.Password); encryptPasswordError != nil {
		ginContext.AbortWithError(http.StatusInternalServerError, encryptPasswordError)
		return
	}

	var user = model.User{
		Username: requestPayload.Username,
		Password: requestPayload.Password,
	}

	if result := SaveNewUser(h.DB, &user, context); result.Error != nil {
		ginContext.AbortWithError(http.StatusInternalServerError, result.Error)
		return
	}
	ginContext.JSON(http.StatusOK, user)
}

func SaveNewUser(db *gorm.DB, user *model.User, ctx context2.Context) *gorm.DB {
	user.BaseEntity.Active = true
	user.BaseEntity.CreatedAt = time.Now()
	user.BaseEntity.UpdatedAt = time.Now()
	user.BaseEntity.CreatedBy = ctx.Value("username").(string)
	user.BaseEntity.UpdatedBy = ctx.Value("username").(string)
	return db.Save(user)
}
