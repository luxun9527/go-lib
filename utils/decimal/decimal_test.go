package decimal

import (
	"github.com/shopspring/decimal"
	"log"
	"testing"
)

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
func TestDecimal(t *testing.T) {
	//使用float类型来表示小数 和-2的53 ~ 2的53范围外的数，是有精度问题的。示例
	var fv float64 = 74080999549435905
	log.Printf("%v", int64(fv))
	//74080999549435904
	//涉及到浮点数的计算，最好都使用decimal库，转为大数来处理。
	d1, _ := decimal.NewFromString("120.5")
	d := decimal.New(1000, 3)
	log.Printf("%v", d)
	d2, _ := decimal.NewFromString("120.11")
	d2.GreaterThan(d1)
	log.Printf("result %v", d2.RoundBank(3).String())
}
