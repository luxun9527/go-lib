package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID          int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username    string `gorm:"column:username;not null;comment:用户名" json:"username"`
	Password    string `gorm:"column:password;not null;comment:密码" json:"password"`
	PhoneNumber int64  `gorm:"column:phone_number;not null;comment:手机号" json:"phone_number"`
	Status      int32  `gorm:"column:status;not null;comment:用户状态，1正常2锁定" json:"status"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

var (
	c = flag.String("c", "E:\\demoproject\\go-lib\\utils\\docker\\config.yaml", "")
)

func main() {
	flag.Parse()
	r := gin.Default()
	config := initConfig()
	db := config.GormConf.MustNewGormClient()
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Panicf("auto migrate failed err %v", err)
	}
	err = db.Create(&User{
		ID:          0,
		Username:    "zhangsan",
		Password:    "121212",
		PhoneNumber: 012,
		Status:      120,
	}).Error
	if err != nil {
		log.Panicf(" create failed err %v", err)
	}
	r.GET("/get_user_list", func(c *gin.Context) {
		users := make([]*User, 0, 10)
		if db.Find(&users).Error != nil {
			c.JSON(200, gin.H{"err": err})
			return
		}
		c.JSON(200, users)
	})

	r.Run(":" + cast.ToString(config.Port))
}

func initConfig() Config {
	viper.SetConfigFile(*c)
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("ReadInConfig failed err =%v", err)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Panicf("init config faile err =%v", err)
	}
	return config
}

type Config struct {
	GormConf GormConf
	Port     int32
}

type GormConf struct {
	Ip           string `json:"ip"`
	Port         int32  `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DbName       string `json:"dbname"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

func (gc *GormConf) dns() string {
	return gc.Username + ":" + gc.Password + "@tcp(" + gc.Ip + ":" + cast.ToString(gc.Port) + ")/" + gc.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
}
func (gc *GormConf) MustNewGormClient() *gorm.DB {
	if db, err := gorm.Open(mysql.Open(gc.dns()), gc.gormConfig()); err != nil {
		log.Panicf("init gorm failed err =%v", err)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(gc.MaxIdleConns)
		sqlDB.SetMaxOpenConns(gc.MaxOpenConns)
		return db
	}
}

func (gc *GormConf) gormConfig() *gorm.Config {
	config := &gorm.Config{}
	config.SkipDefaultTransaction = true

	config.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	return config
}
