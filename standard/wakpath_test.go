package standard

import (
	"flag"
	cb "github.com/emirpasic/gods/queues/circularbuffer"
	"io/fs"
	"log"
	"path/filepath"
	"testing"
)

func TestWalk(t *testing.T) {
	filepath.Walk("S:\\go-lib\\standard", func(path string, info fs.FileInfo, err error) error {
		println(filepath.Ext(info.Name()))
		println(info.Name())
		return nil
	})
}
func TestReadDir(t *testing.T) {

	s := flag.String("GOPATH", "", "")
	flag.Parse()
	log.Println(*s)
}

type u struct {
	Name string
}

func TestReadDir1(t *testing.T) {
	queue := cb.New(10)
	ulist := make([]u, 0, 10)
	log.Printf("%p %v", ulist, ulist)
	queue.Enqueue(ulist)
	ulist = append(ulist, u{Name: "test"})
	log.Printf("%p %v", ulist, ulist)
	log.Printf("%p", queue.Values()[0])
	log.Printf("%v", queue.Values()[0])
}
func TestReadDir2(t *testing.T) {
	queue := cb.New(5)
	for i := 0; i < 10; i++ {
		queue.Enqueue(i)
	}
	log.Printf("%v", queue.Values())
}
