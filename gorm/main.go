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
	"time"
)

func main() {
	InitGorm()
	//insert()
	//update()
	//find()
	//CustomType()
	hasOne()
}

var DB *gorm.DB

func update() {
	//https://gorm.io/zh_CN/docs/update.html#%E6%9B%B4%E6%96%B0%E9%80%89%E5%AE%9A%E5%AD%97%E6%AE%B5

	if err := DB.Model(&User1{ID: 1}).Updates(&User1{
		ID:   0,
		Name: "",
		Age:  12,
	}).Error; err != nil {
		log.Println(err)
	}
	// UPDATE `user1` SET `age`=12,`updated_at`='2022-11-10 22:54:54.793' WHERE `id` = 1 字段名称为UpdatedAt的会自动更新。字段的默认值不更新

}

func find() {
	var u User1
	if err := DB.Where(User1Columns.ID+" = ?", 1).Find(&u).Error; err != nil {
		log.Println("err", err)
	}
	log.Println(u)
}
func insert() {
	//refer https://gorm.io/zh_CN/docs/create.html
	user := &User1{
		Name: "zhangsan",
		Age:  12,
	}
	if err := DB.Create(user).Error; err != nil {
		log.Println("insert err", err)
		return
	}

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlBb, err := db.DB()
	sqlBb.SetMaxIdleConns(100)
	sqlBb.SetMaxOpenConns(1000)
	DB = db
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
	DB.Table("test_custom").Create(c)
	var c2 Custom
	DB.Table("test_custom").Where("id =1 ").Find(&c2)
	log.Printf("%+v", c2)
}
func hasOne() {
	var u []*model.User
	if err := DB.Model(model.User{}).Preload("UserProfile").Find(&u).Error; err != nil {
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
	if err := DB.Model(model.User{}).Preload("UserProfile").Find(&u).Error; err != nil {
		log.Println("err", err)
		return
	}
	for _, v := range u {
		log.Println(*v)
	}
}
