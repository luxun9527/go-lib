// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameDomin = "domin"

// Domin mapped from table <domin>
type Domin struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:id" json:"id"` // id
	DomainName string `gorm:"column:domain_name;not null;comment:域名" json:"domain_name"`    // 域名
	Remark     string `gorm:"column:remark;not null;comment:备注" json:"remark"`              // 备注
	UpdatedBy  int64  `gorm:"column:updated_by;not null;comment:最后修改人id" json:"updated_by"` // 最后修改人id
	CreatedAt  int64  `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`    // 创建时间
	UpdatedAt  int64  `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"`    // 修改时间
	DeletedAt  int64  `gorm:"column:deleted_at;not null;comment:删除时间" json:"deleted_at"`    // 删除时间
}

// TableName Domin's table name
func (*Domin) TableName() string {
	return TableNameDomin
}