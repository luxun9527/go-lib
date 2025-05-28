package standard

import (
	"log"
	"strings"
	"testing"
)

// 验证问题
/* 切片的本质
type slice struct {
    array unsafe.Pointer // 指向底层数组的指针
    len   int            // 当前切片的长度（元素个数）
    cap   int            // 切片的容量（从起始位置到底层数组末尾的元素总数）
}
*/
// 切片传递
func TestCopy(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{5, 4, 3}
	slice3 := []int{7, 8, 9}
	copy(slice2, slice1) // 只会复制slice1的前3个元素到slice2中
	t.Log("slice2", slice2)
	copy(slice1, slice3) // 只会复制slice2的3个元素到slice1的前3个位置
	t.Log("slice1", slice1)
	s := make([]int32, 1, 10)
	copy(s, []int32{1})
	t.Log(s)
}
func TestString(t *testing.T) {
	b := []byte("1234567890000000")
	log.Println(b)
	log.Println(string(b))
}
func TestSliceCut(t *testing.T) {
	m := []byte{97, 98, 100, 0}
	log.Printf("%v= ", string(m[:4]))
	log.Println(m[1:])
	b := []byte{0}
	log.Println(string(b))
	f := []byte{0, 0, 0, 0, 97}

	f1 := strings.ReplaceAll(string(f), string([]byte{0}), "")

	log.Printf(" len(f1) =%v f1 =%v", len(f1), f1)

}
func TestSlice1(t *testing.T) {
	str := "a	b	23	s	c	d"
	s := strings.Split(str, "	")
	s = s[1 : len(s)-2]
	sr := strings.Join(s, "\t")
	log.Println(sr + "2")
}
func TestRangeSlice(t *testing.T) {
	s := []uint8{1, 2, 3, 4}

	for _, v := range s {
		s = append(s, v)
	}
	log.Println(s)
}
func TestRangeSlice1(t *testing.T) {
	type Person struct {
		Name string
	}
	s := []Person{{Name: ""}, {Name: ""}}

	for i, _ := range s {
		s[i].Name = "1"
	}
	log.Println(s)
}
func TestSliceMake(t *testing.T) {
	m := make([]byte, 2, 10)
	log.Println(m)
	log.Println(len(m), cap(m))
	m = m[:4]
	log.Println(len(m), cap(m))
	log.Println(m)
}
func TestSliceCopy(t *testing.T) {
	m := make([]byte, 4)
	m[1] = 2
	m[2] = 9
	n := copy(m, m[1:3])
	log.Println(n)
	log.Println(m)
}
func TestResize(t *testing.T) {
	arr := make([]int32, 0, 5)
	log.Printf("%p", arr)
	arr = append(arr, []int32{12, 3, 2}...)
	log.Printf("%p", arr)
	arr = append(arr, []int32{12, 3, 223, 3, 2, 3, 2}...)
	log.Printf("%p", arr)
	/*
		2023/05/23 23:04:42 0xc00000c1f8
		2023/05/23 23:04:42 0xc00000c1f8
		2023/05/23 23:04:42 0xc00013a6c0
	*/
}

func TestPass(t *testing.T) {

	userIdList := make([]int32, 0, 10)
	//本质是值传递。 len cap属性改变不会影响其他变量。
	log.Printf("userIdList 地址:%p 值:%v", userIdList, userIdList)
	userIdList1 := &userIdList
	log.Printf("userIdList1 地址:%p 值:%v", userIdList1, userIdList1)

	userIdList = append(userIdList, 1)

	log.Printf("修改后 userIdList 地址:%p 值:%v", userIdList, userIdList)
	log.Printf("修改后 userIdList1 地址:%p 值:%v", userIdList1, userIdList1)

}
func TestPass2(t *testing.T) {

	userIdList := make([]int32, 4)
	log.Printf("userIdList 地址:%p 值:%v", userIdList, userIdList)
	userIdList1 := userIdList
	log.Printf("userIdList1 地址:%p 值:%v", userIdList1, userIdList1)
	//本质是值传递。 len cap属性改变不会影响其他变量。
	userIdList = make([]int32, 0, 3)
	userIdList = append(userIdList, 1)
	log.Printf("修改后 userIdList 地址:%p 值:%v", userIdList, userIdList)
}
