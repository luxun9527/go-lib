package validatorx

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
)

type Validator struct {
	validate *validator.Validate
	trans ut.Translator
}
func NewValidator()*Validator{
	en := zh.New()
	uni := ut.New( en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	zh_translations.RegisterDefaultTranslations(validate, trans)
	return &Validator{validate: validate,trans:trans}
}

func (v *Validator)Validate(r *http.Request, data any) error{
	if err := v.validate.Struct(data);err!=nil{
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errs,ok := err.(validator.ValidationErrors)
		if ok{
			if len(errs) > 0 {
			return 	errors.New( errs[0].Translate(v.trans))
			}
		}
		return err
	}

	return nil
}