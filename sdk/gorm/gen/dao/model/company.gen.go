// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameCompany = "company"

// Company mapped from table <company>
type Company struct {
	ID   int32  `gorm:"column:id;primaryKey" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
}

// TableName Company's table name
func (*Company) TableName() string {
	return TableNameCompany
}