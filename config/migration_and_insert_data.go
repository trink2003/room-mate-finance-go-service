package config

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"room-mate-finance-go-service/model"
	"room-mate-finance-go-service/utils"
	"time"
)

func MigrationAndInsertDate(db *gorm.DB) {
	usersMigrateErr := db.AutoMigrate(&model.Users{})
	if usersMigrateErr != nil {
		panic(usersMigrateErr)
	}

	debitUserMigrateErr := db.AutoMigrate(&model.DebitUser{})
	if debitUserMigrateErr != nil {
		panic(debitUserMigrateErr)
	}

	listOfExpensesMigrateErr := db.AutoMigrate(&model.ListOfExpenses{})
	if listOfExpensesMigrateErr != nil {
		panic(listOfExpensesMigrateErr)
	}

	rolesMigrateErr := db.AutoMigrate(&model.Roles{})
	if rolesMigrateErr != nil {
		panic(rolesMigrateErr)
	}

	usersRolesMigrateErr := db.AutoMigrate(&model.UsersRoles{})
	if usersRolesMigrateErr != nil {
		panic(usersRolesMigrateErr)
	}
}

func InsertData(db *gorm.DB) {
	pass, _ := utils.EncryptPassword("admin")

	var room = model.Rooms{
		BaseEntity: model.BaseEntity{
			Active:    utils.GetPointerOfAnyValue(true),
			UUID:      uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: "system",
			UpdatedBy: "system",
		},
		RoomCode: "ADMIN_ROOM",
		RoomName: "Default room for administrator user",
	}

	db.Save(&room)

	adminUser := model.Users{
		BaseEntity: model.BaseEntity{
			Active:    utils.GetPointerOfAnyValue(true),
			UUID:      uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: "system",
			UpdatedBy: "system",
		},
		Username: "admin",
		Password: pass,
		UserUid:  uuid.New().String(),
		RoomsID:  room.BaseEntity.Id,
	}

	db.Create(&adminUser)

	roles := []model.Roles{
		{
			BaseEntity: model.BaseEntity{
				Active:    utils.GetPointerOfAnyValue(true),
				UUID:      uuid.New().String(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				CreatedBy: "system",
				UpdatedBy: "system",
			},
			RoleName: "ADMIN",
		},
		{
			BaseEntity: model.BaseEntity{
				Active:    utils.GetPointerOfAnyValue(true),
				UUID:      uuid.New().String(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				CreatedBy: "system",
				UpdatedBy: "system",
			},
			RoleName: "USER",
		},
	}

	db.Create(&roles)

	db.Create(&model.UsersRoles{
		BaseEntity: model.BaseEntity{
			Active:    utils.GetPointerOfAnyValue(true),
			UUID:      uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: "system",
			UpdatedBy: "system",
		},
		UsersId: adminUser.BaseEntity.Id,
		RolesId: roles[0].BaseEntity.Id,
	})
}
