package standard

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
)

func TestCal(t *testing.T) {
	for i := 101; i > 0; i -= 100 {
		log.Println(i)
	}
	var i float32 = 2.2

	log.Println(float64(i))
}

func TestFloat(t *testing.T) {
	//小数的二进制算法。小数位乘以2 取整数位
	//1.3
	/*
					整数部分无符号0000 0001
					小数部分
						0.3 * 2=0.6整数位 0 .0
						0.6 * 2=1.2整数位 1 .01
						0.2 * 2 0.4整数位 0 .010
						0.4 * 2 0.8整数位 0 .0100
						0.8 * 2 1.6整数位 1 .01001
						0.6 * 2 1.2整数位 1 .010011
				        1.0100110011001100110011001100110011001100110011001101
						https://www.cnblogs.com/xkfz007/archive/2012/02/27/2370357.html
						https://polarisxu.studygolang.com/posts/basic/diagram-float-point/
				    	https://baseconvert.com/ieee-754-floating-point

					浮点数的内存表示
				     单精度
				      符号位   指数位		尾数23位
						120.5的内存标识
				    二进制表示 1111000.1
					1.1110001 *2的6次方
					正数无符号 0
					6+127=133 对应1000 0101
					尾数 1110001
					0 10000101 11100010000000000000000

				128 64 32 16 8 4 2 1
		                           2的0
				74080999549435905
				双精度
				 符号位   指数位11位		     				 尾数52位
				 0      0000 0000 000 这个位置为2的53次方     0000 0000 0000 0000 0000 000
				二进制
				1 00000111 00110000 01000110 11001110 11000000 00000000 00000001
			 	转为二进制的科学技术法
			 	i := 1.00000111 00110000 01000110 11001110 11000000 00000000 00000001 * 56
				符号位0
				指数位12位 56+1023  1079 =10000110111
				尾数00000111 00110000 01000110 11001110 11000000 00000000 00000001
				最终内存表示
				0 10000110111 00000111 00110000 01000110 11001110 11000000 00000000 00000001
				由于float64占内存8个字节 尾数要舍弃，如果尾数大于2的53次方，会有精度丢失的风险。
				1 00000111 00110000 01000110 11001110 11000000 00000000 0000
				10进制表示 74080999549435904
	*/
	var fv float64 = 74080999549435905

	log.Println(strconv.FormatFloat(fv, 'f', -1, 64))
	//74080999549435900 strconv.FormatFloat方法会四舍五入 float64能精确表示的范围-2的53 ~ 2的53
	var iv int64 = int64(fv)
	log.Println(iv)
	//74080999549435904
}

// 将整数部分转换为二进制字符串
func intToBinary(integerPart int64) string {
	return strconv.FormatInt(integerPart, 2)
}

// 将小数部分转换为二进制字符串
func fracToBinary(fractionalPart float64, precision int) string {
	var binaryBuilder strings.Builder
	for i := 0; i < precision; i++ {
		fractionalPart *= 2
		if fractionalPart >= 1 {
			binaryBuilder.WriteString("1")
			fractionalPart -= 1
		} else {
			binaryBuilder.WriteString("0")
		}
		if fractionalPart == 0 {
			break
		}
	}
	return binaryBuilder.String()
}

// 将十进制浮点数转换为二进制字符串
func floatToBinary(decimal float64, precision int) string {
	// 分离整数部分和小数部分
	integerPart := int64(decimal)
	fractionalPart := decimal - float64(integerPart)

	// 转换整数部分和小数部分
	binaryInteger := intToBinary(integerPart)
	binaryFraction := fracToBinary(fractionalPart, precision)

	// 组合结果
	if binaryFraction == "" {
		return binaryInteger
	}
	return binaryInteger + "." + binaryFraction
}
func TestFloat1(t *testing.T) {
	decimal := 0.1
	precision := 100 // 小数部分精度

	binaryRepresentation := floatToBinary(decimal, precision)
	fmt.Printf("The binary representation of %.10f is %s\n", decimal, binaryRepresentation)
}

//0.0001100110011001100110011001100110011001100110011001101
//0 01111111011 1001100110011001100110011001100110011001100110011010
