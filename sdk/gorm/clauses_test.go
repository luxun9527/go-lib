package main

import (
	"context"
	"database/sql/driver"
	"github.com/spf13/cast"
	"go-lib/sdk/gorm/gen/dao/query"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
	"log"
	"testing"
)

// Like whether string matches regular expression
type LikeBinary struct {
	Table string
	Column string
	Value interface{}

}
//实现build接口，可以自定义表达式
func (like LikeBinary) Build(builder clause.Builder) {
	builder.WriteQuoted(like.Table)
	builder.WriteByte(46)
	builder.WriteQuoted(like.Column)
	builder.WriteString(" LIKE BINARY ")
	builder.AddVar(builder, like.Value)
}


//演示gorm子句 https://gorm.io/zh_CN/docs/data_types.html
func TestGormClause(t *testing.T){
	var u User
	db.Where(clause.Like{Column: "c1", Value: "v1"}).Find(u)
	// SELECT * FROM `user` WHERE `c1` LIKE 'v1'
	db.Where(LikeBinary{
		Table:  "user",
		Column: "c1",
		Value:  1,
	}).Find(u)
	//SELECT * FROM `user` WHERE `user`.`c1` LIKE BINARY 1
}



//gen 自定义clause 子句
type StringInt int64

func (s StringInt) Value() (driver.Value, error) {
	return cast.ToString(int64(s)), nil
}
// gen自定义子句，gen自定义一些灵活的子句的时候不好实现，可以通过子查询的方式来覆盖条件  where 支持下面几种类型	case *condContainer, field.Expr, SubQuery:
type StringBinaryLike struct {
	gen.DO
	TableName string
	Column    string
	Value     string
}

func (s StringBinaryLike) BeCond() interface{} {
	return s
}

func (s StringBinaryLike) Build(builder clause.Builder) {
	builder.WriteQuoted(s.TableName)
	builder.WriteByte(46)
	builder.WriteQuoted(s.Column)
	builder.WriteString(" LIKE binary ")
	builder.AddVar(builder, s.Value)
}


func TestGenCondition1(t *testing.T) {
	dao := query.Use(db)
	u := dao.User
	s := StringInt(1)
	result, err := u.WithContext(context.Background()).Where(field.NewField(u.TableName(), u.Age.ColumnName().String()).Like(s)).First()
	if err != nil {
		log.Println(err)
	}
	//SELECT * FROM `user` WHERE `user`.`age` LIKE '1' ORDER BY `user`.`id` LIMIT 1
	log.Printf("%+v", result)

}
func TestGenCondition2(t *testing.T) {
	dao := query.Use(db)
	u := dao.User
	var s gen.SubQuery = &StringBinaryLike{
		TableName: "table",
		Column:    "column",
		Value:     "v1",
	}
	_, _ = u.WithContext(context.Background()).Where(s).Where(u.ID.Like(1)).First()
	//SELECT * FROM `user` WHERE `table`.`column` LIKE binary 'v1' AND `user`.`id` LIKE 1 ORDER BY `user`.`id` LIMIT 1
}

