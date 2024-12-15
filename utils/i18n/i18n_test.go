package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"testing"
)

func TestI18n1(t *testing.T) {
	//初始化设置默认语言
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	//加载语言包和MustLocalize解析不是并发安全的，如果如果希望多语言文件热更新，需要加锁
	bundle.MustLoadMessageFile("./zh.toml")
	_, _ = bundle.LoadMessageFile("./en.toml")
	{
		localizer := i18n.NewLocalizer(bundle, "en-US")
		msg, err := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "100001",
			},
		})

		if err != nil {
			log.Printf("translate en-US msg failed %v", err)
			return
		}
		log.Printf("translate en-US success msg:%v", msg)
		//default value
	}
	{
		localizer := i18n.NewLocalizer(bundle, "en-US")
		msg, err := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "noExistedKey",
				Other: "default value", //如果没有找到则使用这个，默认语言不会报错误，非默认语言会报错
			},
		})

		if err != nil {
			log.Printf("translate msg failed %v", err)
			return
		}
		log.Printf("translate success msg:%v", msg)
		//default value
	}
	{
		localizer := i18n.NewLocalizer(bundle, "zh")
		msg, err := localizer.Localize(&i18n.LocalizeConfig{
			TemplateData: nil,
			PluralCount:  nil,
			DefaultMessage: &i18n.Message{
				ID:    "noExistedKey",
				Other: "default value",
			},
			Funcs:          nil,
			TemplateParser: nil,
		})
		// //如果没有找到则使用这个，默认语言不会报错误，非默认语言会报错i18n.MessageNotFoundErr,但是msg还是有值
		if err != nil {
			log.Printf("translate zh msg failed %v,msg %s", err, msg)
			return
		}
		log.Printf("translate zh msg success msg %s", msg)
	}

}
