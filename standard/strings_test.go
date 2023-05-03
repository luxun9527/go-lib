package standard

import (
	"fmt"
	"log"
	"testing"
	"unicode"
)

func TestSpecial(t *testing.T) {
	str := "("
	for _, v := range str {
		if v == '(' {
			log.Println(true)
		}
		log.Println(v)
	}
	log.Println(byte('('))
}

// 如果存在特殊字符，直接在特殊字符前添加\
/**
判断是否为字母： unicode.IsLetter(v)
判断是否为十进制数字： unicode.IsDigit(v)
判断是否为数字： unicode.IsNumber(v)
判断是否为空白符号： unicode.IsSpace(v)
判断是否为Unicode标点字符 :unicode.IsPunct(v)
判断是否为中文：unicode.Han(v)

————————————————
*/
func SpecialLetters(letter rune) (bool, []rune) {
	if unicode.IsPunct(letter) || unicode.IsSymbol(letter) || unicode.Is(unicode.Han, letter) {
		var chars []rune
		chars = append(chars, '\\', letter)
		return true, chars
	}
	return false, nil
}

func TestSpecailLetter(t *testing.T) {
	str := `Admin123!@#$%^&*()-=_+[]{};'\:"|,./<>?`
	var chars []rune
	for _, letter := range str {
		ok, letters := SpecialLetters(letter)
		if ok {
			chars = append(chars, letters...)
		} else {
			chars = append(chars, letter)
		}
	}
	fmt.Println(string(chars))
}
