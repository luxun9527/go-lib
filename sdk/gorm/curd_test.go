package main

import (
	"database/sql/driver"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm/model"
	"log"
	"os"
	"testing"
	"time"
)

var db *gorm.DB

func TestCreate(t *testing.T) {
	InitGorm()
	user := &model.User{Username: "zhangsan", Fav: "篮球"}
	//if err := db.Create(user).Error; err != nil {
	//	log.Println("create err", err)
	//}
	//INSERT INTO `user` (`username`,`fav`) VALUES ('zhangsan','')
	//	INSERT INTO `user` (`fav`) VALUES ('篮球')
	if err := db.Select("fav").Create(user).Error; err != nil {
		log.Println("create err", err)
	}
	//INSERT INTO `user` (`username`) VALUES ('zhangsan')
	if err := db.Omit("fav").Create(user).Error; err != nil {
		log.Println("create err", err)
	}

}
func TestUpdate(t *testing.T) {
	InitGorm()
	u := &model.User{ID: 1}
	//更新单个列 只有非空的字段会被修改UPDATE `user` SET `username`='hello' WHERE `id` = 1
	db.Model(&u).Update("username", "hello")
	//更新多个列,只有非空的字段会被修改 UPDATE `user` SET `username`='hello12' WHERE `id` = 1
	db.Model(&u).Updates(model.User{Username: "hello12"})
	//select 更新指定了列  UPDATE `user` SET `username`='hello12',`fav`='finish' WHERE `id` = 1
	db.Model(&u).Select("username").Updates(&model.User{
		Username: "",
		Fav:      "",
	})
	//忽略指定的列，会忽略空值。
	//UPDATE `user` SET `fav`='1' WHERE `id` = 1
	db.Model(&u).Omit("username").Updates(&model.User{
		Username: "",
		Fav:      "1",
	})
	//指定条件更新，会忽略空值。 UPDATE `user` SET `username`='1' WHERE username = 'admin'
	db.Model(model.User{}).Where("username = ?", "admin").Updates(model.User{Username: "1"})
	//默认没有指定条件不会全局更新。db.Model(&User{}).Update("name", "jinzhu").Error gorm.ErrMissingWhereClause

	//如果要更新零字段，要使用map或者select指定字段。
	// Select all fields (select all fields include zero value fields)
	//db.Model(&user).Select("*").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})

	// Select all fields but omit Role (select all fields include zero value fields)
	//db.Model(&user).Select("*").Omit("Role").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})
}

func find() {
	//var u User1
	//if err := db.Where(User1Columns.ID+" = ?", 1).Find(&u).Error; err != nil {
	//	log.Println("err", err)
	//}
	//log.Println(u)
}
func insert() {
	//refer https://gorm.io/zh_CN/docs/create.html

}
func InitGorm() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color

		},
	)

	dsn := "root:123456@tcp(192.168.2.99:3306)/gormtest?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlBb, err := db.DB()
	sqlBb.SetMaxIdleConns(100)
	sqlBb.SetMaxOpenConns(1000)
	db = db
}

type CustomTime time.Time

func (t CustomTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006/01/02 15:04:05"), nil
}
func (t CustomTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	tStr := tTime.Format("2006-01-02 15:04:05") // 设置格式
	// 注意 json 字符串风格要求
	tStr = "\"" + tStr + "\""
	return []byte(tStr), nil
}

func (t *CustomTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case string:
		// 字符串转成 time.Time 类型
		tTime, _ := time.Parse("2006/01/02 15:04:05", vt)
		*t = CustomTime(tTime)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

type Custom struct {
	ID         uint64     `gorm:"column:id;" json:"id"`
	CustomTime CustomTime `gorm:"column:custom_time;" json:"custom_time"`
}

func CustomType() {
	c := &Custom{
		CustomTime: CustomTime(time.Now()),
	}
	db.Table("test_custom").Create(c)
	var c2 Custom
	db.Table("test_custom").Where("id =1 ").Find(&c2)
	log.Printf("%+v", c2)
}
func hasOne() {
	var u []*model.User
	if err := db.Model(model.User{}).Preload("UserProfile").Find(&u).Error; err != nil {
		log.Println("err", err)
		return
	}
	for _, v := range u {
		log.Println(*v)
	}
}
func manyToMany() {

}
func hasMany() {
	var u []*model.User
	if err := db.Model(model.User{}).Preload("UserProfile").Find(&u).Error; err != nil {
		log.Println("err", err)
		return
	}
	for _, v := range u {
		log.Println(*v)
	}
}
