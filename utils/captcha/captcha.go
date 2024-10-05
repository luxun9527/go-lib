package main

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"log"
)

func main() {
	//生成验证
	var store = base64Captcha.DefaultMemStore
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err := cp.Generate()
	fmt.Printf("id: %v, base64s: %v, answer: %v, err: %v", id, b64s, answer, err)
	//验证验证码
	ok := store.Verify(id, answer, true)
	if ok {
		log.Printf("验证成功")
	} else {
		log.Printf("验证失败")
	}

}
