package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	//	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	en_trans "github.com/go-playground/validator/v10/translations/en"

	"reflect"
	"strings"
	"sync"
)

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
	trans    ut.Translator
}

type SliceValidationError []error

// Error concatenates all error elements in SliceValidationError into a single string separated by \n.
func (err SliceValidationError) Error() string {
	n := len(err)
	switch n {
	case 0:
		return ""
	default:
		var b strings.Builder
		if err[0] != nil {
			fmt.Fprintf(&b, "[%d]: %s", 0, err[0].Error())
		}
		if n > 1 {
			for i := 1; i < n; i++ {
				if err[i] != nil {
					b.WriteString("\n")
					fmt.Fprintf(&b, "[%d]: %s", i, err[i].Error())
				}
			}
		}
		return b.String()
	}
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *DefaultValidator) validateStruct(obj interface{}) error {
	if err := v.validate.Struct(obj); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		firstError := errs[0]
		return errors.New(firstError.Translate(v.trans))
	}
	return nil

}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://pkg.go.dev/github.com/go-playground/validator/v10
func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		validate := validator.New()
		//enTranslator := en.New()
		//zhTranslator := zh.New()
		//uni := ut.New(enTranslator, zhTranslator)
		//
		//// this is usually know or extracted from http 'Accept-Language' header
		//// also see uni.FindTranslator(...)
		//trans, _ := uni.GetTranslator("en")
		e := en.New()
		uni := ut.New(e, e)

		// this is usually know or extracted from http 'Accept-Language' header
		// also see uni.FindTranslator(...)
		trans, _ := uni.GetTranslator("en")
		en_trans.RegisterDefaultTranslations(validate, trans)
		v.validate = validate
		v.trans = trans
		v.validate.SetTagName("binding")
	})
}
