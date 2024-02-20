package model

import "gorm.io/plugin/soft_delete"

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID        uint32                `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username  string                `gorm:"column:username;not null" json:"username"`
	Password  string                `gorm:"column:password;not null" json:"password"`
	CreatedAt uint32                `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt uint32                `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:unix" json:"deleted_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
