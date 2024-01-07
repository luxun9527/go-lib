package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"
	"os"
	"testing"
	"time"
)

const TableNameCard = "card"

// Card mapped from table <card>
type Card struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	No        int32  `gorm:"column:no;not null;comment:卡号" json:"no"`                   // 卡号
	UserID    int32  `gorm:"column:user_id;not null;comment:用户id" json:"user_id"`       // 用户id
	Amount    string `gorm:"column:amount;not null;comment:金额" json:"amount"`           // 金额
	CreatedAt int64  `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"` // 修改时间
	DeletedAt int64  `gorm:"column:deleted_at;not null;comment:删除时间" json:"deleted_at"` // 删除时间
}

// TableName Card's table name
func (*Card) TableName() string {
	return TableNameCard
}

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID        int32   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username  string  `gorm:"column:username;not null;comment:用户名" json:"username"`      // 用户名
	Age       int32   `gorm:"column:age;not null;comment:年龄" json:"age"`                 // 年龄
	Fav       string  `gorm:"column:fav;not null;comment:爱好" json:"fav"`                 // 爱好
	CreatedAt int64   `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt int64   `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"` // 修改时间
	DeletedAt int64   `gorm:"column:deleted_at;not null;comment:删除时间" json:"deleted_at"` // 删除时间
	Cards     []*Card `gorm:"foreignKey:UserID;references:ID"`
	Profile   Profile `gorm:"foreignKey:UserID;references:ID"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

const TableNameProfile = "profile"

// Profile mapped from table <profile>
type Profile struct {
	ID        int32  `gorm:"column:id;primaryKey" json:"id"`
	UserID    int32  `gorm:"column:user_id;not null;comment:用户id" json:"user_id"`       // 用户id
	Nickname  string `gorm:"column:nickname;not null;comment:昵称" json:"nickname"`       // 昵称
	Desc      string `gorm:"column:desc;not null;comment:描述" json:"desc"`               // 描述
	CreatedAt int64  `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;not null;comment:更新时间" json:"updated_at"` // 更新时间
	DeletedAt int64  `gorm:"column:deleted_at;not null;comment:删除时间" json:"deleted_at"` // 删除时间
}

// TableName Profile's table name
func (*Profile) TableName() string {
	return TableNameProfile
}

var db *gorm.DB

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	dsn := "root:root@tcp(192.168.11.185:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlBb, err := db.DB()
	sqlBb.SetMaxOpenConns(100)
	sqlBb.SetMaxIdleConns(10)
	sqlBb.SetConnMaxIdleTime(time.Hour * 5)
	//db.AutoMigrate(&User{})
	//db.AutoMigrate(&Profile{})
	//db.AutoMigrate(&Card{})
	//db.AutoMigrate(&CustomField{})

}
func TestCreate(t *testing.T) {
	user := &User{
		Username:  "",
		Age:       0,
		Fav:       "",
		CreatedAt: 0,
		UpdatedAt: 0,
	}
	//INSERT INTO `user` (`username`,`age`,`fav`,`created_at`,`updated_at`) VALUES ('',0,'',1692947238,1692947238)
	if err := db.Create(user).Error; err != nil {
		fmt.Println("create err", err)
	}
	fmt.Printf("user = %+v", user)
	//INSERT INTO `user` (`fav`,`created_at`,`updated_at`) VALUES ('',1692947405,1692947405)
	if err := db.Select("fav").Create(user).Error; err != nil {
		fmt.Println("create err", err)
	}
	//INSERT INTO `user` (`username`,`age`,`created_at`,`updated_at`,`id`) VALUES ('',0,1692947813,1692947813,14)
	if err := db.Omit("fav").Create(user).Error; err != nil {
		fmt.Println("create err", err)
	}
	//https://gorm.io/zh_CN/docs/models.html
	//官方的文档很全主要是验证一些模糊的地方
	//1、零值会被插入。
	//2、created_at updated_at会被填充当前时间插入。
	//3、插入后会将主键赋值回来。
}
func TestUpdate(t *testing.T) {
	u := &User{ID: 1}
	//更新单个列 只有非空的字段会被修改
	//UPDATE `user` SET `username`='hello',`updated_at`=1692948530 WHERE `id` = 1
	db.Model(&u).Update("username", "hello")
	//更新多个列,只有非空的字段会被修改
	// UPDATE `user` SET `username`='hello12',`updated_at`=1692948530 WHERE `id` = 1
	db.Model(&u).Updates(&User{Username: "hello12"})
	//select 更新指定了列,不会忽略零值
	// UPDATE `user` SET `username`='',`updated_at`=1692948530 WHERE `id` = 1
	db.Model(&u).Select("username").Updates(&User{
		Username: "",
		Fav:      "",
	})
	//忽略指定的列，会忽略空值。
	//UPDATE `user` SET `updated_at`=1695796004 WHERE `id` = 1
	db.Model(&u).Omit("fav").Updates(&User{
		Username: "",
		Fav:      "1",
	})
	//UPDATE `user` SET `id`=0,`username`='',`age`=0,`created_at`=0,`updated_at`=1695796443,`deleted_at`=0 WHERE `id` = 1
	//忽略指定字段其他还会被更新
	db.Select("*").Omit("fav").Updates(&User{
		ID:       1,
		Username: "",
		Fav:      "1",
	})
	//指定条件更新，会忽略空值
	// UPDATE `user` SET `id`=1,`username`='1',`updated_at`=1692948928 WHERE username = 'admin' AND `id` = 1
	db.Where("username = ?", "admin").Updates(User{Username: "1", ID: 1})
	//默认没有指定条件不会全局更新。db.Model(&User{}).Update("name", "jinzhu").Error gorm.ErrMissingWhereClause

	//如果要更新零字段，要使用map或者select指定字段。
	// UPDATE `user` SET `id`=1,`username`='',`age`=0,`fav`='',`created_at`=0,`updated_at`=1692948946 WHERE `id` = 1
	db.Select("*").Updates(User{
		ID:        1,
		Username:  "",
		Age:       0,
		Fav:       "",
		CreatedAt: 0,
		UpdatedAt: 0,
	})
}

func TestDelete(t *testing.T) {
	db.Where("id = ?", 1).Delete(&User{})
}

func TestSelect(t *testing.T) {
	//普通的查询，官方文档上写的比较全,主要是显示预加载。
	//SELECT * FROM `user` WHERE id = 6
	var users1 []*User
	if err := db.Model(&User{}).Where("id = ?", 6).Find(&users1).Error; err != nil {
		log.Println(err)
		return
	}
	//SELECT * FROM `card` WHERE `card`.`user_id` = 6
	//SELECT * FROM `user` WHERE id = 6
	//SELECT * FROM `profile` WHERE `profile`.`user_id` = 6
	var users2 []*User
	if err := db.Model(&User{}).Where("id = ?", 6).Preload("Cards").Preload("Profile").Find(&users2).Error; err != nil {
		log.Println(err)
		return
	}
	for _, v := range users2 {
		log.Printf("%+v", v)
	}
}

// 存到数据库是时间戳，取出来是字符串
const TableNameCustomField = "custom_field"

// CustomField mapped from table <custom_field>
type CustomField struct {
	ID          int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CreatedTime CustomTime `gorm:"column:created_time;not null" json:"created_time"`
}

// TableName CustomField's table name
func (*CustomField) TableName() string {
	return TableNameCustomField
}

// 演示自定义类型，实现scan和value接口。
func TestCustomField(t *testing.T) {
	cf := &CustomField{
		CreatedTime: CustomTime(cast.ToString(time.Now().Unix())),
	}
	db.Create(cf)
	var c []CustomField
	if err := db.Find(&c).Error; err != nil {
		log.Fatal(err)
	}
	log.Printf("data = %+v", c)
	//output  [{ID:1 CreatedTime:2023/08/25 16:46:18} {ID:2 CreatedTime:2023/08/25 16:46:42}]
}

type CustomTime string

func (t CustomTime) Value() (driver.Value, error) {
	return cast.ToInt64(string(t)), nil
}

// 数据库 反序列化到值
func (t *CustomTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case []byte:
		// 字符串转成 time.Time 类型
		s := time.Unix(cast.ToInt64(string(vt)), 0).Format("2006/01/02 15:04:05")
		*t = CustomTime(s)

	default:
		return errors.New("类型处理错误")
	}
	return nil
}
func (loc CustomTime) GormDataType() string {
	return "int64"
}

const TableNameUserJSON = "user_json"

// UserJSON mapped from table <user_json>
type UserJSON struct {
	ID         int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserConfig []byte `gorm:"column:user_config;not null" json:"user_config"`
}

// TableName UserJSON's table name
func (*UserJSON) TableName() string {
	return TableNameUserJSON
}
func TestJsonMap(t *testing.T) {
	var u UserJSON
	db.Where("id=?", 1).Select("user_config").Find(&u)
	log.Println(string(u.UserConfig))
	db.Create(&u)
}
