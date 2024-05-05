package main

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
)

type CustomDateModel struct {
	ID         int32 `gorm:"column:id;not null;" json:"user_id"`
	CustomDate Date  `gorm:"column:custom_date;not null;" json:"custom_date"`
}

// TableName Card's table name
func (*CustomDateModel) TableName() string {
	return "custom_date"
}

func TestCustomType(t *testing.T) {
	db.Create(&CustomDateModel{CustomDate: Date(time.Now())})
	//INSERT INTO `custom_date` (`custom_date`) VALUES ('2024-05-05 00:00:00')
}

type Date time.Time

// 将数据库中的数据转换为Date类型
func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

// 将Date类型转换为数据库中的数据
func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}

func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

func (date Date) MarshalJSON() ([]byte, error) {
	return time.Time(date).MarshalJSON()
}

func (date *Date) UnmarshalJSON(b []byte) error {
	return (*time.Time)(date).UnmarshalJSON(b)
}
