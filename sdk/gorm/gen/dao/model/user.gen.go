// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username  string `gorm:"column:username;not null;comment:用户名" json:"username"`      // 用户名
	Age       uint32 `gorm:"column:age;not null;comment:年龄" json:"age"`                 // 年龄
	Fav       string `gorm:"column:fav;not null;comment:爱好" json:"fav"`                 // 爱好
	CompanyID int32  `gorm:"column:company_id;not null;comment:公司Id" json:"company_id"` // 公司Id
	CreatedAt uint64 `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt uint64 `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"` // 修改时间
	DeletedAt uint64 `gorm:"column:deleted_at;not null" json:"deleted_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
