package i18n

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestI18n1(t *testing.T) {
	//初始化设置默认语言
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	//加载语言包和MustLocalize解析不是并发安全的，如果如果希望多语言文件热更新，需要加锁
	bundle.MustLoadMessageFile("./zh.toml")
	bundle.LoadMessageFile("./en.toml")

	{
		localizer := i18n.NewLocalizer(bundle, "en-US")
		msg, err := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "default1",
				Other: "test", //如果没有找到则使用这个，默认语言不会报错误，非默认语言会报错
			},
		})

		if err != nil {
			log.Printf("%v", err)
			return
		}
		fmt.Println(msg)
		//test
	}
	{
		localizer := i18n.NewLocalizer(bundle, "zh")
		msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "HelloWorld1"})
		if err != nil {
			log.Printf("%v", err)
			//message "HelloWorld1" not found in language "zh"
			return
		}
		fmt.Println(msg)
	}

	{
		localizer := i18n.NewLocalizer(bundle, "en-US")
		msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "HelloWorld"})
		if err != nil {
			return
		}
		fmt.Println(msg)
	}

}

// curl http://localhost:8080/getUserInfo?userId=11&lang=zh
func TestHttpError(t *testing.T) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	//加载语言包和MustLocalize解析不是并发安全的，如果希望热更新，需要加锁
	bundle.MustLoadMessageFile("./zh.toml")
	bundle.MustLoadMessageFile("./en.toml")
	e := gin.New()
	e.GET("/getUserInfo", func(c *gin.Context) {
		lang := c.Query("lang")
		userId := c.Query("userId")

		l := i18n.NewLocalizer(bundle, lang)
		if userId == "1" {
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "",
				"data": gin.H{"username": "张三"},
			})
			return
		}

		msg, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: "100001",
		})
		if err != nil {
			return
		}
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  msg,
		})

	})
	e.Run(":8080")
}
func TestHttpConcurrence(t *testing.T) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	//加载语言包和MustLocalize解析不是并发安全的，如果希望热更新，需要加锁
	bundle.MustLoadMessageFile("./zh.toml")
	bundle.MustLoadMessageFile("./en.toml")
	e := gin.New()
	e.GET("/getUserInfo", func(c *gin.Context) {
		lang := c.Query("lang")
		userId := c.Query("userId")

		l := i18n.NewLocalizer(bundle, lang)

		if userId == "1" {
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "",
				"data": gin.H{"username": "张三"},
			})
			return
		}

		msg, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: "100001",
		})
		if err != nil {
			return
		}
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  msg,
		})

	})
	e.Run(":8080")
}
func TestN(t *testing.T) {

	_ = filepath.WalkDir("E:\\demoproject\\go-lib\\utils\\i18n", func(path string, d os.DirEntry, err error) error {
		log.Printf("%s\n", path)
		log.Printf("%s\n", d.Name())
		log.Printf("%s\n", filepath.Dir(d.Name()))
		return err
	})
}
