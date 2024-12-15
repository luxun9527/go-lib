package main

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"testing"
)

func TestTranslate(t *testing.T) {
	validate = validator.New()

	translator, err := NewTranslator(validate)
	if err != nil {
		log.Printf("new translator faile err %v", err)
	}
	type User struct {
		Username string `validate:"required"`
		Password string `validate:"required"`
	}
	var user User
	err = validate.Struct(user)
	if err != nil {

		var errs validator.ValidationErrors
		errors.As(err, &errs)

		for _, e := range errs {
			// can translate each error one at a time.
			msg := translator.Translate("zh", e)
			log.Printf("%s\n", msg)
		}
	}

}
