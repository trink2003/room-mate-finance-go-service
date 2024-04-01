package model

import "time"

type BaseEntity struct {
	Id        int64     `json:"id" gorm:"column:ID;primaryKey;"`
	Active    bool      `json:"active" gorm:"column:ACTIVE;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:CREATED_AT;"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:UPDATED_AT;"`
	CreatedBy string    `json:"createdBy" gorm:"column:CREATED_BY;"`
	UpdatedBy string    `json:"updatedBy" gorm:"column:UPDATED_BY;"`
}

type User struct {
	BaseEntity BaseEntity `gorm:"embedded"`
	Username   string     `json:"username" gorm:"column:USERNAME;"`
	Password   string     `json:"password" gorm:"column:PASSWORD;"`
}
