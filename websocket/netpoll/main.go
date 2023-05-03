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
	"github.com/cloudwego/netpoll"
	"time"
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

type Connection struct {
	isUpgrade bool
}

func prepare(connection netpoll.Connection) context.Context {
	return context.Background()
}

func handle(ctx context.Context, connection netpoll.Connection) error {

	reader := connection.Reader()
	defer reader.Release()
	msg, _ := reader.ReadString(reader.Len())
	println(msg)
	return nil
}
