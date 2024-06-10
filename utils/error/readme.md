# go error实践

相关代码地址，[代码地址](https://github.com/luxun9527/go-lib/tree/master/utils/errors)，如果觉得对你有帮助，欢迎给我的 GitHub 仓库点个 Star ⭐！你的支持是我持续改进和发布更多优质内容的动力。感谢你的关注和支持！

## 1、error相关api

https://www.cnblogs.com/YLTFY1998/p/16741285.html

### is/as

```
func Is(err, target error) bool
```

判断两个错误是否是一样的，如果错误有wrap,会递归unwrap,拿到原始错误再比较。如果错误有实现Is()接口，会调用来比较

```
func As(err error, target interface{}) bool
```

errors.As 的用法，判断err能否转为指定类型。如果能则转换，	接受两种类型，error和interface类型，interface类型会被反射直接赋值，一定能转。error类型，会执行第一个参数自定义的as方法如果有的话。



### wrap/unwrap

```
github.com/pkg/errors
```

errors.Wrap()/fmt.Errorf() 使用这两个方法，可以给err包装一层，给error提供更多的信息。使用Unwrap可以获取原始的错误。

**以下是示例代码**

custom_error.go

```protobuf
package error

import (
	"errors"
	"fmt"
)

type CustomError struct {
	Code    int32
	Message string
}

// 自定义错误类型需要实现error接口
func (e CustomError) Error() string {
	return e.Message
}
func (e CustomError) String() string {
	return fmt.Sprintf(`{code:"%v",message:"%v"}`, e.Code, e.Message)
}

// errors.Is()会调用这个方法来判断是否是同一个错误类型
func (e CustomError) Is(err error) bool {
	var customError CustomError
	if ok := errors.As(err, &customError); !ok {
		return false
	}
	return e.Code == customError.Code

}

// 提取指定类型的错误，判断包装的 error 链中，某一个 error 的类型是否与 target 相同，并提取第一个符合目标类型的错误的值，将其赋值给 target。
func (e CustomError) As(target interface{}) bool {

	t, ok := target.(**E)
	if ok {
		*t = &E{}
	}
	return true
}

// 解析错误，获取原始错误
func (e CustomError) Unwrap() error {
	return e
}

func NewCustomError(code int32, message string) CustomError {
	return CustomError{
		Code:    code,
		Message: message,
	}

}

type E struct {
}

func (E) Error() string {
	return "E type error"
}
```



```protobuf
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
```

## 2、go error实践

**以下是我在个人在日常开发中觉得比较好的实践。**

1、预定义错误使用grpc的status来预定义,预定义错误，错误码按照模块区分。api和rpc的业务错误都使用预定的grpc status错误，这样api不用关心grpc的错误直接抛出就可以，在api统一处理。

2、后端的grpc返回业务错误的时候使用error来返回，而不是在message中定义code,message的形式

3、多语言场景，翻译在api层统一对错误处理。如果需要翻译支持动态修改，修改无需重启程序可以使用fsnotify，etcd，nacos来实现。

4、错误统一存放，统一格式，最好使用工具生成。如go:generate 加上stringer，如果不能满足自己的需求可以使用ast包来自己自定义。



```go
package error

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Code codes.Code

func (c Code) Translate(lang string) string {
	d, ok := _m[lang]
	if !ok {
		d = _m["zh-CN"] //默认中文
	}
	result, ok := d[c]
	if !ok {
		return "服务器繁忙,请稍后再试"
	}
	return result
}

func (c Code) Error(msg string) error {
	return status.New(codes.Code(c), msg).Err()
}

var _m = map[string]map[Code]string{
	"zh-CN": {
		UserNotFoundCode: "用户不存在",
		InternalCode:     "服务器繁忙,请稍后再试",
	},
}
//go:generate errgen -p common.go

package error

const (
    CommonCodeInit Code = 100000 * (iota + 1)
    AccountCodeInit
)

const (
    InternalCode = CommonCodeInit + iota + 1
    RedisErrCode
    ExecSqlFailedCode
    ParamValidateFailedCode
    RecordNotFoundErrCode
    DuplicateDataErrCode
    MongoErrCode
    KafkaErrCode
    EtcdErrCode
    DtmErrCode
    //未知错误 如果没有多语言的需求，且不出文件中加载错误文案，
    //可以使用ast解析注释的方式,将注释解析为文案。
    PulsarErrCode
)

var (
    Internal            = InternalCode.Error("")
    RedisErr            = RedisErrCode.Error("")
    ExecSqlFailed       = ExecSqlFailedCode.Error("")
    ParamValidateFailed = ParamValidateFailedCode.Error("")
    RecordNotFoundErr   = RecordNotFoundErrCode.Error("")
    DuplicateDataErr    = DuplicateDataErrCode.Error("")
    MongoErr            = MongoErrCode.Error("")
    KafkaErr            = KafkaErrCode.Error("")
    EtcdErr             = EtcdErrCode.Error("")
    DtmErr              = DtmErrCode.Error("")
    PulsarErr           = PulsarErrCode.Error("未知错误")
)
```