

## http详解

当你请求请求一次接口的时候go会开几个协程？

服务端keep-alive？

上传大文件的时候，会不会把大文件全部加载到内存中？

从http的源码分析。

### http从tcp中读出数据到执行我们的业务代码的流程

**1、go http默认的读缓存是4096，在读取一次之后，进行http各种相关信息的各种解析。**

w, err := c.readRequest(ctx)--> readRequest()-->tp.ReadLine()

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1713193587241-df356db1-ac95-49fc-b026-dcbd83a5aa5b.png)

在这一步会从tcp连接中读取4096个字节。进行一系列处理后给request，response赋值

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1713193688815-ff2408e7-669d-4c7c-a694-121e5622286b.png)



**2、执行我们的业务代码**

serverHandler{c.server}.ServeHTTP(w, w.req)如果在执行我们的业务代码的时候buf中的4096不足还是会去tcp连接中读取的。

**3、在一个tcp连接上一直循环这个过程。**

### 当你请求请求一次接口的时候go会开几个协程？

**先说结论：一次请求会开两个协程。**

```go
go install github.com/link1st/go-stress-testing@latest

mv $GOPATH/bin/go-stress-testing.exe  $GOPATH/bin/stress.exe

stress -c 1000 -n 1000 -u http://localhost:8888/
func TestHttpGoroutine(t *testing.T) {
    go func() {
        http.ListenAndServe("0.0.0.0:8899", nil)
    }()
    if err := http.ListenAndServe(":8888", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        select {
            case <-request.Context().Done():
            log.Println("request canceled")
        }
    })); err != nil {
        log.Fatal(err)
    }
}
```

 pprof 分析 http://192.168.2.109:8899/debug/pprof/ 可以看到当我们以1000个协程请求一次，对应的服务端会启动2000个协程，一次请求会启动两个协程。

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1713111021022-c3c65afa-8052-40c2-ae74-0b5ce98c2b9a.png)

```go
http.ListenAndServe-->server.ListenAndServe()-->srv.Serve(ln)-->rw, err := l.Accept()
-->go c.serve(connCtx)
for {
    rw, err := l.Accept()
    //拿到一个连接
    if err != nil {
        if srv.shuttingDown() {
            return ErrServerClosed
        }
        if ne, ok := err.(net.Error); ok && ne.Temporary() {
            if tempDelay == 0 {
                tempDelay = 5 * time.Millisecond
            } else {
                tempDelay *= 2
            }
            if max := 1 * time.Second; tempDelay > max {
                tempDelay = max
            }
            srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
            time.Sleep(tempDelay)
            continue
        }
        return err
    }
    connCtx := ctx
    if cc := srv.ConnContext; cc != nil {
        connCtx = cc(connCtx, rw)
        if connCtx == nil {
            panic("ConnContext returned nil")
        }
    }
    tempDelay = 0
    c := srv.newConn(rw)
    c.setState(c.rwc, StateNew, runHooks) // before Serve can return
    //开启协程处理
    go c.serve(connCtx)
}
serve(ctx context.Context)


	for {
		w, err := c.readRequest(ctx)
		if c.r.remain != c.server.initialReadLimitSize() {
			// If we read any bytes off the wire, we're active.
			c.setState(c.rwc, StateActive, runHooks)
		}

        //在读取完一次请求所有内容之后,，开启一个协程在后台读取。保证连接，当连接出现问题时候
        //及时通知。
		if requestBodyRemains(req.Body) {
			registerOnHitEOF(req.Body, w.conn.r.startBackgroundRead)
		} else {
			w.conn.r.startBackgroundRead()
		}
		serverHandler{c.server}.ServeHTTP(w, w.req)
		inFlightResponse = nil
		w.cancelCtx()
		if c.hijacked() {
			return
		}
		w.finishRequest()
		c.rwc.SetWriteDeadline(time.Time{})
		if !w.shouldReuseConnection() {
			if w.requestBodyLimitHit || w.closedRequestBodyEarly() {
				c.closeWriteAndWait()
			}
			return
		}
    }
```



```go
func (cr *connReader) startBackgroundRead() {
    cr.lock()
    defer cr.unlock()
    if cr.inRead {
        panic("invalid concurrent Body.Read call")
    }
    if cr.hasByte {
        return
    }
    cr.inRead = true
    cr.conn.rwc.SetReadDeadline(time.Time{})
    go cr.backgroundRead()
}

func (cr *connReader) backgroundRead() {
	n, err := cr.conn.rwc.Read(cr.byteBuf[:])
	cr.lock()
	if n == 1 {
		cr.hasByte = true
		// We were past the end of the previous request's body already
		// (since we wouldn't be in a background read otherwise), so
		// this is a pipelined HTTP request. Prior to Go 1.11 we used to
		// send on the CloseNotify channel and cancel the context here,
		// but the behavior was documented as only "may", and we only
		// did that because that's how CloseNotify accidentally behaved
		// in very early Go releases prior to context support. Once we
		// added context support, people used a Handler's
		// Request.Context() and passed it along. Having that context
		// cancel on pipelined HTTP requests caused problems.
		// Fortunately, almost nothing uses HTTP/1.x pipelining.
		// Unfortunately, apt-get does, or sometimes does.
		// New Go 1.11 behavior: don't fire CloseNotify or cancel
		// contexts on pipelined requests. Shouldn't affect people, but
		// fixes cases like Issue 23921. This does mean that a client
		// closing their TCP connection after sending a pipelined
		// request won't cancel the context, but we'll catch that on any
		// write failure (in checkConnErrorWriter.Write).
		// If the server never writes, yes, there are still contrived
		// server & client behaviors where this fails to ever cancel the
		// context, but that's kinda why HTTP/1.x pipelining died
		// anyway.
	}
   
	if ne, ok := err.(net.Error); ok && cr.aborted && ne.Timeout() {
         //当读取超时的时候走到这里。当serverHandler{c.server}.ServeHTTP(w, w.req)执行完
    //正常是这个调用链触发finishRequest()--> abortPendingRead()-->	cr.conn.rwc.SetReadDeadline(aLongTimeAgo)
		// Ignore this error. It's the expected error from
		// another goroutine calling abortPendingRead.
	} else if err != nil {
        //当出现读取错误的时候走到这里。出现比较多的情况是客户端取消，连接断开。
		cr.handleReadError(err)
	}
	cr.aborted = false
	cr.inRead = false
	cr.unlock()
	cr.cond.Broadcast()
}


// handleReadError is called whenever a Read from the client returns a
// non-nil error.
//
// The provided non-nil err is almost always io.EOF or a "use of
// closed network connection". In any case, the error is not
// particularly interesting, except perhaps for debugging during
// development. Any error means the connection is dead and we should
// down its context.
//
// It may be called from multiple goroutines.
func (cr *connReader) handleReadError(_ error) {
    //调用cancel通知结束。监听http的cancelCtx的协程会收到通知。
	cr.conn.cancelCtx()
	cr.closeNotify()
}
```

### 服务端keep-alive

go 服务端是默认开启http keep-alive的。`http.Server{}.SetKeepAlivesEnabled()`可以通过这个设置禁用

```go
for {
    w, err := c.readRequest(ctx)
    if c.r.remain != c.server.initialReadLimitSize() {
        // If we read any bytes off the wire, we're active.
        c.setState(c.rwc, StateActive, runHooks)
    }

    //可以开启这个才会一直在这个tcp连接中读取
    if !w.conn.server.doKeepAlives() {
        // We're in shutdown mode. We might've replied
        // to the user without "Connection: close" and
        // they might think they can send another
        // request, but such is life with HTTP/1.1.
        return
    }
}
```

### 上传大文件的时候，会不会把大文件全部加载到内存中？

**先说答案：不会，但是会创几个32MB的buf也是会占用一点内存。**

源码分析调用链

r.FormFile("uploadfile")

err := r.ParseMultipartForm(defaultMaxMemory)

f, err := mr.ReadForm(maxMemory)

readForm(maxMemory int64) 

p, err := r.nextPart(false, maxMemoryBytes, maxHeaders)

```go
package http

import (
    "crypto/md5"
    "fmt"
    "html/template"
    "io"
    "net/http"
    _ "net/http/pprof"
    "os"
    "strconv"
    "testing"
    "time"
)

func TestUpload(t *testing.T) {
    go func() {
        http.ListenAndServe("0.0.0.0:8899", nil)
    }()

    http.HandleFunc("/upload", upload)
    http.ListenAndServe(":8888", nil)
}

// 处理 /upload  逻辑
func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) // 获取请求的方法
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))
        t, _ := template.ParseFiles("upload.gtpl")
        t.Execute(w, token)
    } else {
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
        fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
    }
}
//将p可以理解为tcp连接中拷贝32MB到内存b中
n, err := io.CopyN(&b, p, maxFileMemoryBytes+1)
if err != nil && err != io.EOF {
    return nil, err
}
//如果文件大于32MB会创建一个临时文件来保存
	if n > maxFileMemoryBytes {
			if file == nil {
				file, err = os.CreateTemp(r.tempDir, "multipart-")
				if err != nil {
					return nil, err
				}
			}
			numDiskFiles++
			if _, err := file.Write(b.Bytes()); err != nil {
				return nil, err
			}
			if copyBuf == nil {
				copyBuf = make([]byte, 32*1024) // same buffer size as io.Copy uses
			}
			// os.File.ReadFrom will allocate its own copy buffer if we let io.Copy use it.
			type writerOnly struct{ io.Writer }
            //拷贝文件从tcp连接中拷贝数据到临时文件中，使用32MB的缓冲区。
			remainingSize, err := io.CopyBuffer(writerOnly{file}, p, copyBuf)
			if err != nil {
				return nil, err
			}
			fh.tmpfile = file.Name()
			fh.Size = int64(b.Len()) + remainingSize
			fh.tmpoff = fileOff
			fileOff += fh.Size
			if !combineFiles {
				if err := file.Close(); err != nil {
					return nil, err
				}
				file = nil
			}
		} else {
			fh.content = b.Bytes()
			fh.Size = int64(len(fh.content))
			maxFileMemoryBytes -= n
			maxMemoryBytes -= n
		}
```

文件大于32MB

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1713196945179-f232492f-0f10-4139-b75d-895a79465a9c.png)

文件小于32MB

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1713196992970-2167c3aa-b59f-4c6c-9d31-581729f478f8.png)