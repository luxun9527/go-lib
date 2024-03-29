package main

import "log"

// 读取者接口
type Reader interface {
	Read([]byte) error
}

var (
	//demo全局变量
	_globalVar int32 = 1
)

const TableNameCard = "card"

// Card mapped from table <card>
type Card struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	No        int32  `gorm:"column:no;not null;Comment:卡号" json:"no"`                   // 卡号
	UserID    int32  `gorm:"column:user_id;not null;Comment:用户id" json:"user_id"`       // 用户id
	Amount    string `gorm:"column:amount;not null;Comment:金额" json:"amount"`           // 金额
	CreatedAt int64  `gorm:"column:created_at;not null;Comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;not null;Comment:修改时间" json:"updated_at"` // 修改时间
	DeletedAt int64  `gorm:"column:deleted_at;not null;Comment:删除时间" json:"deleted_at"` // 删除时间
}

// TableName Card's table name
func (*Card) TableName() string {
	return TableNameCard
}

/*
这是一个注释
*/
func DemoSwitch() {
	switch {
	case true:
		log.Println("true")
	}
}
