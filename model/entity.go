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
	BaseEntity     BaseEntity       `gorm:"embedded" json:"baseInfo"`
	Username       string           `json:"username" gorm:"column:username;"`
	Password       string           `json:"-" gorm:"column:password;"`
	UserUid        string           `json:"userUid" gorm:"column:user_uid;"`
	ListOfExpenses []ListOfExpenses `gorm:"foreignKey:UsersID"`
	UserToPaid     []DebitUser      `gorm:"foreignKey:UserToPaidID"`
	PaidToUser     []DebitUser      `gorm:"foreignKey:PaidToUserID"`
}

type ListOfExpenses struct {
	BaseEntity BaseEntity `gorm:"embedded" json:"baseInfo"`
	Purpose    string     `json:"purpose" gorm:"column:purpose"`
	Amount     float64    `json:"amount" gorm:"column:amount"`
	UsersID    int        `json:"-"`
	Users      Users      `json:"boughtByUser" gorm:"column:bought_by_user;references:id"`
}

type DebitUser struct {
	BaseEntity   BaseEntity `gorm:"embedded" json:"baseInfo"`
	UserToPaidID int        `json:"-" gorm:"column:user_to_paid"`
	PaidToUserID int        `json:"-" gorm:"column:paid_to_user"`
	UserToPaid   Users      `json:"userToPaid" gorm:"column:user_to_paid;references:id"`
	PaidToUser   Users      `json:"paidToUser" gorm:"column:paid_to_user;references:id"`
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
