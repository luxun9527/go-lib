package main

import (
	"context"
	"database/sql/driver"
	"github.com/spf13/cast"
	"go-lib/sdk/gorm/gentool/dao/query"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
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

type StringInt int64

func (s StringInt) Value() (driver.Value, error) {
	return cast.ToString(int64(s)), nil
}

type StringBinary string

func (s StringBinary) underlyingDB() *gorm.DB {
	return nil
}
func (s StringBinary) underlyingDO() *gorm.DB {
	return nil
}
func (s StringBinary) BeCond() interface{} {
	return nil
}
func (s StringBinary) CondError() error {
	return nil
}

type StringBinaryLike struct {
	gen.DO
	TableName string
	Column    string
	Value     string
}

func (s StringBinaryLike) BeCond() interface{} {
	return s
}

func (like StringBinaryLike) Build(builder clause.Builder) {
	builder.WriteQuoted(like.TableName)
	builder.WriteByte(46)
	builder.WriteQuoted(like.Column)
	builder.WriteString(" LIKE binary ")
	builder.AddVar(builder, like.Value)
}

func TestGenCondition1(t *testing.T) {
	InitGorm()
	dao := query.Use(db)
	u := dao.User
	s := StringInt(1)
	result, err := u.WithContext(context.Background()).Where(field.NewField(u.TableName(), u.Age.ColumnName().String()).Like(s)).First()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", result)

}
func TestGenCondition2(t *testing.T) {
	InitGorm()
	dao := query.Use(db)
	u := dao.User
	var s gen.SubQuery = &StringBinaryLike{
		TableName: "table",
		Column:    "column",
		Value:     "v1",
	}
	_, _ = u.WithContext(context.Background()).Where(s).Where(u.ID.Like(1)).First()
}
