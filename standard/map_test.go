package standard

import (
	"log"
	"sync"
	"testing"
)

func TestMapAssign(t *testing.T) {
	type Person struct {
		Name string
	}
	m := map[string]Person{"zhangsan": {Name: ""}}
	p1 := m["str"]
	log.Printf("%p\n", &p1)
	log.Printf("%p\n", &m)
	m2 := map[string]interface{}{"zhangsna": 1}
	delete(m2, "test")
}
func TestSyncMap(t *testing.T) {
	var m sync.Map
	m.Store(1, 2)
	m.Store(2, 3)
	m.Range(func(key, value any) bool {
		m.Delete(key)
		return true
	})

	m.Range(func(key, value any) bool {
		log.Println(key, value)
		return true
	})
}
func TestRangeDelete(t *testing.T) {
	m := map[int32]interface{}{1: 2, 3: 4}
	for k, _ := range m {
		delete(m, k)
	}
	for k, v := range m {
		log.Println(k, v)
	}
}
func TestSyncMap1(t *testing.T) {
	var m sync.Map
	m.Store(1, 2)
	m.Store(2, 3)
	for  {
		store, loaded := m.LoadOrStore(3, 4)
		log.Println(store)
		if loaded{
			break
		}
	}
}

func TestName(t *testing.T) {
	var s = map[int32]string{1:"1"}
	s2 := s[2]
	log.Println(s2)
}