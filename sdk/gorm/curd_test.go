package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/soft_delete"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

const TableNameCard = "card"

// Card mapped from table <card>
type Card struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	No        int32  `gorm:"column:no;not null;comment:卡号" json:"no"`
	UserID    int32  `gorm:"column:user_id;not null;comment:用户id" json:"user_id"`
	Amount    string `gorm:"column:amount;not null;default:0.000000000000000000;comment:金额" json:"amount"`
	CreatedAt int64  `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"`
	DeletedAt int64  `gorm:"column:deleted_at;not null;comment:删除时间" json:"deleted_at"`
}

// TableName Card's table name
func (*Card) TableName() string {
	return TableNameCard
}

const TableNameProfile = "profile"

// Profile mapped from table <profile>
type Profile struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID    int32  `gorm:"column:user_id;not null;comment:用户id" json:"user_id"`
	Nickname  string `gorm:"column:nickname;not null;comment:昵称" json:"nickname"`
	Desc      string `gorm:"column:desc;not null;comment:描述" json:"desc"`
	CreatedAt int64  `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at;not null;comment:更新时间" json:"updated_at"`
	DeletedAt int64  `gorm:"column:deleted_at;not null;comment:删除时间" json:"deleted_at"`
}

// TableName Profile's table name
func (*Profile) TableName() string {
	return TableNameProfile
}

// User mapped from table <user>
type User struct {
	ID        int32                 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username  string                `gorm:"column:username;not null;comment:用户名" json:"username"`
	Age       uint32                `gorm:"column:age;not null;comment:年龄" json:"age"`
	Fav       string                `gorm:"column:fav;not null;comment:爱好" json:"fav"`
	CompanyID int32                 `gorm:"column:company_id;not null;comment:公司Id" json:"company_id"`
	CreatedAt uint64                `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
	UpdatedAt uint64                `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"`
	Cards     []*Card               `gorm:"foreignKey:UserID;references:ID"`
	Profile   Profile               `gorm:"foreignKey:UserID;references:ID"`
	Company   Company               `gorm:"references:ID"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;not null" json:"deleted_at"`
}

const TableNameUser = "user"

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

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

// TableName User's table name

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

	dsn := "root:root@tcp(192.168.2.159:3308)/demo?charset=utf8mb4&parseTime=True&loc=Local"
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
	//db.AutoMigrate(&Card{})
	//db.AutoMigrate(&Profile{})

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
		log.Panicf("create no select  err %v", err)
	}
	log.Printf("user1 = %+v", user)
	//INSERT INTO `user` (`fav`,`created_at`,`updated_at`) VALUES ('',1692947405,1692947405)
	if err := db.Select("fav").Create(user).Error; err != nil {
		log.Panicf("create select field err %v", err)
	}
	log.Printf("user2 = %+v", user)
	//INSERT INTO `user` (`username`,`age`,`created_at`,`updated_at`,`id`) VALUES ('',0,1692947813,1692947813,14)
	if err := db.Omit("fav").Create(user).Error; err != nil {
		log.Panicf("create Omit field err %v", err)
	}
	log.Printf("user2 = %+v", user)

	//https://gorm.io/zh_CN/docs/models.html https://gorm.io/zh_CN/docs/create.html#%E9%BB%98%E8%AE%A4%E5%80%BC
	//官方的文档很全主要是验证一些模糊的地方
	//1、不使用select指定，零值字段也会被插入。使用select只会插入指定的字段。可以使用tag指定默认值 `gorm:"default:18"`
	//2、created_at updated_at会被填充当前时间插入。
	//3、插入后会将主键赋值回来，使用select的方式插入不会赋值回来。
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
	//select 所有再忽略指定的列，不会忽略零值
	db.Select("*").Omit("fav").Updates(&User{
		ID:       1,
		Username: "",
		Fav:      "1",
	})
	//指定条件更新，会忽略空值
	// UPDATE `user` SET `id`=1,`username`='1',`updated_at`=1692948928 WHERE username = 'admin' AND `id` = 1
	db.Where("username = ?", "admin").Updates(User{Username: "1", ID: 1})
	//默认没有指定条件不会全局更新。db.Model(&User{}).Update("name", "jinzhu").Error gorm.ErrMissingWhereClause

	//如果要更新零字段，要使用map或者select指定字段。 select 所有不会忽略零值
	// UPDATE `user` SET `id`=1,`username`='',`age`=0,`fav`='',`created_at`=0,`updated_at`=1692948946 WHERE `id` = 1
	db.Select("*").Updates(User{
		ID:        1,
		Username:  "",
		Age:       0,
		Fav:       "篮球",
		CreatedAt: 0,
		UpdatedAt: 0,
	})
}

func TestDelete(t *testing.T) {
	//UPDATE `user` SET `deleted_at`=1713714861 WHERE id = 1 AND `user`.`deleted_at` = 0
	db.Where("id = ?", 3).Delete(&User{})

	//DELETE FROM `user` WHERE id = 3
	db.Unscoped().Where("id = ?", 3).Delete(&User{})
}

func TestSelect2(t *testing.T) {
	var user []*User

	queryDB := db.Where("name = ?", "jinzhu")

	// First query
	queryDB.Where("age > ?", 10).First(&user)
	//SELECT * FROM `user` WHERE name = 'jinzhu' AND age > 10 AND `user`.`deleted_at` = 0 ORDER BY `user`.`id` LIMIT 1

	// Second query with unintended compounded condition
	queryDB.Where("age > ?", 20).First(&user)
}
func TestSelect(t *testing.T) {
	//普通的查询，官方文档上写的比较全,主要是显示预加载。
	//========================================一对一 一个用户有一个profile======================
	var users1 []*User

	if err := db.Model(&User{}).Where("id = ?", 6).Preload("Profile").Find(&users1).Error; err != nil {
		log.Println(err)
	}
	//SELECT * FROM `card` WHERE `card`.`user_id` = 6
	//SELECT * FROM `user` WHERE id = 6
	//SELECT * FROM `profile` WHERE `profile`.`user_id` = 6
	//===================================================一对多，一个用户有多个cards=================================
	var users2 []*User
	if err := db.Model(&User{}).Where("id = ?", 6).Preload("Cards").Preload("Profile").Find(&users2).Error; err != nil {
		log.Println(err)

	}
	for _, v := range users2 {
		log.Printf("%+v", v)
	}
	/*
		[0.529ms] [rows:0] SELECT * FROM `card` WHERE `card`.`user_id` = 6

		2024/04/23 22:57:55 E:/demoproject/go-lib/sdk/gorm/curd_test.go:215
		[1.098ms] [rows:0] SELECT * FROM `profile` WHERE `profile`.`user_id` = 6

		2024/04/23 22:57:55 E:/demoproject/go-lib/sdk/gorm/curd_test.go:215
		[3.268ms] [rows:1] SELECT * FROM `user` WHERE id = 6 AND `user`.`deleted_at` = 0
		2024/04/23 22:57:55 &{ID:6 Username: Age:0 Fav: CompanyID:1 CreatedAt:1713713017 UpdatedAt:1713713017 Cards:[] Profile:{ID:0 UserID:0 Nickname: Desc: CreatedAt:0 UpdatedAt:0 DeletedAt:0} Company:{ID:0 Name:} DeletedAt:0}

	*/
	//===============================belong to 一个用户属于一家公司=====================
	var (
		user    User
		company Company
	)
	//查用户所属的公司。
	//SELECT * FROM `user` WHERE id = 1 AND `user`.`deleted_at` = 0 LIMIT 1
	db.Where("id = ?", 6).Take(&user)
	//SELECT * FROM `company` WHERE `company`.`id` = 1
	if err := db.Debug().Model(&user).Association("Company").Find(&company); err != nil {
		log.Println("belong to error", err)
	}
	user.Company = company
	log.Printf("user =%+v", user)
	//user ={ID:6 Username: Age:0 Fav: CompanyID:0 CreatedAt:1713713017 UpdatedAt:1713713017  Company:{ID:1 Name:test} DeletedAt:0}

	//=====================================join 预加载==========================================================
	printLine("join 预加载")
	var u User
	if err := db.Joins("Company").Take(&u, 6).Error; err != nil {
		log.Panicf("join error %v", err)
	}
	//SELECT `user`.`id`,`user`.`username`,`user`.`age`,`user`.`fav`,`user`.`company_id`,`user`.`created_at`,`user`.`updated_at`,`user`.`deleted_at`,
	//`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `user`
	//LEFT JOIN `company` `Company` ON `user`.`company_id` = `Company`.`id` WHERE `user`.`id` = 4 AND `user`.`deleted_at` = 0 LIMIT 1
	log.Printf("join u =%+v", u)
}

// 存到数据库是时间戳，取出来是字符串
const TableNameCustomField = "custom_field"

func printLine(content ...string) {
	builder := strings.Builder{}
	str := strings.Repeat("=", 50)
	builder.WriteString(str)
	for _, v := range content {
		builder.Write([]byte(v))
	}
	builder.WriteString(str)
	str = strings.Replace(builder.String(), " ", "", -1)
	log.Println(str)
}
