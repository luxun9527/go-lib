package standard

import "testing"

func TestCompare(t *testing.T) {
	type P struct {
		Name string
	}
	p1 := P{Name: "sasdf"}
	p2 := P{Name: "sassdf"}
	if p1 == p2 {
		t.Log(true)
	} else {

	}

}
