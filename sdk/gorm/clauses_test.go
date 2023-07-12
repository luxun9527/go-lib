package main

import (
	"gorm.io/gorm/clause"
	"testing"
)

func TestClauses(t *testing.T) {
	u := &User{
		ID:   0,
		Age:  120,
		Name: "0",
	}
	InitGorm()
	db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)
}

type User struct {
	ID   int    `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	Age  int    `gorm:"column:age;" json:"age"`
	Name string `gorm:"column:name;" json:"name"`
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}
