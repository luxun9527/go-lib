package reflectdemo

import (
	"fmt"
	"go-lib/standard/reflectdemo/model"
	"log"
	"reflect"
	"testing"
	"unsafe"
)

type student struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func TestReflectStruct(t *testing.T){
	stu1 := student{
		Name:  "小王子",
		Score: 90,
	}

	t2 := reflect.TypeOf(stu1)
	fmt.Println(t2.Name(), t2.Kind()) // student struct
	// 通过for循环遍历结构体的所有字段信息
	for i := 0; i < t2.NumField(); i++ {
		field := t2.Field(i)
		fmt.Printf("name:%s index:%d type:%v json tag:%v\n", field.Name, field.Index, field.Type, field.Tag.Get("json"))
	}

	// 通过字段名获取指定结构体字段信息
	if scoreField, ok := t2.FieldByName("Score"); ok {
		fmt.Printf("name:%s index:%d type:%v json tag:%v\n", scoreField.Name, scoreField.Index, scoreField.Type, scoreField.Tag.Get("json"))
	}

}

//修改接口未导出的字段
func TestUnExport(t *testing.T) {
	 sourceCopy := GetStructPtrUnExportedField(model.E, "unexport")
	//sourceCopy.
	log.Println(sourceCopy.String(), sourceCopy.Type())
	name := sourceCopy.FieldByName("name")
	s := reflect.NewAt( name.Type(), unsafe.Pointer(name.UnsafeAddr())).Elem()
	//reflect.NewAt( sourceCopy.Type(), unsafe.Pointer(sourceCopy.UnsafeAddr())).Elem()
	s.SetString("111")
	log.Println(model.E)
	//指针通过Elem()函数，获得CanSet和CanAddress能力 结构体则需要 先通过reflect.New函数获得一个一样的 reflect.Value
}

func SetStructPtrUnExportedStrField(source interface{}, fieldName string, fieldVal interface{}) (err error) {
	v := GetStructPtrUnExportedField(source, fieldName)
	rv := reflect.ValueOf(fieldVal)
	if v.Kind() != rv.Kind() {
		return fmt.Errorf("invalid kind: expected kind %v, got kind: %v", v.Kind(), rv.Kind())
	}

	v.Set(rv)
	return nil
}

func SetStructUnExportedStrField(source interface{}, fieldName string, fieldVal interface{}) (addressableSourceCopy reflect.Value, err error) {
	var accessableField reflect.Value
	accessableField, addressableSourceCopy = GetStructUnExportedField(source, fieldName)
	rv := reflect.ValueOf(fieldVal)
	if accessableField.Kind() != rv.Kind() {
		return addressableSourceCopy, fmt.Errorf("invalid kind: expected kind %v, got kind: %v", addressableSourceCopy.Kind(), rv.Kind())
	}
	accessableField.Set(rv)
	return
}

func GetStructPtrUnExportedField(source interface{}, fieldName string) reflect.Value {
	v := reflect.ValueOf(source).Elem().FieldByName(fieldName)
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

func GetStructUnExportedField(source interface{}, fieldName string) (accessableField, addressableSourceCopy reflect.Value) {
	v := reflect.ValueOf(source)
	// since source is not a ptr, get an addressable copy of source to modify it later
	log.Println(v.Type())
	addressableSourceCopy = reflect.New(v.Type()).Elem()
	addressableSourceCopy.Set(v)
	accessableField = addressableSourceCopy.FieldByName(fieldName)
	accessableField = reflect.NewAt(accessableField.Type(), unsafe.Pointer(accessableField.UnsafeAddr())).Elem()
	return
}