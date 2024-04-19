package utils

import (
	"context"
	"net/http"
	"room-mate-finance-go-service/constant"

	"github.com/gin-gonic/gin"
)

func PrepareContext(c *gin.Context) (context.Context, bool) {
	ctx := context.Background()

	currentUser, isCurrentUserExist := GetCurrentUsername(c)

	if isCurrentUserExist != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			ReturnResponse(
				c,
				constant.Unauthorized,
				nil,
				isCurrentUserExist.Error(),
			),
		)
		return ctx, false
	}
	ctx = context.WithValue(ctx, "username", *currentUser)
	ctx = context.WithValue(ctx, "traceId", GetTraceId(c))

	return ctx, true
}

func GetCurrentUsernameFromContextForInsertOrUpdateDataInDb(ctx context.Context) string {
	var currentUsernameInsertOrUpdateData = ""
	var usernameFromContext = ctx.Value("username")
	if usernameFromContext != nil {
		currentUsernameInsertOrUpdateData = usernameFromContext.(string)
	}
	return currentUsernameInsertOrUpdateData
}
