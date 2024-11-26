package i18n

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sync"
	"testing"
)

func TestTranslator_Translate(t *testing.T) {
	translator, err := NewTranslatorFormFile("./")
	if err != nil {
		t.Logf("err:%v\n", err)
		return
	}
	SetDefaultTranslator(translator)
	{
		msg := Translate("zh", "100001")
		log.Printf("msg:%v\n", msg)
		// msg:用户不存在
	}
	{
		msg := Translate("zh", "100002")
		log.Printf("msg:%v\n", msg)
		//msg:内部错误
	}
	{
		msg := Translate("zh", "100003")
		log.Printf("msg:%v\n", msg)
	}
	{
		msg := Translate("en", "100001")
		log.Printf("msg:%v\n", msg)
		//msg:User not found
	}
	{
		msg := Translate("en", "100003")
		log.Printf("msg:%v\n", msg)
		//msg:internal error
	}

	{
		msg := Translate("fr", "100001")
		log.Printf("msg:%v\n", msg)
		//msg:internal error
	}
	{
		msg := Translate("fr", "100001")
		log.Printf("msg:%v\n", msg)
		//msg:User not found 语言没找会使用默认语言。
	}
}
func TestTranslator_ConnTranslate(t *testing.T) {
	// 假设 Translator 实例
	translator1, _ := NewTranslatorFormFile("./")
	translator2, _ := NewTranslatorFormFile("./")

	// 初始化默认 Translator
	SetDefaultTranslator(translator1)

	// 并发测试
	const numGoroutines = 1000
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(2 * numGoroutines)

	// 并发执行 Translate
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = Translate("en", "999999") // 使用默认的 msgId
			}
		}(i)
	}

	// 并发执行 SetDefaultTranslator
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				if j%2 == 0 {
					SetDefaultTranslator(translator1)
				} else {
					SetDefaultTranslator(translator2)
				}
			}
		}(i)
	}

	// 等待所有 Goroutine 完成
	wg.Wait()

	fmt.Println("Concurrent test completed")

}
func TestGin(t *testing.T) {
	e := gin.New()
	translator, _ := NewTranslatorFormFile("./")
	SetDefaultTranslator(translator)
	e.GET("/getUserInfo", func(c *gin.Context) {
		lang := c.Query("lang")
		userId := c.Query("userId")

		if userId == "1" {
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "",
				"data": gin.H{"username": "张三"},
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 500,
			"msg":  Translate(lang, "100001"),
		})
	})
	e.POST("/updateFileLangInfo", func(c *gin.Context) {
		_ = os.WriteFile("./fr-FR.toml", []byte(`100001="Utilisateur introuvable"`), 0644)
		t1, _ := NewTranslatorFormFile("./")
		SetDefaultTranslator(t1)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "",
		})
	})
	e.POST("/updateBytesLangInfo", func(c *gin.Context) {
		lang := c.PostForm("lang")
		data := c.PostForm("data")
		d := []*LangData{{
			Lang: lang,
			Data: []byte(data),
		}}
		tr, err := NewTranslatorFormBytes(d)
		if err != nil {
			return
		}
		SetDefaultTranslator(tr)

		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "",
		})

	})
	e.Run(":8080")
}
