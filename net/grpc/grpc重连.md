# GRPC断线重连

refer

https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md

https://github.com/jpillora/backoff

https://zacard.net/2021/01/13/log-agent-backoff/

https://stackoverflow.com/questions/37125975/how-to-debug-grpc-call

当grpc服务端关闭，过了一会重启后，客户端不用重启，总是能够再次连接到服务端。

从源码分析，探究其中发生了什么。



## 先说结论

当服务端关闭，客户端并不会去主动建立连接，会更新picker，为下次建立连接做好准备。只有当客户端下次调用接口的时候，客户端才会开一个新的协程主动去建立连接（后面调用接口，不会开启协程创建连接），当建立失败的时候，会往一个chan中塞一个函数，另一个协程收到函数后，会调用此函数重试建立连接。具体的时候间隔采用指数退避的方式，总体就是失败的越多，下次等待的时候越长。具体的算法参考。ac.dopts.bs.Backoff(ac.backoffIdx)

## 前置操作

grpc版本v1.5.9

### 设置为grpc logging 设置为debug

设置环境变量

export GRPC_GO_LOG_VERBOSITY_LEVEL=99
export GRPC_GO_LOG_SEVERITY_LEVEL=info

### 具体grpc服务端客户端代码

```go
type GrpcDemoServer struct {
    grpcdemo.UnimplementedGrpcDemoServer
}
func (GrpcDemoServer) Call(ctx context.Context,req *grpcdemo.NoticeReaderReq) (*emptypb.Empty, error) {
    return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func TestServer(t *testing.T) {
    listener, err := net.Listen("tcp", "0.0.0.0:8899")
    if err != nil {
        log.Println("net listen err ", err)
        return
    }
    s := grpc.NewServer()
    grpcdemo.RegisterGrpcDemoServer(s,new(GrpcDemoServer))
    if err := s.Serve(listener); err != nil {
        log.Println("failed to serve...", err)
        return
    }
}
func TestClient(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
   conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  1.0 * time.Second,
			Multiplier: 1.6,
			Jitter:     0.2,
			MaxDelay:   30 * time.Second,
		},
	}))
if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
    for {
        time.Sleep(time.Second * 40)
        cli := grpcdemo.NewGrpcDemoClient(conn)
        result, err := cli.Call(context.Background(), &grpcdemo.NoticeReaderReq{
            Msg:       "",
            NoticeWay: &grpcdemo.NoticeReaderReq_Email{Email: "test"},
        })
        if err != nil {
            log.Printf("Call  failed %v", err)
        }
        log.Printf("================================result================================== %v", result)
    }
}
```

## 

### 日志

```go
2023/11/26 23:49:27 INFO: [core] [Channel #1 SubChannel #2] Subchannel created
2023/11/26 23:49:29 INFO: [core] [Channel #1] Channel Connectivity change to CONNECTING
2023/11/26 23:49:29 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to CONNECTING
2023/11/26 23:49:29 INFO: [core] [Channel #1 SubChannel #2] Subchannel picks a new address "127.0.0.1:8899" to connect
2023/11/26 23:49:29 INFO: [core] [pick-first-lb 0xc00027ea80] Received SubConn state update: 0xc00027ebd0, {ConnectivityState:CONNECTING ConnectionError:<nil>}
2023/11/26 23:49:30 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to READY
2023/11/26 23:49:30 INFO: [core] [pick-first-lb 0xc00027ea80] Received SubConn state update: 0xc00027ebd0, {ConnectivityState:READY ConnectionError:<nil>}
2023/11/26 23:49:31 INFO: [core] [Channel #1] Channel Connectivity change to READY
//关闭服务端
2023/11/26 23:49:38 INFO: [transport] [client-transport 0xc000350000] Closing: connection error: desc = "error reading from server: read tcp 127.0.0.1:60811->127.0.0.1:8899: wsarecv: An existing connection was forcibly closed by the remote host."
2023/11/26 23:49:38 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to IDLE
2023/11/26 23:49:38 INFO: [core] [pick-first-lb 0xc00027ea80] Received SubConn state update: 0xc00027ebd0, {ConnectivityState:IDLE ConnectionError:<nil>}
2023/11/26 23:49:39 INFO: [transport] [client-transport 0xc000350000] loopyWriter exiting with error: transport closed by client
2023/11/26 23:49:39 INFO: [core] [Channel #1] Channel Connectivity change to IDLE

//客户端调用方法
2023/11/26 23:53:17 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to CONNECTING
2023/11/26 23:53:17 INFO: [core] [Channel #1 SubChannel #2] Subchannel picks a new address "127.0.0.1:8899" to connect
2023/11/26 23:53:17 INFO: [core] [pick-first-lb 0xc00027ea80] Received SubConn state update: 0xc00027ebd0, {ConnectivityState:CONNECTING ConnectionError:<nil>}
2023/11/26 23:53:18 INFO: [core] [Channel #1] Channel Connectivity change to CONNECTING
2023/11/26 23:53:20 INFO: [core] Creating new client transport to "{Addr: \"127.0.0.1:8899\", ServerName: \"127.0.0.1:8899\", }": connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:8899: connectex: No connection could be made because the target machine actively refused it."
2023/11/26 23:53:20 WARNING: [core] [Channel #1 SubChannel #2] grpc: addrConn.createTransport failed to connect to {Addr: "127.0.0.1:8899", ServerName: "127.0.0.1:8899", }. Err: connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:8899: connectex: No connection could be made because the target machine actively refused it."
2023/11/26 23:53:20 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to TRANSIENT_FAILURE, last error: connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:8899: connectex: No connection could be made because the target machine actively refused it."
2023/11/26 23:53:20 INFO: [core] [pick-first-lb 0xc00027ea80] Received SubConn state update: 0xc00027ebd0, {ConnectivityState:TRANSIENT_FAILURE ConnectionError:connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:8899: connectex: No connection could be made because the target machine actively refused it."}
2023/11/26 23:53:23 INFO: [core] [Channel #1] Channel Connectivity change to TRANSIENT_FAILURE
2023/11/26 23:53:23 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to IDLE, last error: connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:8899: connectex: No connection could be made because the target machine actively refused it."
2023/11/26 23:53:23 Call  failed rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:8899: connectex: No connection could be made because the target machine actively refused it."
```

**从日志中分析，当关闭连接的时候客户端并没有主动去重试建立连接，只有当调用了具体的方法的时候客户端才主动去建立的建立连接。**

## 调用链分析

### 关闭服务端，断开连接

当关闭服务端，断开连接。具体调用链如下

```go
func (t *http2Client) reader(errCh chan<- error) {
	defer close(t.readerDone)

	if err := t.readServerPreface(); err != nil {
		errCh <- err
		return
	}
	close(errCh)
	if t.keepaliveEnabled {
		atomic.StoreInt64(&t.lastRead, time.Now().UnixNano())
	}

	// loop to keep reading incoming messages on this transport.
    //循环读取
	for {
		t.controlBuf.throttle()
		frame, err := t.framer.fr.ReadFrame()
		if t.keepaliveEnabled {
			atomic.StoreInt64(&t.lastRead, time.Now().UnixNano())
		}
		if err != nil {
			// Abort an active stream if the http2.Framer returns a
			// http2.StreamError. This can happen only if the server's response
			// is malformed http2.
			if se, ok := err.(http2.StreamError); ok {
				t.mu.Lock()
				s := t.activeStreams[se.StreamID]
				t.mu.Unlock()
				if s != nil {
					// use error detail to provide better err message
					code := http2ErrConvTab[se.Code]
					errorDetail := t.framer.fr.ErrorDetail()
					var msg string
					if errorDetail != nil {
						msg = errorDetail.Error()
					} else {
						msg = "received invalid frame"
					}
					t.closeStream(s, status.Error(code, msg), true, http2.ErrCodeProtocol, status.New(code, msg), nil, false)
				}
				continue
			} else {
				// Transport error.
                  //===============调用此函数 出现错误触发此函数===================
				t.Close(connectionErrorf(true, err, "error reading from server: %v", err))
				return
			}
		}
// Close kicks off the shutdown process of the transport. This should be called
// only once on a transport. Once it is called, the transport should not be
// accessed any more.
func (t *http2Client) Close(err error) {
	t.mu.Lock()
	// Make sure we only close once.
	if t.state == closing {
		t.mu.Unlock()
		return
	}
	if t.logger.V(logLevel) {
		t.logger.Infof("Closing: %v", err)
	}
	// Call t.onClose ASAP to prevent the client from attempting to create new
	// streams.
	if t.state != draining {
          //===============调用此函数===================
		t.onClose(GoAwayInvalid)
	}
	t.state = closing
	streams := t.activeStreams
	t.activeStreams = nil
	if t.kpDormant {
		// If the keepalive goroutine is blocked on this condition variable, we
		// should unblock it so that the goroutine eventually exits.
		t.kpDormancyCond.Signal()
	}
	t.mu.Unlock()
	t.controlBuf.finish()
	t.cancel()
	t.conn.Close()
```



```go
	onClose := func(r transport.GoAwayReason) {
		ac.mu.Lock()
		defer ac.mu.Unlock()
		// adjust params based on GoAwayReason
		ac.adjustParams(r)
		if ctx.Err() != nil {
			// Already shut down or connection attempt canceled.  tearDown() or
			// updateAddrs() already cleared the transport and canceled hctx
			// via ac.ctx, and we expected this connection to be closed, so do
			// nothing here.
			return
		}
		hcancel()
		if ac.transport == nil {
			// We're still connecting to this address, which could error.  Do
			// not update the connectivity state or resolve; these will happen
			// at the end of the tryAllAddrs connection loop in the event of an
			// error.
			return
		}
		ac.transport = nil
		// Refresh the name resolver on any connection loss.
		ac.cc.resolveNow(resolver.ResolveNowOptions{})
		// Always go idle and wait for the LB policy to initiate a new
		// connection attempt.
        //===============更新状态===================
		ac.updateConnectivityState(connectivity.Idle, nil)
	}
```



```go
// Note: this requires a lock on ac.mu.
func (ac *addrConn) updateConnectivityState(s connectivity.State, lastErr error) {
	if ac.state == s {
		return
	}
	// When changing states, reset the state change channel.
	close(ac.stateChan)
	ac.stateChan = make(chan struct{})
	ac.state = s
	if lastErr == nil {
		channelz.Infof(logger, ac.channelzID, "Subchannel Connectivity change to %v", s)
	} else {
		channelz.Infof(logger, ac.channelzID, "Subchannel Connectivity change to %v, last error: %s", s, lastErr)
	}
       //===============更新服务状态===================
	ac.cc.handleSubConnStateChange(ac.acbw, s, lastErr)
}
// updateSubConnState is invoked by grpc to push a subConn state update to the
// underlying balancer.
func (ccb *ccBalancerWrapper) updateSubConnState(sc balancer.SubConn, s connectivity.State, err error) {
	ccb.mu.Lock()
     //===============添加此函数作为信号===================
	ccb.serializer.Schedule(func(_ context.Context) {
		// Even though it is optional for balancers, gracefulswitch ensures
		// opts.StateListener is set, so this cannot ever be nil.
          
		sc.(*acBalancerWrapper).stateListener(balancer.SubConnState{ConnectivityState: s, ConnectionError: err})
	})
	ccb.mu.Unlock()
}
// Schedule adds a callback to be scheduled after existing callbacks are run.
//
// Callbacks are expected to honor the context when performing any blocking
// operations, and should return early when the context is canceled.
//
// Return value indicates if the callback was successfully added to the list of
// callbacks to be executed by the serializer. It is not possible to add
// callbacks once the context passed to NewCallbackSerializer is cancelled.
func (cs *CallbackSerializer) Schedule(f func(ctx context.Context)) bool {
	cs.closedMu.Lock()
	defer cs.closedMu.Unlock()

	if cs.closed {
		return false
	}
      //===============添加此函数作为信号===================
	cs.callbacks.Put(f)
	return true
}
// Put adds t to the unbounded buffer.
func (b *Unbounded) Put(t any) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closing {
		return errBufferClosed
	}
	if len(b.backlog) == 0 {
		select {
		case b.c <- t:
			return nil
		default:
		}
	}
	b.backlog = append(b.backlog, t)
	return nil
}
func (cs *CallbackSerializer) run(ctx context.Context) {
	defer close(cs.done)

	// TODO: when Go 1.21 is the oldest supported version, this loop and Close
	// can be replaced with:
	//
	// context.AfterFunc(ctx, cs.callbacks.Close)
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			// Do nothing here. Next iteration of the for loop will not happen,
			// since ctx.Err() would be non-nil.
            //===============收到信号，获取推过来的函数===================
		case cb := <-cs.callbacks.Get():
			cs.callbacks.Load()
             //===============调用推过来的函数===================
			cb.(func(context.Context))(ctx)
		}
	}

	// Close the buffer to prevent new callbacks from being added.
	cs.callbacks.Close()

	// Run all pending callbacks.
	for cb := range cs.callbacks.Get() {
		cs.callbacks.Load()
		cb.(func(context.Context))(ctx)
	}
}
```



```go
// updateSubConnState is invoked by grpc to push a subConn state update to the
// underlying balancer.
func (ccb *ccBalancerWrapper) updateSubConnState(sc balancer.SubConn, s connectivity.State, err error) {
	ccb.mu.Lock()
	ccb.serializer.Schedule(func(_ context.Context) {
		// Even though it is optional for balancers, gracefulswitch ensures
		// 这个是具体的函数。
		sc.(*acBalancerWrapper).stateListener(balancer.SubConnState{ConnectivityState: s, ConnectionError: err})
	})
	ccb.mu.Unlock()
}
// updateSubConnState forwards the update to the appropriate child.
func (gsb *Balancer) updateSubConnState(sc balancer.SubConn, state balancer.SubConnState, cb func(balancer.SubConnState)) {
	gsb.currentMu.Lock()
	defer gsb.currentMu.Unlock()
	gsb.mu.Lock()
	// Forward update to the appropriate child.  Even if there is a pending
	// balancer, the current balancer should continue to get SubConn updates to
	// maintain the proper state while the pending is still connecting.
	var balToUpdate *balancerWrapper
	if gsb.balancerCurrent != nil && gsb.balancerCurrent.subconns[sc] {
		balToUpdate = gsb.balancerCurrent
	} else if gsb.balancerPending != nil && gsb.balancerPending.subconns[sc] {
		balToUpdate = gsb.balancerPending
	}
	if balToUpdate == nil {
		// SubConn belonged to a stale lb policy that has not yet fully closed,
		// or the balancer was already closed.
		gsb.mu.Unlock()
		return
	}
	if state.ConnectivityState == connectivity.Shutdown {
		delete(balToUpdate.subconns, sc)
	}
	gsb.mu.Unlock()
	if cb != nil {
        //======================走到了这里=========================
		cb(state)
	} else {
		balToUpdate.UpdateSubConnState(sc, state)
	}
}
func (b *pickfirstBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	addrs := state.ResolverState.Addresses
	if len(addrs) == 0 {
		// The resolver reported an empty address list. Treat it like an error by
		// calling b.ResolverError.
		if b.subConn != nil {
			// Shut down the old subConn. All addresses were removed, so it is
			// no longer valid.
			b.subConn.Shutdown()
			b.subConn = nil
		}
		b.ResolverError(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	}

	// We don't have to guard this block with the env var because ParseConfig
	// already does so.
	cfg, ok := state.BalancerConfig.(pfConfig)
	if state.BalancerConfig != nil && !ok {
		return fmt.Errorf("pickfirst: received illegal BalancerConfig (type %T): %v", state.BalancerConfig, state.BalancerConfig)
	}
	if cfg.ShuffleAddressList {
		addrs = append([]resolver.Address{}, addrs...)
		grpcrand.Shuffle(len(addrs), func(i, j int) { addrs[i], addrs[j] = addrs[j], addrs[i] })
	}

	if b.logger.V(2) {
		b.logger.Infof("Received new config %s, resolver state %s", pretty.ToJSON(cfg), pretty.ToJSON(state.ResolverState))
	}

	if b.subConn != nil {
		b.cc.UpdateAddresses(b.subConn, addrs)
		return nil
	}

	var subConn balancer.SubConn
	subConn, err := b.cc.NewSubConn(addrs, balancer.NewSubConnOptions{
		StateListener: func(state balancer.SubConnState) {
            //===============走到了这里=================
			b.updateSubConnState(subConn, state)
		},
	})
func (b *pickfirstBalancer) updateSubConnState(subConn balancer.SubConn, state balancer.SubConnState) {
	if b.logger.V(2) {
		b.logger.Infof("Received SubConn state update: %p, %+v", subConn, state)
	}
	if b.subConn != subConn {
		if b.logger.V(2) {
			b.logger.Infof("Ignored state change because subConn is not recognized")
		}
		return
	}
	if state.ConnectivityState == connectivity.Shutdown {
		b.subConn = nil
		return
	}

	switch state.ConnectivityState {
	case connectivity.Ready:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{result: balancer.PickResult{SubConn: subConn}},
		})
	case connectivity.Connecting:
		if b.state == connectivity.TransientFailure {
			// We stay in TransientFailure until we are Ready. See A62.
			return
		}
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
		})
	case connectivity.Idle:
		if b.state == connectivity.TransientFailure {
			// We stay in TransientFailure until we are Ready. Also kick the
			// subConn out of Idle into Connecting. See A62.
			b.subConn.Connect()
			return
		}
        //===========走到了这里================
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &idlePicker{subConn: subConn},
		})
	case connectivity.TransientFailure:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{err: state.ConnectionError},
		})
	}
	b.state = state.ConnectivityState
}
func (ccb *ccBalancerWrapper) UpdateState(s balancer.State) {
    if ccb.isIdleOrClosed() {
        return
    }

    // Update picker before updating state.  Even though the ordering here does
    // not matter, it can lead to multiple calls of Pick in the common start-up
    // case where we wait for ready and then perform an RPC.  If the picker is
    // updated later, we could call the "connecting" picker when the state is
    // updated, and then call the "ready" picker after the picker gets updated.
    //走到了这里更新了picker
    ccb.cc.pickerWrapper.updatePicker(s.Picker)
    ccb.cc.csMgr.updateState(s.ConnectivityState)
}
```

小结，当出现服务端关闭的时候，会触发更新picker，当我们下次调用grpc接口的时候会重连。

### 客户端调用接口

当我们调用grpc的接口的时候，调用invoke方法的调用链

```go
//不重要的地方简单说明
func newClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (_ ClientStream, err error) {}

func newClientStreamWithParams(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, mc serviceconfig.MethodConfig, onCommit, doneFunc func(), opts ...CallOption) (_ iresolver.ClientStream, err error) {
// Pick the transport to use and create a new stream on the transport.
	// Assign cs.attempt upon success.
	op := func(a *csAttempt) error {
        //走到了这里
		if err := a.getTransport(); err != nil {
			return err
		}
		if err := a.newStream(); err != nil {
			return err
		}
		// Because this operation is always called either here (while creating
		// the clientStream) or by the retry code while locked when replaying
		// the operation, it is safe to access cs.attempt directly.
		cs.attempt = a
		return nil
	}
	if err := cs.withRetry(op, func() { cs.bufferForRetryLocked(0, op) }); err != nil {
		return nil, err
	}
	
    
}

func (pw *pickerWrapper) pick(ctx context.Context, failfast bool, info balancer.PickInfo) (transport.ClientTransport, balancer.PickResult, error) {


    		// If the channel is set, it means that the pick call had to wait for a
		// new picker at some point. Either it's the first iteration and this
		// function received the first picker, or a picker errored with
		// ErrNoSubConnAvailable or errored with failfast set to false, which
		// will trigger a continue to the next iteration. In the first case this
		// conditional will hit if this call had to block (the channel is set).
		// In the second case, the only way it will get to this conditional is
		// if there is a new picker.
		if ch != nil {
			for _, sh := range pw.statsHandlers {
				sh.HandleRPC(ctx, &stats.PickerUpdated{})
			}
		}

		ch = pw.blockingCh
		p := pw.picker
		pw.mu.Unlock()
    	//走到了这里
		pickResult, err := p.Pick(info)
}
```



```go
func (i *idlePicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	//走到了这里
    i.subConn.Connect()
	return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
}
```



```go
func (acbw *acBalancerWrapper) Connect() {
    go acbw.ac.connect()
}
// connect starts creating a transport.
// It does nothing if the ac is not IDLE.
// TODO(bar) Move this to the addrConn section.
func (ac *addrConn) connect() error {
    ac.mu.Lock()
    if ac.state == connectivity.Shutdown {
        if logger.V(2) {
            logger.Infof("connect called on shutdown addrConn; ignoring.")
        }
        ac.mu.Unlock()
        return errConnClosing
    }
    if ac.state != connectivity.Idle {
        if logger.V(2) {
            logger.Infof("connect called on addrConn in non-idle state (%v); ignoring.", ac.state)
        }
        ac.mu.Unlock()
        return nil
    }
    ac.mu.Unlock()
//走到了这里
    ac.resetTransport()
    return nil
}
func (ac *addrConn) resetTransport() {
	ac.mu.Lock()
	acCtx := ac.ctx
	if acCtx.Err() != nil {
		ac.mu.Unlock()
		return
	}

	addrs := ac.addrs
    ////走到了这里，指数避退，失败次数越多下一次等待调用的时间越长
	backoffFor := ac.dopts.bs.Backoff(ac.backoffIdx)
	// This will be the duration that dial gets to finish.
	dialDuration := minConnectTimeout
	if ac.dopts.minConnectTimeout != nil {
		dialDuration = ac.dopts.minConnectTimeout()
	}

	if dialDuration < backoffFor {
		// Give dial more time as we keep failing to connect.
		dialDuration = backoffFor
	}
	// We can potentially spend all the time trying the first address, and
	// if the server accepts the connection and then hangs, the following
	// addresses will never be tried.
	//
	// The spec doesn't mention what should be done for multiple addresses.
	// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md#proposed-backoff-algorithm
	connectDeadline := time.Now().Add(dialDuration)

	ac.updateConnectivityState(connectivity.Connecting, nil)
	ac.mu.Unlock()
	// =======================重连=================================
	if err := ac.tryAllAddrs(acCtx, addrs, connectDeadline); err != nil {
        //如果失败还是走到这里
		ac.cc.resolveNow(resolver.ResolveNowOptions{})
		ac.mu.Lock()
		if acCtx.Err() != nil {
			// addrConn was torn down.
			ac.mu.Unlock()
			return
		}
		// After exhausting all addresses, the addrConn enters
		// TRANSIENT_FAILURE.
		ac.updateConnectivityState(connectivity.TransientFailure, err)

		// Backoff.
		b := ac.resetBackoff
		ac.mu.Unlock()

		timer := time.NewTimer(backoffFor)
		select {
		case <-timer.C:
			ac.mu.Lock()
			ac.backoffIdx++
			ac.mu.Unlock()
		case <-b:
			timer.Stop()
		case <-acCtx.Done():
			timer.Stop()
			return
		}

		ac.mu.Lock()
		if acCtx.Err() == nil {
            //重连失败了。走到了这里
			ac.updateConnectivityState(connectivity.Idle, err)
		}
		ac.mu.Unlock()
		return
	}
// Note: this requires a lock on ac.mu.
func (ac *addrConn) updateConnectivityState(s connectivity.State, lastErr error) {
	if ac.state == s {
		return
	}
	// When changing states, reset the state change channel.
	close(ac.stateChan)
	ac.stateChan = make(chan struct{})
	ac.state = s
	if lastErr == nil {
		channelz.Infof(logger, ac.channelzID, "Subchannel Connectivity change to %v", s)
	} else {
		channelz.Infof(logger, ac.channelzID, "Subchannel Connectivity change to %v, last error: %s", s, lastErr)
	}
    //============================走到了这里======================
	ac.cc.handleSubConnStateChange(ac.acbw, s, lastErr)
}
```





```go
// updateSubConnState is invoked by grpc to push a subConn state update to the
// underlying balancer.
func (ccb *ccBalancerWrapper) updateSubConnState(sc balancer.SubConn, s connectivity.State, err error) {
	ccb.mu.Lock()
	//========================走到了这里，发送信号，触发重试。============================
    // 收到函数
    ccb.serializer.Schedule(func(_ context.Context) {
		// Even though it is optional for balancers, gracefulswitch ensures
		// opts.StateListener is set, so this cannot ever be nil.
		//
        sc.(*acBalancerWrapper).stateListener(balancer.SubConnState{ConnectivityState: s, ConnectionError: err})
	})
	ccb.mu.Unlock()
}
func (cs *CallbackSerializer) run(ctx context.Context) {
	defer close(cs.done)

	// TODO: when Go 1.21 is the oldest supported version, this loop and Close
	// can be replaced with:
	//
	// context.AfterFunc(ctx, cs.callbacks.Close)
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			// Do nothing here. Next iteration of the for loop will not happen,
			// since ctx.Err() would be non-nil.
            //===============收到信号，获取推过来的函===================
		case cb := <-cs.callbacks.Get():
			cs.callbacks.Load()
             //===============执行推过来的函数，具体的调用链和断线执行的时候执行差不多===================
			//只有再执行下面的函数有不同
            cb.(func(context.Context))(ctx)
		}
	}
func (b *pickfirstBalancer) updateSubConnState(subConn balancer.SubConn, state balancer.SubConnState) {
	if b.logger.V(2) {
		b.logger.Infof("Received SubConn state update: %p, %+v", subConn, state)
	}
	if b.subConn != subConn {
		if b.logger.V(2) {
			b.logger.Infof("Ignored state change because subConn is not recognized")
		}
		return
	}
	if state.ConnectivityState == connectivity.Shutdown {
		b.subConn = nil
		return
	}

	switch state.ConnectivityState {
	case connectivity.Ready:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{result: balancer.PickResult{SubConn: subConn}},
		})
	case connectivity.Connecting:
		if b.state == connectivity.TransientFailure {
			// We stay in TransientFailure until we are Ready. See A62.
			return
		}
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
		})
	case connectivity.Idle:
		if b.state == connectivity.TransientFailure {
			// We stay in TransientFailure until we are Ready. Also kick the
			// subConn out of Idle into Connecting. See A62.
			 //===========走到了这里，重新建立连接================
            b.subConn.Connect()
			return
		}
       
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &idlePicker{subConn: subConn},
		})
	case connectivity.TransientFailure:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{err: state.ConnectionError},
		})
	}
	b.state = state.ConnectivityState
}
```

整个调用链完成。当建立连接失败的时候会往chan中塞一个函数来重试。





## backoff机制

```go
func (bc Exponential) Backoff(retries int) time.Duration {
	if retries == 0 {
		return bc.Config.BaseDelay
	}
	backoff, max := float64(bc.Config.BaseDelay), float64(bc.Config.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.Config.Multiplier
		retries--
	}
	if backoff > max {
		backoff = max
	}
	// Randomize backoff delays so that if a cluster of requests start at
	// the same time, they won't operate in lockstep.
    //会随机加或减随机时间。
	backoff *= 1 + bc.Config.Jitter*(grpcrand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}

// DefaultConfig is a backoff configuration with the default values specfied
// at https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md.
//
// This should be useful for callers who want to configure backoff with
// non-default values only for a subset of the options.
var DefaultConfig = Config{
	BaseDelay:  1.0 * time.Second, //第一次失败等待多久
	Multiplier: 1.6, //乘以倍数
	Jitter:     0.2, //随机的范围。
	MaxDelay:   120 * time.Second, //最大等待时间
}
```



```go
//可以通过这个grpc.WithConnectParams参数设置。
conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithConnectParams(grpc.ConnectParams{
    Backoff: backoff.Config{
        BaseDelay:  1.0 * time.Second,
        Multiplier: 1.6,
        Jitter:     0.2,
        MaxDelay:   30 * time.Second,
    },
}))
if err != nil {
    log.Printf("DialContext failed %v", err)
    return
}
```