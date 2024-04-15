package utils

import (
	"context"
	"github.com/google/uuid"
	"room-mate-finance-go-service/model"
	"time"
)

func GenerateNewBaseEntity(ctx context.Context) model.BaseEntity {
	var currentUsernameInsertOrUpdateData = GetCurrentUsernameFromContextForInsertOrUpdateDataInDb(ctx)
	return model.BaseEntity{
		Active:    GetPointerOfAnyValue(true),
		UUID:      uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: currentUsernameInsertOrUpdateData,
		UpdatedBy: currentUsernameInsertOrUpdateData,
	}
}
