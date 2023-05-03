package main

//使用tcp 传递稀疏文件
func main() {

}

//类似于一个rpc
//定义数据包格式，比较关键的信息

type Message struct {
	offset uint64
	len    uint64
	data   []byte
}

//1、seek 起始位置 2、数据的长度 前10个字节表示起始位置 后10个字节表示数据的长度。
func (m Message) encode() []byte {

}
func (m Message) decode() {

}
