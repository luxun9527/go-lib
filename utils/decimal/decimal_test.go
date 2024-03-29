package decimal

import (
	"github.com/shopspring/decimal"
	"log"
	"testing"
)

func TestDecimal(t *testing.T) {
	d, err := decimal.NewFromString("1.32000")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(d)
	t.Log(d.RoundCeil(3).String())
	t.Log(d.Round(3).String())
	t.Log(d.RoundFloor(3).String())
	d1, _ := decimal.NewFromString("1.1101")
	d2, _ := decimal.NewFromString("0.01")
	d3 := d1.Mod(d2)
	t.Log(d3)

	d4 := decimal.New(1, 18)
	log.Println(d4)
	d5:=decimal.New(2, 12)
	log.Println(d5)
	result := d5.Div(d4).StringFixedBank(18)
	log.Println(result)


}



