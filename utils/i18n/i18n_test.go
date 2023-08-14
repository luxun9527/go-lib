package i18n

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"testing"
)

func TestI18n(t *testing.T) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("en.json")
	bundle.MustLoadMessageFile("el.json")
	loc := i18n.NewLocalizer(bundle, language.English.String(), language.French.String())

	translation := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "messages",
		TemplateData: map[string]interface{}{
			"Name":  "Theo",
			"Count": 2,
		},
		PluralCount: 2,
	})
	log.Println(translation)
}
