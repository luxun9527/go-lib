package error

import (
	"github.com/pkg/errors"
	"log"
	"testing"
)

func TestCustomError(t *testing.T) {
	err := NewCustomError(400, "Invalid request parameters")
	log.Printf("Error: %v", err)
	//Error: Invalid request parameters
	//参考https://www.cnblogs.com/YLTFY1998/p/16741285.html
	//=================================================演示errors.IS的用法比较错误是否相同=============================================
	err1 := NewCustomError(400, "Invalid request parameters")
	err2 := NewCustomError(400, "Invalid request parameters1")
	log.Printf("Error: %v", errors.Is(err2, err1))

	//=====================================================演示errors.As错误装换=================================================
	//err3 := errors.New("invalid request parameters")
	err4 := NewCustomError(400, "Invalid request parameters")
	var err5 *E
	//errors.As 的用法，判断err能否转为指定类型。
	//接受两种类型，error和interface类型，interface类型会被反射直接赋值，error类型，会执行第一个参数自定义的as方法如果有的话。
	//将错误转为另一种错误类型。
	r2 := errors.As(err4, &err5)
	log.Printf("as result r2: %v err: %v", r2, err5)
	//2024/05/18 23:06:55 as result r2: true err: E type error

	//演示将错误转换为interface类型，i会被反射赋值。
	var i interface{}
	r3 := errors.As(err4, &i)
	log.Printf("as result r3: %v err: %v", r3, i)
	//2024/05/18 23:06:55 as result r3: true err: Invalid request parameters
	//==========================================wrap用法==================================================
	//wrap给错误包装了一层堆栈,errors.Unwrap 可以获取原始错误
	err6 := NewCustomError(500, "internal error")
	//使用fmt.Errorf有同样的效果,errors.Is 和 errors.As都是会逐层unwrap的来判断原始错误。
	err7 := errors.Wrap(err6, "wrap error")
	log.Printf("Error: err7 %v", err7)
	//2024/05/18 23:06:55 Error: err7 wrap error: internal error

	r4 := errors.Is(err7, err6)
	log.Printf("result r3 = %v", r4)
	//2024/05/18 23:16:11 result r3 = true
}
