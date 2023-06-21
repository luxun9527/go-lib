package standard

import (
	"fmt"
	"math/big"
	"testing"
)

func TestBigInt(t *testing.T) {
	n := new(big.Int)
	n, ok := n.SetString("1234567890987.923232323", 10)
	if !ok {
		fmt.Println("SetString: error")
		return
	}
	fmt.Println(n)
}
