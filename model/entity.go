package model

import (
	"time"
)

type BaseEntity struct {
	Id        int64     `json:"id" gorm:"column:id;primaryKey;not null"`
	UUID      string    `json:"uuid" gorm:"column:uuid;not null"`
	Active    *bool     `json:"active" gorm:"column:active;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;not null"`
	CreatedBy string    `json:"createdBy" gorm:"column:created_by;not null"`
	UpdatedBy string    `json:"updatedBy" gorm:"column:updated_by;not null"`
}

type Rooms struct {
	BaseEntity BaseEntity `gorm:"embedded" json:"baseInfo"`
	RoomName   string     `json:"roomName" gorm:"column:room_name;"`
	RoomCode   string     `json:"roomCode" gorm:"column:room_code;unique;not null"`
	Users      []Users    `json:"usersInRoom" gorm:"foreignKey:RoomsID"`
}

type Users struct {
	BaseEntity     BaseEntity       `gorm:"embedded" json:"baseInfo"`
	Username       string           `json:"username" gorm:"column:username;"`
	Password       string           `json:"-" gorm:"column:password;"`
	UserUid        string           `json:"userUid" gorm:"column:user_uid;"`
	RoomsID        int64            `json:"-" gorm:"column:room_id;references:id"`
	Rooms          Rooms            `json:"room" gorm:"->;<-:false;-:migration;foreignKey:Id"`
	ListOfExpenses []ListOfExpenses `json:"listOfExpenses" gorm:"foreignKey:BoughtByUserID"`
	UserToPaid     []DebitUser      `json:"userToPaid" gorm:"foreignKey:UserToPaidID"`
	PaidToUser     []DebitUser      `json:"paidToUser" gorm:"foreignKey:PaidToUserID"`
	UsersRoles     []UsersRoles     `json:"usersRoles" gorm:"foreignKey:UsersId;"`
}

type Roles struct {
	BaseEntity BaseEntity   `gorm:"embedded" json:"baseInfo"`
	RoleName   string       `json:"roleName" gorm:"column:role_name;"`
	UsersRoles []UsersRoles `json:"usersRoles" gorm:"foreignKey:RolesId;"`
}

type UsersRoles struct {
	BaseEntity BaseEntity `gorm:"embedded" json:"baseInfo"`
	UsersId    int64      `json:"-" gorm:"column:users_id;references:id"`
	Users      Users      `json:"users" gorm:"->;<-:false;-:migration;foreignKey:Id"`
	RolesId    int64      `json:"-" gorm:"column:roles_id;references:id"`
	Roles      Roles      `json:"roles" gorm:"->;<-:false;-:migration;foreignKey:Id"`
}

type ListOfExpenses struct {
	BaseEntity     BaseEntity  `gorm:"embedded" json:"baseInfo"`
	Purpose        string      `json:"purpose" gorm:"column:purpose"`
	Amount         float64     `json:"amount" gorm:"column:amount"`
	BoughtByUserID int64       `json:"-" gorm:"column:bought_by_user;references:id"`
	Users          Users       `json:"expenseOwner" gorm:"->;<-:false;-:migration;foreignKey:Id"`
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
func (Roles) TableName() string {
	return "roles"
}
func (Rooms) TableName() string {
	return "rooms"
}
