package i18n

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/atomic"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultMsgId  = "999999"
	defaultFormat = "toml"
	unKnownMsg    = "internal error"
)

var (
	_defaultLang             = language.English
	_defaultAtomicTranslator = atomic.Pointer[Translator]{}
)

func SetDefaultTranslator(translator *Translator) {
	_defaultAtomicTranslator.Store(translator)
}

type Translator struct {
	*i18n.Bundle
	defaultMsgs map[string]string
}
type LangData struct {
	Lang string
	Data []byte
}

func NewTranslatorFormBytes(data []*LangData) (*Translator, error) {
	bundle := i18n.NewBundle(_defaultLang)
	bundle.RegisterUnmarshalFunc(defaultFormat, toml.Unmarshal)
	var defaultMsgs = map[string]string{}
	for _, v := range data {
		msgFile, err := bundle.ParseMessageFileBytes(v.Data, v.Lang+"."+defaultFormat)
		if err != nil {
			return nil, err
		}
		for _, v := range msgFile.Messages {
			if v.ID == defaultMsgId {
				lang := msgFile.Tag.String()
				defaultMsgs[lang] = v.Other
				break
			}

		}
	}
	return &Translator{Bundle: bundle, defaultMsgs: defaultMsgs}, nil

}

func NewTranslatorFormFile(langFilePath string) (*Translator, error) {
	bundle := i18n.NewBundle(_defaultLang)
	bundle.RegisterUnmarshalFunc(defaultFormat, toml.Unmarshal)
	var defaultMsgs = map[string]string{}
	err := filepath.WalkDir(langFilePath, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if ext := strings.TrimLeft(filepath.Ext(d.Name()), "."); ext != defaultFormat {
			return nil
		}
		msgFile, err := bundle.LoadMessageFile(path)
		if err != nil {
			return err
		}

		for _, v := range msgFile.Messages {
			if v.ID == defaultMsgId {
				lang := msgFile.Tag.String()
				defaultMsgs[lang] = v.Other
				break
			}

		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Translator{Bundle: bundle, defaultMsgs: defaultMsgs}, nil

}

func (t *Translator) Translate(lang, msgId string) string {
	localizer := i18n.NewLocalizer(t.Bundle, lang)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    msgId,
			Other: t.defaultMsgs[lang],
		},
	})
	var e *i18n.MessageNotFoundErr
	if err == nil || (errors.As(err, &e) && msg != "") {
		return msg
	}
	return unKnownMsg
}

func Translate(lang, msgId string) string {
	return _defaultAtomicTranslator.Load().Translate(lang, msgId)
}
