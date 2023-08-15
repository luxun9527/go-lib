package main

import (
	"context"
	"database/sql/driver"
	"go-lib/sdk/gorm/gentool/dao/query"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
	"log"
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

/*
	Scan(src interface{}) error   // sql.Scanner
	Value() (driver.Value, error) // driver.Valuer
*/
// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}

type SpecInt32 int64

func (s SpecInt32) Value() (driver.Value, error) {
	return int64(s), nil
}
func (s *SpecInt32) Scan(src interface{}) error {
	*s = SpecInt32(src.(int64))
	return nil
}
func TestGenCondition(t *testing.T) {
	InitGorm()
	u := query.Use(db).User
	s := SpecInt32(1)
	gen.Cond()
	//like := clause.Like{Column: u.Age.RawExpr(), Value: "value"}

	result, err := u.WithContext(context.Background()).Where(field.NewField(u.TableName(), u.Age.ColumnName().String()).Like(s)).First()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", result)

}
