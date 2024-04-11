package service

import (
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

type AuthHandler struct {
	DB *gorm.DB
}

type ExpenseHandler struct {
	DB *gorm.DB
}

type DebitHandler struct {
	DB *gorm.DB
}
