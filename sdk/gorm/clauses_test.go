package main

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/spf13/cast"
	"go-lib/sdk/gorm/gen/dao/query"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
	"log"
	"testing"
	"time"
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
	u.WithContext(context.Background()).Select(u.ID).First()
}


type JsonFunc struct {
	field.Expr
}
func (j JsonFunc)BuildWithArgs(){

}

const TableNameEntrustOrder = "entrust_order"

// EntrustOrder mapped from table <entrust_order>
type EntrustOrder struct {
	ID             int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	OrderID        string    `gorm:"column:order_id;comment:uuid" json:"order_id"`
	UserID         int32     `gorm:"column:user_id;not null;comment:用户id" json:"user_id"`
	SymbolID       int32     `gorm:"column:symbol_id;not null;comment:交易对ID" json:"symbol_id"`
	SymbolName     string    `gorm:"column:symbol_name;not null;comment:交易对名称" json:"symbol_name"`
	Qty            string    `gorm:"column:qty;not null;comment:下单数量" json:"qty"`
	Price          string    `gorm:"column:price;not null;comment:价格" json:"price"`
	Side           int32     `gorm:"column:side;not null;comment:方向1买 2卖" json:"side"`
	Amount         string    `gorm:"column:amount;not null;comment:金额" json:"amount"`
	Status         int32     `gorm:"column:status;not null;comment:状态1新订单2部分成交 3全部成交，4撤销，5无效订单" json:"status"`
	OrderType      int32     `gorm:"column:order_type;not null;comment:订单类型1市价单2限价单" json:"order_type"`
	FilledQty      string    `gorm:"column:filled_qty;not null;comment:成交数量" json:"filled_qty"`
	UnFilledQty    string    `gorm:"column:un_filled_qty;not null;comment:未成交数量" json:"un_filled_qty"`
	FilledAvgPrice string    `gorm:"column:filled_avg_price;not null;comment:成交均价" json:"filled_avg_price"`
	FilledAmount   string    `gorm:"column:filled_amount;not null;comment:成交金额" json:"filled_amount"`
	UnFilledAmount string    `gorm:"column:un_filled_amount;not null;comment:未成交金额" json:"un_filled_amount"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"`
}

// TableName EntrustOrder's table name
func (*EntrustOrder) TableName() string {
	return TableNameEntrustOrder
}

const TableNameMatchedOrder = "matched_order"

// MatchedOrder mapped from table <matched_order>
type MatchedOrder struct {
	ID           int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MatchID      string    `gorm:"column:match_id;not null;comment:撮合id" json:"match_id"`
	SymbolID     int32     `gorm:"column:symbol_id;not null;comment:交易对id" json:"symbol_id"`
	SymbolName   string    `gorm:"column:symbol_name;not null;comment:交易对名称" json:"symbol_name"`
	TakerOrderID string    `gorm:"column:taker_order_id;not null;comment:taker订单id" json:"taker_order_id"`
	MakerOrderID string    `gorm:"column:maker_order_id;not null;comment:maker订单id" json:"maker_order_id"`
	Price        string    `gorm:"column:price;not null;comment:价格" json:"price"`
	Qty          string    `gorm:"column:qty;not null;comment:数量(基础币)" json:"qty"`
	Amount       string    `gorm:"column:amount;not null;comment:金额（计价币）" json:"amount"`
	MatchTime    int64     `gorm:"column:match_time;comment:撮合时间" json:"match_time"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"`
}

// TableName MatchedOrder's table name
func (*MatchedOrder) TableName() string {
	return TableNameMatchedOrder
}

const TableNameKline = "kline"

// Kline mapped from table <kline>
type Kline struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	StartTime int64  `gorm:"column:start_time;not null;comment:k线开始时间" json:"start_time"`
	EndTime   int64  `gorm:"column:end_time;not null;comment:k线结束时间" json:"end_time"`
	Symbol    string `gorm:"column:symbol;not null;comment:交易对" json:"symbol"`
	SymbolID  int32  `gorm:"column:symbol_id;comment:交易对id" json:"symbol_id"`
	KlineType int32  `gorm:"column:kline_type;not null;comment:k线类型1分钟 5分钟" json:"kline_type"`
	Open      string `gorm:"column:open;not null;comment:开盘价" json:"open"`
	High      string `gorm:"column:high;not null;comment:k线内最高价" json:"high"`
	Low       string `gorm:"column:low;not null;comment:k线内最低价" json:"low"`
	Close     string `gorm:"column:close;not null;comment:收盘价" json:"close"`
	Volume    string `gorm:"column:volume;not null;comment:成交量(基础币数量)" json:"volume"`
	Turnover  string `gorm:"column:turnover;not null;comment:成交额(计价币数量)" json:"turnover"`
	Range     string `gorm:"column:range;not null;comment:涨跌幅" json:"range"`
}

// TableName Kline's table name
func (*Kline) TableName() string {
	return TableNameKline
}

const TableNameUser1 = "user"

// User mapped from table <user>
type User1 struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username    string    `gorm:"column:username;not null;comment:用户名" json:"username"`
	Password    string    `gorm:"column:password;not null;comment:密码" json:"password"`
	PhoneNumber int64     `gorm:"column:phone_number;not null;comment:手机号" json:"phone_number"`
	Status      int32     `gorm:"column:status;not null;comment:用户状态，1正常2锁定" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;comment:更新时间" json:"updated_at"`
}

// TableName User's table name
func (*User1) TableName() string {
	return TableNameUser
}

const TableNameAsset = "asset"

// Asset mapped from table <asset>
type Asset struct {
	ID           int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID       int32     `gorm:"column:user_id;not null;comment:用户ID" json:"user_id"`
	Username     string    `gorm:"column:username;comment:用户名" json:"username"`
	CoinID       int32     `gorm:"column:coin_id;not null;comment:数字货币ID" json:"coin_id"`
	CoinName     string    `gorm:"column:coin_name;not null;comment:数字货币名称" json:"coin_name"`
	AvailableQty string    `gorm:"column:available_qty;not null;comment:可用余额" json:"available_qty"`
	FrozenQty    string    `gorm:"column:frozen_qty;not null;comment:冻结金额" json:"frozen_qty"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"`
}

// TableName Asset's table name
func (*Asset) TableName() string {
	return TableNameAsset
}
func TestM(t *testing.T) {
	fmt.Print(1)
}

