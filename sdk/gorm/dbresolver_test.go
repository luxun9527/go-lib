package main

import (
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stringx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"testing"
	"time"
)

const TableNameRunoobTbl = "runoob_tbl"

// RunoobTbl mapped from table <runoob_tbl>
type RunoobTbl struct {
	RunoobID       int32     `gorm:"column:runoob_id;primaryKey;autoIncrement:true" json:"runoob_id"`
	RunoobTitle    string    `gorm:"column:runoob_title;not null" json:"runoob_title"`
	RunoobAuthor   string    `gorm:"column:runoob_author;not null" json:"runoob_author"`
	SubmissionDate time.Time `gorm:"column:submission_date" json:"submission_date"`
}

// TableName RunoobTbl's table name
func (*RunoobTbl) TableName() string {
	return TableNameRunoobTbl
}

const (
	masterDsn = "root:root@tcp(192.168.2.99:33307)/test1?charset=utf8mb4&parseTime=True&loc=Local"
	slaveDsn  = "root:root@tcp(192.168.2.99:33306)/test1?charset=utf8mb4&parseTime=True&loc=Local"
)

// gentool --dsn="root:root@tcp(192.168.2.99:33307)/test1?charset=utf8mb4&parseTime=True&loc=Local" --onlyModel=true --db=mysql --tables=runoob_tbl -outPath=./ -fieldMap="decimal:string;tinyint:int32;"
func TestDbresolve(t *testing.T) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color

		},
	)
	masterDB, err := gorm.Open(mysql.New(mysql.Config{
		DSN: masterDsn,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := masterDB.Use(
		dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.New(mysql.Config{
				DSN: masterDsn,
			})},
			Replicas: []gorm.Dialector{mysql.New(mysql.Config{
				DSN: slaveDsn,
			})},
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}).
			SetMaxIdleConns(10).
			SetConnMaxLifetime(time.Hour).
			SetMaxOpenConns(200),
	); err != nil {
		log.Fatal(err)
	}
	d := &RunoobTbl{
		RunoobTitle:    cast.ToString(time.Now().Unix()),
		RunoobAuthor:   stringx.Randn(10),
		SubmissionDate: time.Now(),
	}
	masterDB.Create(d)

	if masterDB.Where("runoob_id=?", 2).Find(d).Error != nil {
		log.Println(err)
	}
	log.Println(d)
}
