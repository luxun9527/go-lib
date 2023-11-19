package standard

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestRange(t *testing.T) {
	type student struct {
		name string
	}
	//range是值覆盖的方式。
	students,target :=make([]student,0,10) ,make(map[string]*student, 10)
	students = append(students, student{name: "1"},student{name: "2"})
	for _,v := range students {
		target[v.name]=&v
		log.Printf("%p",&v)
	}
	for _,v := range target {
		log.Printf("%+v",v)
	}
}

func TestR(t *testing.T) {
	fi, err := os.Open("D:\\project\\go-lib\\standard\\test.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		s := fmt.Sprintf("delete from ac_reward where inviter_kai_id = %s and reward_amount = 0.3000000000 limit 1;", string(a))
		fmt.Println(s)
	}

}