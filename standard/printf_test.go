package standard

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	//var f1 float64 = 0
	//fmt.Printf("%.2f\n", f1)
	fmt.Println(time.Now())
	parseTime, err := time.Parse("2006-01-02T15:04:05.999999999Z", "2020-12-21T06:33:14.000000000Z")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(parseTime)
}
