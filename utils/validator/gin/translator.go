package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/fa"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/id"
	"github.com/go-playground/locales/it"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/lv"
	"github.com/go-playground/locales/nl"
	"github.com/go-playground/locales/pt"
	"github.com/go-playground/locales/pt_BR"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/locales/tr"
	"github.com/go-playground/locales/vi"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	es_translations "github.com/go-playground/validator/v10/translations/es"
	fa_translations "github.com/go-playground/validator/v10/translations/fa"
	fr_translations "github.com/go-playground/validator/v10/translations/fr"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	it_translations "github.com/go-playground/validator/v10/translations/it"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
	lv_translations "github.com/go-playground/validator/v10/translations/lv"
	nl_translations "github.com/go-playground/validator/v10/translations/nl"
	pt_translations "github.com/go-playground/validator/v10/translations/pt"
	ptBR_translations "github.com/go-playground/validator/v10/translations/pt_BR"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
	tr_translations "github.com/go-playground/validator/v10/translations/tr"
	vi_translations "github.com/go-playground/validator/v10/translations/vi"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	zhTW_translations "github.com/go-playground/validator/v10/translations/zh_tw"
)

type Translator struct {
	translators map[string]ut.Translator
}

var (
	_registerFuncs = map[string]func(v *validator.Validate, trans ut.Translator) (err error){
		"zh":         zh_translations.RegisterDefaultTranslations,
		"en":         en_translations.RegisterDefaultTranslations,
		"es":         es_translations.RegisterDefaultTranslations,
		"fa":         fa_translations.RegisterDefaultTranslations,
		"fr":         fr_translations.RegisterDefaultTranslations,
		"id":         id_translations.RegisterDefaultTranslations,
		"it":         it_translations.RegisterDefaultTranslations,
		"ja":         ja_translations.RegisterDefaultTranslations,
		"lv":         lv_translations.RegisterDefaultTranslations,
		"nl":         nl_translations.RegisterDefaultTranslations,
		"pt":         pt_translations.RegisterDefaultTranslations,
		"pt_BR":      ptBR_translations.RegisterDefaultTranslations,
		"ru":         ru_translations.RegisterDefaultTranslations,
		"tr":         tr_translations.RegisterDefaultTranslations,
		"vi":         vi_translations.RegisterDefaultTranslations,
		"zh_Hant_TW": zhTW_translations.RegisterDefaultTranslations,
	}
)

func NewTranslator(validate *validator.Validate) (*Translator, error) {
	var translator Translator
	m := make(map[string]ut.Translator, 20)
	var (
		zhLang   = zh.New()
		enLang   = en.New()
		esLang   = es.New()
		faLang   = fa.New()
		frLang   = fr.New()
		idLang   = id.New()
		itLang   = it.New()
		jaLang   = ja.New()
		lvLang   = lv.New()
		nlLang   = nl.New()
		ptLang   = pt.New()
		ptBrLang = pt_BR.New()
		ruLang   = ru.New()
		trLang   = tr.New()
		viLang   = vi.New()
		twLang   = zh_Hant_TW.New()
	)

	uni := ut.New(enLang, enLang, zhLang, esLang, faLang, frLang, idLang, itLang, jaLang, lvLang, nlLang, ptLang, ptBrLang, ruLang, trLang, viLang, twLang)
	for lang, registerFunc := range _registerFuncs {
		tran, ok := uni.GetTranslator(lang)
		if !ok {
			return nil, fmt.Errorf("%s not found", lang)
		}
		if err := registerFunc(validate, tran); err != nil {
			return nil, err
		}
		m[lang] = tran
	}
	translator.translators = m
	return &translator, nil
}

func toString(messages validator.ValidationErrorsTranslations) string {
	msg := ""
	for k, v := range messages {
		msg += fmt.Sprintf("%s:%s\n", k, v)
	}
	return msg
}

func (tl *Translator) Translate(locale string, err error) string {
	translator, ok := tl.translators[locale]
	if !ok {
		return err.Error()
	}

	var (
		e1 validator.ValidationErrors
		e2 validator.FieldError
	)
	if ok := errors.As(err, &e1); ok {
		m := e1.Translate(translator)
		return toString(m)
	}

	if ok := errors.As(err, &e2); ok {
		m := e2.Translate(translator)
		return m
	}
	return err.Error()

}

func (tl *Translator) TranslateFirst(locale string, err error) string {
	var (
		e1 validator.ValidationErrors
	)
	if ok := errors.As(err, &e1); !ok {
		return err.Error()
	}
	if len(e1) > 0 {
		return tl.Translate(locale, e1[0])
	}
	return err.Error()

}
