package main

import (
	"time"
)

// User1 [...]
type User1 struct {
	ID        int       `gorm:"primaryKey;column:id;type:int;not null" json:"-"`
	Name      string    `gorm:"column:name;type:varchar(50);not null;default:''" json:"name"`
	Age       int       `gorm:"column:age;type:int;not null;default:0" json:"age"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName get sql table name.获取数据库表名
func (m *User1) TableName() string {
	return "user1"
}

// User1Columns get sql column name.获取数据库列名
var User1Columns = struct {
	ID        string
	Name      string
	Age       string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Name:      "name",
	Age:       "age",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
