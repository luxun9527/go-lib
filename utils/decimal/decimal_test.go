package decimal

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestDecimal(t *testing.T) {
	d, err := decimal.NewFromString("1.3223")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(d.RoundCeil(3).String())
	t.Log(d.Round(3).String())
	t.Log(d.RoundFloor(3).String())
}
