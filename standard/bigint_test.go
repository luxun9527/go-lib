package standard

import (
	"crypto/rand"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"log"
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
func TestLottery(t *testing.T) {
	p :=[]RewardGrade{{
		Grade:       1,
		Probability: "0.1",
	}, {
		Grade:       2,
		Probability: "0.2",
	}, {
		Grade:       3,
		Probability: "0.7",
	},
	}
	Lottery(p)
}

type RewardGrade struct {
	Grade int32
	Probability string

}

type Range struct {
	Grade int32
	start,end int64
}
func Lottery(p []RewardGrade){
	//概率最小为小数点后两位 则总数对应为10000 则生成随机数在[0-10000) 区间内
	//一等奖 抽中的概率为 0.01% 万份之一 即在1-10000 随机数选中0
	//二等奖 抽中的概率为 1% 万份之一 即在1-10000 随机数选中 [1-101)

	max,start := decimal.NewFromInt(10000) , decimal.NewFromInt(0)
	var rangeList []*Range
	for _,v := range p {
		probability, _ := decimal.NewFromString(v.Probability)
		end := probability.Mul(max).Add(start)
		r := &Range{
			Grade: v.Grade,
			start: cast.ToInt64(start.String()),
			end:   cast.ToInt64(end.String()),
		}
		start = end
		rangeList = append(rangeList, r)
	}
	for i := 0; i < 10; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10000))
		for _,v := range rangeList {
			if v.start <= n.Int64() && n.Int64() < v.end {
				log.Printf("rand value=%v in grade=%v grade start=%v end=%v",n.Int64(),v.Grade,v.start,v.end)
				break
			}

		}
	}


	
}