package model

import (
	"time"
)

type BaseEntity struct {
	Id        int64     `json:"id" gorm:"column:id;primaryKey;"`
	Active    bool      `json:"active" gorm:"column:active;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;"`
	CreatedBy string    `json:"createdBy" gorm:"column:created_by;"`
	UpdatedBy string    `json:"updatedBy" gorm:"column:updated_by;"`
}

type Users struct {
	BaseEntity     BaseEntity       `gorm:"embedded" json:"baseInfo"`
	Username       string           `json:"username" gorm:"column:username;"`
	Password       string           `json:"-" gorm:"column:password;"`
	UserUid        string           `json:"userUid" gorm:"column:user_uid;"`
	ListOfExpenses []ListOfExpenses `json:"listOfExpenses" gorm:"foreignKey:BoughtByUserID"`
	UserToPaid     []DebitUser      `json:"userToPaid" gorm:"foreignKey:UserToPaidID"`
	PaidToUser     []DebitUser      `json:"paidToUser" gorm:"foreignKey:PaidToUserID"`
}

type ListOfExpenses struct {
	BaseEntity     BaseEntity  `gorm:"embedded" json:"baseInfo"`
	Purpose        string      `json:"purpose" gorm:"column:purpose"`
	Amount         float64     `json:"amount" gorm:"column:amount"`
	BoughtByUserID int64       `json:"-" gorm:"column:bought_by_user;references:id"`
	Users          Users       `json:"users" gorm:"->;<-:false;-:migration;foreignKey:Id"`
	DebitUser      []DebitUser `json:"debitUser" gorm:"foreignKey:ListOfExpensesID"`
}

type DebitUser struct {
	BaseEntity       BaseEntity `gorm:"embedded" json:"baseInfo"`
	Amount           float64    `json:"amount" gorm:"column:amount"`
	ListOfExpensesID int64      `json:"-" gorm:"column:expense;references:id"`
	UserToPaidID     int64      `json:"-" gorm:"column:user_to_paid;references:id"`
	PaidToUserID     int64      `json:"-" gorm:"column:paid_to_user;references:id"`
	UserToPaid       Users      `json:"-" gorm:"->;<-:false;-:migration;foreignKey:Id"`
	PaidToUser       Users      `json:"-" gorm:"->;<-:false;-:migration;foreignKey:Id"`
}

type Tabler interface {
	TableName() string
}

func (Users) TableName() string {
	return "users"
}

func (ListOfExpenses) TableName() string {
	return "list_of_expenses"
}

func (DebitUser) TableName() string {
	return "debit_user"
}
