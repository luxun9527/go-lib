// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameT1 = "t1"

// T1 mapped from table <t1>
type T1 struct {
	ID   int32  `gorm:"column:id;primaryKey" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Age  int32  `gorm:"column:age;not null" json:"age"`
}

// TableName T1's table name
func (*T1) TableName() string {
	return TableNameT1
}