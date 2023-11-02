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
	"reflect"
	"testing"
	"unsafe"
)

// Like whether string matches regular expression
type LikeBinary struct {
	Table  string
	Column string
	Value  interface{}
}

// 实现build接口，可以自定义表达式
func (like LikeBinary) Build(builder clause.Builder) {
	builder.WriteQuoted(like.Table)
	builder.WriteByte(46)
	builder.WriteQuoted(like.Column)
	builder.WriteString(" LIKE BINARY ")
	builder.AddVar(builder, like.Value)
}

// 演示gorm子句 https://gorm.io/zh_CN/docs/data_types.html
func TestGormClause(t *testing.T) {
	var u User
	db.Where(clause.Like{Column: "c1", Value: "v1"}).Find(u)
	// SELECT * FROM `user` WHERE `c1` LIKE 'v1'
	db.Where(LikeBinary{
		Table:  "user",
		Column: "c1",
		Value:  1,
	}).Find(u)
	//SELECT * FROM `user` WHERE `user`.`c1` LIKE BINARY 1
	//使用聚合函数

}

// gen 自定义clause 子句
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
	result, err := u.WithContext(context.Background()).Select().Where(field.NewField(u.TableName(), u.Age.ColumnName().String()).Like(s)).First()
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
	u.WithContext(context.Background()).Select(u.ID.Max()).First()

}

func TestGenCondition3(t *testing.T) {
	dao := query.Use(db)
	u := dao.User

	//customField := GenCustomField{Expr: field.EmptyExpr()}

	//expr := reflect.ValueOf(field.EmptyExpr())
	//
	//addressableSourceCopy := reflect.New(expr.Type()).Elem()
	//addressableSourceCopy.Set(expr)
	//
	//e := addressableSourceCopy.FieldByName("e")
	//
	//newE := reflect.NewAt(e.Type(), unsafe.Pointer(e.UnsafeAddr())).Elem()
	//sql := newE.FieldByName("SQL")
	//
	//newSql := reflect.NewAt(sql.Type(), unsafe.Pointer(sql.UnsafeAddr())).Elem()
	//newSql.SetString("test")
	//c := clause.Expr{
	//	SQL:                "112121",
	//	Vars:               nil,
	//	WithoutParentheses: false,
	//}
	//var cb CustomColumn
	//newE.Set(reflect.ValueOf(cb))
	//cs := &GenCustomField{Expr: field.EmptyExpr()}

	//v := reflect.ValueOf(cs).Elem().FieldByName("Expr")
	//v1 := reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	//v1.Set(addressableSourceCopy)
	//log.Printf("%+v", expr)
	//expr = addressableSourceCopy
	//column := CustomColumn{Column: "JSON_EXTRACT( activity_config, '$.backup_page_id' )"}
	c := NewCustomColumn("JSON_EXTRACT( activity_config, '$.backup_page_id' )")
	_, _ = u.WithContext(context.Background()).Select(c).Where(u.ID.Like(1)).First()
	//SELECT JSON_EXTRACT( activity_config, '$.backup_page_id' ) FROM `user` WHERE `user`.`id` LIKE 1 ORDER BY `user`.`id` LIMIT 1

}


/*
只用在你想自定义列的时候
当你想 使用json函数的时候 JSON_EXTRACT( activity_config, '$.backup_page_id' ) gen 没有这个函数
*/
type CustomColumn struct {
	Column string
	field.Expr
}
func NewCustomColumn(column string)CustomColumn{
	c := CustomColumn{
		Column: column,
	}
	c.replace()
	return c
}
//由于gen没有开放替换e(clause.Expression)的操作，只能通过反射的方式将一个 expr 中的 e 替换成我们自己的。
func (c *CustomColumn)replace(){
	expr := reflect.ValueOf(field.EmptyExpr())
	//获取一个新的expr,新建的expr有CANSET和CANADDR能力。
	newExpr := reflect.New(expr.Type()).Elem()
	newExpr.Set(expr)
	// 获取e
	e := newExpr.FieldByName("e")
	newE := reflect.NewAt(e.Type(), unsafe.Pointer(e.UnsafeAddr())).Elem()
	newE.Set(reflect.ValueOf(*c))
	//修改e 为我们自定义的e

	//修改 c的 field.Expr为刚才自定义e。
	v := reflect.ValueOf(c).Elem().FieldByName("Expr")
	v1 := reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	v1.Set(newExpr)
}


func (c CustomColumn) Build(builder clause.Builder) {
	result, err := builder.WriteString(c.Column)
	if err!=nil{
		log.Printf("err = %v valule =%v",err,result)
	}
}