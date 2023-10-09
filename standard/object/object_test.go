package object

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

type myint int32

// https://cyent.github.io/golang/basic/type_custom/
func (m *myint) print() {
	fmt.Println(*m)
}

func TestCustom(t *testing.T) {
	//自定义类型，是一种新的类型,需要强制转换,且原始的类型也可以是自定义类型
	int32Type := int32(1)
	p := myint(int32Type)
	p.print()

}

/*
例中，tempString是string的别名，其本质上与string是同一个类型。类型别名只会在代码中存在，编译完成后不会有如tempString一样的类型别名。所以变量s的类型是string。
字符类型中的byte和rune就是类型别名：

type byte = uint8
type rune = int32
类型别名这个功能非常有用，鉴于go中有些类型写起来非常繁琐，比如json相关的操作中，经常用到map[string]interface {}这种类型，写起来是不是很繁琐，没关系，给它起个简单的别名!这样用起来爽多了。

type strMap2Any = map[string]interface {}*/

type tempString = string

func TestAlisa(t *testing.T) {
	//自定义类型不用强装
	var stringType string

	stringType = "foo"

	var s tempString
	s = stringType

	fmt.Println(s)        // 我是s
	fmt.Printf("%T\n", s) // string
}

type alisaMyInt = myint

func (alisaMyInt) printValue() {

}

// 这段代码编译会报错，自定义类型只能在原始类型所在的包中定义方法
// func (tempString) print(){
//
// }
func TestDefer(t *testing.T) {
	err := func1()
	if err != nil {
		fmt.Println("after func err =", err)
	}

}
func func1() (err error) {
	//defer 在return 后执行
	defer func() {
		if err != nil {
			log.Println("defer err =", err)
			err = nil
		}
	}()
	return errors.New("not found")
}

// interface object
type Mover interface {
	move()
}

type dog struct{}

// 指针接收 对象是值类型不能实现接口。
func (d *dog) move() {
	fmt.Println("狗会动")
}
func TestImp(t *testing.T) {

	//var m Mover = dog{} //不能直接使用dog类型

	//var fugui Mover = &dog{}  // 富贵是*dog类型

}







//演示内嵌继承，内嵌继承拥有被内嵌所有的属性。
type F struct {
	Embed
}

func ( F )Print() {
	fmt.Println("f")
}

type Embed struct {
	
}

func (e Embed )Print(){
	fmt.Println("embed")
}
func TestEmbed(t *testing.T) {
	embed := Embed{}
	embed.Print()
	f := F{}
	f.Print()
}



