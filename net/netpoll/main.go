/*
 * Copyright 2021 CloudWeGo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"time"

	"github.com/cloudwego/netpoll"
)

func main() {
	network, address := "tcp", ":8080"
	listener, _ := netpoll.CreateListener(network, address)

	eventLoop, _ := netpoll.NewEventLoop(
		handle,
		netpoll.WithOnPrepare(prepare),
		netpoll.WithReadTimeout(time.Minute),
	)
	netpoll.SetNumLoops(4)
	// start listen loop ...
	eventLoop.Serve(listener)
}

var _ netpoll.OnPrepare = prepare
var _ netpoll.OnRequest = handle

func prepare(connection netpoll.Connection) context.Context {
	return context.Background()
}

//想实现的功能读取指定的长度，如果没有读取到执行的长度，这缓存数据返回，下次从这个地方继续读。
func handle(ctx context.Context, connection netpoll.Connection) error {
	reader := connection.Reader()
	defer reader.Release()
	msg, _ := reader.ReadString(4)
	println(msg)
	return nil
}
