package decimal

import (
	"github.com/shopspring/decimal"
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

}
