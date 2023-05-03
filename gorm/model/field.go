package model

// Role [...]
type Role struct {
	ID       int    `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	RoleName string `gorm:"column:role_name;type:varchar(50);not null;default:''" json:"roleName"`
}

// TableName get sql table name.获取数据库表名
func (m *Role) TableName() string {
	return "role"
}

// RoleColumns get sql column name.获取数据库列名
var RoleColumns = struct {
	ID       string
	RoleName string
}{
	ID:       "id",
	RoleName: "role_name",
}

// User [...]
type User struct {
	ID          int         `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	Username    string      `gorm:"column:username;type:varchar(50);not null;default:''" json:"username"`
	UserProfile UserProfile `gorm:"foreignKey:UserID" json:"userProfile"`
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}

// UserColumns get sql column name.获取数据库列名
var UserColumns = struct {
	ID       string
	Username string
}{
	ID:       "id",
	Username: "username",
}

// UserProfile [...]
type UserProfile struct {
	ID       int    `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	Nickname string `gorm:"column:nickname;type:varchar(50);not null" json:"nickname"`
	UserID   int    `gorm:"column:user_id;type:int;not null" json:"userId"`
}

// TableName get sql table name.获取数据库表名
func (m *UserProfile) TableName() string {
	return "user_profile"
}

// UserProfileColumns get sql column name.获取数据库列名
var UserProfileColumns = struct {
	ID       string
	Nickname string
	Fav      string
	UserID   string
}{
	ID:       "id",
	Nickname: "nickname",
	Fav:      "fav",
	UserID:   "user_id",
}

// UserRole [...]
type UserRole struct {
	UserID int `gorm:"primaryKey;column:user_id;type:int;not null" json:"userId"`
	RoleID int `gorm:"column:role_id;type:int;not null" json:"roleId"`
}

// TableName get sql table name.获取数据库表名
func (m *UserRole) TableName() string {
	return "user_role"
}

// UserRoleColumns get sql column name.获取数据库列名
var UserRoleColumns = struct {
	UserID string
	RoleID string
}{
	UserID: "user_id",
	RoleID: "role_id",
}
