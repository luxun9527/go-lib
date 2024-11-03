# Redis

### 目标

1. redis运维相关

1. 1. redis安装
   2. redis常用配置常用配置。
   3. redis集群

1. redis常用命令。
2. go sdk
3. 面试的常见问题。





### redis运维相关

https://www.cnblogs.com/yidengjiagou/p/17345831.html

https://gitee.com/zhengqingya/docker-compose/tree/master/Linux/redis/redis7.0.5

#### `单节点`

```yaml
version: '3'
services:
  redis:
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/redis:7.0.5                    # 镜像'redis:7.0.5'
    container_name: redis                                                             # 容器名为'redis'
    restart: unless-stopped                                                                   # 指定容器退出后的重启策略为始终重启，但是不考虑在Docker守护进程启动时就已经停止了的容器
    command: redis-server /etc/redis/redis.conf --requirepass 123456 --appendonly yes # 启动redis服务并添加密码为：123456,默认不开启redis-aof方式持久化配置
#    command: redis-server --requirepass 123456 --appendonly yes # 启动redis服务并添加密码为：123456,并开启redis持久化配置
    environment:                        # 设置环境变量,相当于docker run命令中的-e
      TZ: Asia/Shanghai
      LANG: en_US.UTF-8
    volumes:                            # 数据卷挂载路径设置,将本机目录映射到容器目录
      - "./redis/data:/data"
      - "./redis/config/redis.conf:/etc/redis/redis.conf"  # `redis.conf`文件内容`http://download.redis.io/redis-stable/redis.conf`
    ports:                              # 映射端口
      - "6379:6379"
```

#### 常用配置

```toml
#append only配置
appendonly yes

appendfilename "appendonly.aof"

# For convenience, Redis stores all persistent append-only files in a dedicated
# directory. The name of the directory is determined by the appenddirname
# configuration parameter.

appenddirname "appendonlydir"
# appendfsync always
appendfsync everysec

# Unless specified otherwise, by default Redis will save the DB:
#   * After 3600 seconds (an hour) if at least 1 change was performed
#   * After 300 seconds (5 minutes) if at least 100 changes were performed
#   * After 60 seconds if at least 10000 changes were performed
#
# You can set these explicitly by uncommenting the following line.
# rdb存储方式
save 3600 1 300 100 60 10000
# 密码
requirepass 123456
# 数据存储的位置
dir /data
#rdb+aof混合存储
aof-rewrite-incremental-fsync yes
```

#### 集群

https://github.com/luxun9527/devops/tree/main/middleware/redis/sharding-cluster

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1730040753430-69203593-a74a-4a5d-9659-95e1f215023d.png)

https://gitee.com/zhengqingya/docker-compose/tree/master/Linux/redis/redis7.0.5

https://blog.csdn.net/Go_Bin/article/details/136322685

常用的集群模式，分片+多副本。

```yaml
version: '3'
services:

  redis-node1:  # 服务名
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/redis:7.0.5                   # 镜像'redis:7.0.5'
    container_name: redis-node1 # docker启动的容器名
    ports:  # 映射的端口 7001是redis server使用，17001是集群之间节点通信使用，都必须开放映射，如果不指定17001端口映射的话，创建集群的时候节点之间不能通信，集群会创建失败
      - "7001:7001"
      - "17001:17001"
    volumes:  # 映射的容器卷
      - "./redis-cluster/redis-6381/config/redis.conf:/etc/redis/redis.conf"  # `redis.conf`文件内容`http://download.redis.io/redis-stable/redis.conf`
      - "./redis-cluster/redis-6381/data:/data"
    command: redis-server /etc/redis/redis.conf
    networks:  # 指定使用网络插件名称
      - redis-net

  redis-node2:
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/redis:7.0.5                   # 镜像'redis:7.0.5'
    container_name: redis-node2
    ports:
      - "7002:7002"
      - "17002:17002"
    volumes:
      - "./redis-cluster/redis-6382/config/redis.conf:/etc/redis/redis.conf"  # `redis.conf`文件内容`http://download.redis.io/redis-stable/redis.conf`
      - "./redis-cluster/redis-6382/data:/data"
    command: redis-server /etc/redis/redis.conf
    networks:
      - redis-net

  redis-node3:
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/redis:7.0.5                   # 镜像'redis:7.0.5'
    container_name: redis-node3
    ports:
      - "7003:7003"
      - "17003:17003"
    volumes:
      - "./redis-cluster/redis-6383/config/redis.conf:/etc/redis/redis.conf"  # `redis.conf`文件内容`http://download.redis.io/redis-stable/redis.conf`
      - "./redis-cluster/redis-6383/data:/data"
    command: redis-server /etc/redis/redis.conf
    networks:
      - redis-net

  redis-node4:
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/redis:7.0.5                   # 镜像'redis:7.0.5'
    container_name: redis-node4
    ports:
      - "7004:7004"
      - "17004:17004"
    volumes:
      - "./redis-cluster/redis-6384/config/redis.conf:/etc/redis/redis.conf"  # `redis.conf`文件内容`http://download.redis.io/redis-stable/redis.conf`
      - "./redis-cluster/redis-6384/data:/data"
    command: redis-server /etc/redis/redis.conf
    networks:
      - redis-net

  redis-node5:
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/redis:7.0.5                   # 镜像'redis:7.0.5'
    container_name: redis-node5
    ports:
      - "7005:7005"
      - "17005:17005"
    volumes:
      - "./redis-cluster/redis-6385/config/redis.conf:/etc/redis/redis.conf"  # `redis.conf`文件内容`http://download.redis.io/redis-stable/redis.conf`
      - "./redis-cluster/redis-6385/data:/data"
    command: redis-server /etc/redis/redis.conf
    networks:
      - redis-net

  redis-node6:
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/redis:7.0.5                   # 镜像'redis:7.0.5'
    container_name: redis-node6
    ports:
      - "7006:7006"
      - "17006:17006"
    volumes:
      - "./redis-cluster/redis-6386/config/redis.conf:/etc/redis/redis.conf"  # `redis.conf`文件内容`http://download.redis.io/redis-stable/redis.conf`
      - "./redis-cluster/redis-6386/data:/data"
    command: redis-server /etc/redis/redis.conf
    networks:
      - redis-net



networks:
  redis-net:
```

#### `常用命令`



##### 登录

```
redis-cli -c -h redis-6381 -p 6381 -a 123456
```

-c 指定为集群模式

-h 指定host

##### 创建集群

```
docker exec -it redis-node4 redis-cli -h 192.168.2.159 -p 7001  -a bingo --cluster create 192.168.2.159:7001 192.168.2.159:7002 192.168.2.159:7003 192.168.2.159:7004 192.168.2.159:7005 192.168.2.159:7006 --cluster-replicas 1
```



##### 查找key

```
KEYS user:*
KEYS *foo*
get key
```

\# 查看集群信息

```
cluster info
```

\# 查看集群节点信息

```
cluster nodes
```

\# 查看slots分片

```
cluster slots
```

### go sdk

#### 基础用法

```go
package main

import (
    "context"
    "github.com/redis/go-redis/extra/redisotel/v9"
    "github.com/redis/go-redis/v9"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    tracesdk "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
    "log"
    "sync"

    "testing"

    "time"
)

// v9 对应的服务端是7.0 v8对应的是是6.0
func TestRedis(t *testing.T) {
    ctx := context.Background()

    // 创建 Redis 集群客户端
    rdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "192.168.2.159:7001", // 分片节点1的地址
            "192.168.2.159:7002", // 分片节点2的地址
            "192.168.2.159:7003", // 分片节点3的地址
        },
        Password: "bingo", // 如果设置了密码，可以在这里配置
    })

    // 测试连接
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Error connecting to Redis cluster: %v", err)
    }
    log.Printf("Connected to Redis cluster: %s", pong)

    // 使用示例：设置和获取一个 key
    err = rdb.Set(ctx, "example_key", "Hello, Redis!", 0).Err()
    if err != nil {
        log.Fatalf("Error setting key: %v", err)
    }

    val, err := rdb.Get(ctx, "example_key").Result()
    if err != nil {
        log.Fatalf("Error getting key: %v", err)
    }

    log.Printf("Value of 'example_key': %s", val)

    // 关闭客户端
    if err := rdb.Close(); err != nil {
        log.Printf("Error closing client: %v", err)
    }
}
```



#### redis链路追踪

```go
package main

import (
    "context"
    "github.com/redis/go-redis/extra/redisotel/v9"
    "github.com/redis/go-redis/v9"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    tracesdk "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
    "log"
    "sync"

    "testing"

    "time"
)

const (
    service     = "trace-demo" // 服务名
    environment = "production" // 环境
    id          = 1            // id
)

const (
    serviceName    = "redis-Jaeger-Demo"
    jaegerEndpoint = "192.168.2.159:14268/api/traces"
)

var tracer = otel.Tracer("redis-demo")

// newJaegerTraceProvider 创建一个 Jaeger Trace Provider

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {

    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
    if err != nil {
        return nil, err
    }
    tp := tracesdk.NewTracerProvider(
        // Always be sure to batch in production.
        tracesdk.WithBatcher(exp),
        // Record information about this application in a Resource.
        tracesdk.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(service),
            attribute.String("environment", environment),
            attribute.Int64("ID", id),
        ),
                            ),
        tracesdk.WithSampler(tracesdk.AlwaysSample()),
    )
    return tp, nil
}

// initTracer 初始化 Tracer
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
    tp, err := tracerProvider("http://192.168.2.159:14268/api/traces")
    if err != nil {
        return nil, err
    }

    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(
        propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
    )
    return tp, nil
}
func doSomething(ctx context.Context, rdb *redis.ClusterClient) error {
    if err := rdb.Set(ctx, "name", "Q1mi", time.Minute).Err(); err != nil {
        return err
    }
    if err := rdb.Set(ctx, "tag", "OTel", time.Minute).Err(); err != nil {
        return err
    }
    var wg sync.WaitGroup
    for range []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5} {
        wg.Add(1)
        go func() {
            defer wg.Done()
            val := rdb.Get(ctx, "tag").Val()
            if val != "OTel" {
                log.Printf("%q != %q", val, "OTel")
            }
        }()
    }
    wg.Wait()

    if err := rdb.Del(ctx, "name").Err(); err != nil {
        return err
    }
    if err := rdb.Del(ctx, "tag").Err(); err != nil {
        return err
    }
    log.Println("done!")
    return nil
}
// 链路追踪https://www.liwenzhou.com/posts/Go/redis-otel/

func TestRedisTrace(t *testing.T) {

    ctx := context.Background()

    tp, err := initTracer(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := tp.Shutdown(ctx); err != nil {
            log.Printf("Error shutting down tracer provider: %v", err)
        }
    }()

    // 创建 Redis 集群客户端
    rdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "192.168.2.159:7001", // 分片节点1的地址
            "192.168.2.159:7002", // 分片节点2的地址
            "192.168.2.159:7003", // 分片节点3的地址
        },
        Password: "bingo", // 如果设置了密码，可以在这里配置
    })

    // 启用 tracing
    if err := redisotel.InstrumentTracing(rdb); err != nil {
        panic(err)
    }

    // 启用 metrics
    if err := redisotel.InstrumentMetrics(rdb); err != nil {
        panic(err)
    }

    ctx, span := tracer.Start(ctx, "doSomething")
    defer span.End()

    if err := doSomething(ctx, rdb); err != nil {
        span.RecordError(err) // 记录error
        span.SetStatus(codes.Error, err.Error())
    }
}
```



#### redis分布式锁

https://github.com/zeromicro/go-zero/blob/master/core/stores/redis/redislock.go



### 面试

https://xiaolincoding.com/redis/