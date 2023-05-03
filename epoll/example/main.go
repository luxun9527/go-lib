package main

func main() {
	epollM := NewEpollM()
	//开启监听
	err := epollM.Listen("0.0.0.0", 8088)
	if err != nil {
		panic(err)
	}

	//创建epoll
	err = epollM.CreateEpoll()
	if err != nil {
		panic(err)
	}

	//异步处理epoll
	go func() {
		err := epollM.HandlerEpoll()
		epollM.Close()
		panic(err)
	}()

	//等待client的连接
	err = epollM.Accept()
	epollM.Close()
	panic(err)
}
