package model

import "time"

type BaseEntity struct {
	Id        int64     `json:"id" gorm:"column:id;primaryKey;"`
	Active    bool      `json:"active" gorm:"column:active;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;"`
	CreatedBy string    `json:"createdBy" gorm:"column:created_by;"`
	UpdatedBy string    `json:"updatedBy" gorm:"column:updated_by;"`
}

type Users struct {
	BaseEntity BaseEntity `gorm:"embedded" json:"baseInfo"`
	Username   string     `json:"username" gorm:"column:username;"`
	Password   string     `json:"-" gorm:"column:password;"`
	UserUid    string     `json:"userUid" gorm:"column:user_uid;"`
}

type Tabler interface {
	TableName() string
}

func (Users) TableName() string {
	return "users"
}
