package standard

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestJsonTime(t *testing.T) {
	type TimeFormat struct {
		UpdateTime time.Time `json:"updateTime"`
	}
	//data := `{"updateTime":"2023-03-31T13:03:57+08:00"}`
	//data := `{"updateTime":"2023-03-28T02:48:17.000000000Z"}`
	data := `{"updateTime":"2023-03-02T15:52:53Z"}`
	var tf TimeFormat
	if err := json.Unmarshal([]byte(data), &tf); err != nil {
		log.Fatal(err)
	}
	log.Println(tf)

	d, err := time.Parse("2006-01-02 15:04:05", "2023-03-02 15:52:53")
	if err != nil {
		log.Println(err)
		return
	}
	year, month, day := d.Date()
	log.Printf("year = %v month = %v day = %v", year, int64(month), day)
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println(err)
		return
	}
	//zone := time.FixedZone("UTC-4 America/New_York", -4*60*60)
	log.Println(time.Now().In(time.UTC))
	log.Println(time.Now().In(location))
	log.Println(time.Unix(1688011200, 0).In(location))
	log.Println(time.Unix(1688011200, 0).In(time.Local))
	log.Println(time.Unix(1688011200, 0).In(time.UTC))
	log.Println(time.Unix(1688011200, 0).In(location).Unix())
	log.Println(time.Now().In(location).Unix())
	log.Println(time.Now().Unix())                                           //unix和时区无关，返回的都是unix时区的。
	log.Println(time.Now().Unix() / (24 * 60 * 60 * 1) * (24 * 60 * 60 * 1)) //unix和时区无关，返回的都是unix时区的。
	date := time.Now().Format("2006-01-02")
	inLocation, _ := time.ParseInLocation("2006-01-02", date, location)
	log.Println(inLocation.Unix())

	unix := time.Now().Unix()
	log.Println(unix)
	//loc, _ := time.LoadLocation("Asia/Shanghai")
	log.Println(time.Unix(unix,0).In(time.UTC).Format("2006-01-02 15:04:05"))
}
