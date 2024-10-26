### 目标

1. 核心概念
2. kafka运维相关

1. 1. kafka安装
   2. kafka常用配置常用配置。
   3. kafka集群

1. kafka常用命令。
2. go sdk
3. 面试的常见问题。



代码地址：https://github.com/luxun9527/go-lib/tree/master/sdk/kafka

docker-compose 文件地址：https://github.com/luxun9527/devops/tree/main/middleware/kafka

您的点赞，star，评论都是对我的鼓励

### kafka概念

kafka架构

[https://www.cnblogs.com/hgzero/p/17229564.html#%E7%94%9F%E4%BA%A7%E5%92%8C%E6%B6%88%E8%B4%B9%E5%91%BD%E4%BB%A4](https://www.cnblogs.com/hgzero/p/17229564.html#生产和消费命令)

这个介绍真的太全面了相关概念参考即可。

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1729957825553-cc8061ce-6539-4828-a161-0581dd18d2d1.png)

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1729524628386-e54180bb-0d1c-4fc2-887c-78eb2d75e2d3.png)





### 使用docker-compose安装kafka

https://gitee.com/zhengqingya/docker-compose/blob/master/Linux/kafka/3.4.1/docker-compose-kafka.yml#

https://program.snlcw.com/1009.html

https://www.cnblogs.com/hovin/p/18186821

单节点安装。

各个配置项参考
https://github.com/bitnami/containers/tree/main/bitnami/kafka#configuration

https://github.com/bitnami/containers/tree/main/bitnami/kafka#111-debian-9-r205-220-debian-9-r40-111-ol-7-r286-and-220-ol-7-r53



https://blog.csdn.net/u011618288/article/details/129105777

bitnami/kafka这个镜像都是可以通过环境变量的形式来配置，你也可以首次启动后将配置拷贝出来配置在/opt/bitnami/kafka/config目录下

启用了SASL_PLAINTEXT/SCRAM-SHA-256认证

**单节点**

使用docker-compose来安装。

```yaml
version: '3'

# 网桥 -> 方便相互通讯
networks:
  kafka:
    ipam:
      driver: default
      config:
        - subnet: "172.22.6.0/24"

services:
  zookepper:
    image: docker.1panel.dev/bitnami/zookeeper:3.8                   # 原镜像`bitnami/zookeeper:latest`
    container_name: zookeeper                        # 容器名为'zookeeper-server'
    restart: unless-stopped                                  # 指定容器退出后的重启策略为始终重启，但是不考虑在Docker守护进程启动时就已经停止了的容器
    volumes:                                         # 数据卷挂载路径设置,将本机目录映射到容器目录
      - "/etc/localtime:/etc/localtime"
     # - "./kafka/zookeeper/conf:/opt/bitnami/zookeeper/conf"
      - "./kafka/zookeeper/data:/bitnami/zookeeper/data"
    environment:
       ZOO_ENABLE_AUTH: yes
       ZOO_SERVER_USERS: user
       ZOO_SERVER_PASSWORDS: pass123
       ZOO_CLIENT_USER: user
       ZOO_CLIENT_PASSWORD: pass123
    networks:
      kafka:
        ipv4_address: 172.22.6.12
  kafka:
    image: 'docker.1panel.dev/bitnami/kafka:2.8.1'
    ports:
      - '9093:9093'
      - '9092:9092'
    networks:
      kafka:
        ipv4_address: 172.22.6.11                                   # 容器名为'kafka'
    restart: unless-stopped                                          # 指定容器退出后的重启策略为始终重启，但是不考虑在Docker守护进程启动时就已经停止了的容器
    volumes:                                                 # 数据卷挂载路径设置,将本机目录映射到容器目录
      - "/etc/localtime:/etc/localtime"
      - "./kafka/kafka/config:/opt/bitnami/kafka/config"
      - "./kafka/kafka/data:/bitnami/kafka/data"
    environment:
      # 监听器的 CLIENT 上不要设置 9092，否则日志一直输出下面信息：
      # INFO [SocketServer listenerType=ZK_BROKER, nodeId=1] Failed authentication with /10.0.0.2 (SSL handshake failed)
       KAFKA_CFG_BROKER_ID: 12
       KAFKA_CFG_LISTENERS: INTERNAL://0.0.0.0:9092,CLIENT://0.0.0.0:9093 #监听的地址
       KAFKA_CFG_ADVERTISED_LISTENERS: INTERNAL://192.168.2.159:9092,CLIENT://192.168.2.159:9093 #相对其他BROKER的地址
       KAFKA_CFG_INTER_BROKER_LISTENER_NAME: INTERNAL  #监听的名字
       KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: INTERNAL:SASL_PLAINTEXT,CLIENT:SASL_PLAINTEXT #可用的认证方式
      # 用于 broker 间通信的 SASL 机制
       KAFKA_CFG_SASL_MECHANISM_INTER_BROKER_PROTOCOL: SCRAM-SHA-512 #PLAIN,SCRAM-SHA-256,SCRAM-SHA-512可选 
       KAFKA_CFG_SASL_MECHANISM_CONTROLLER_PROTOCOL: SCRAM-SHA-512 # PLAIN,SCRAM-SHA-256,SCRAM-SHA-512可选
      # 允许使用明文监听，出于安全原因，Bitnami Apache Kafka docker 镜像禁用了 PLAINTEXT 侦听器，但可以通过下面方式开启
       ALLOW_PLAINTEXT_LISTENER: no
       KAFKA_CFG_ZOOKEEPER_CONNECT: zookeeper:2181
       KAFKA_ZOOKEEPER_PROTOCOL: SASL     # 可选值有：PLAINTEXT, SASL, SSL, and SASL_SSL. 默认值: PLAINTEXT
       KAFKA_ZOOKEEPER_USER: user
       KAFKA_ZOOKEEPER_PASSWORD: pass123
      #Inter broker credentials
       KAFKA_INTER_BROKER_USER: interbrokeruser        # Kafka 内部节点通信的用户名，默认值：user
       KAFKA_INTER_BROKER_PASSWORD: interbrokerpass    # Kafka 内部节点通信的密码，默认值：bitnami
      #Client credentials（配置 SASL 认证时，使用下面两个变量来配置用户名和密码）
       KAFKA_CLIENT_USERS: clientuser1       # 使用 SASL 模式处理客户端通信时创建的用户名，用逗号隔开。
       KAFKA_CLIENT_PASSWORDS: pass123       # 使用 SASL 模式处理客户端通信时创建的密码，用逗号隔开。

    depends_on:
      - zookepper

  kafka-map:
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/kafka-map                         # 原镜像`dushixiang/kafka-map:latest`
    container_name: kafka-map                            # 容器名为'kafka-map'
    restart: unless-stopped                                          # 指定容器退出后的重启策略为始终重启，但是不考虑在Docker守护进程启动时就已经停止了的容器
    volumes:
      - "./kafka/kafka-map/data:/usr/local/kafka-map/data"
    environment:
      DEFAULT_USERNAME: admin
      DEFAULT_PASSWORD: 123456
    ports: # 映射端口
      - "9006:8080"
    depends_on: # 解决容器依赖启动先后问题
      - kafka
    networks:
      kafka:
        ipv4_address: 172.22.6.13
```

https://github.com/luxun9527/devops/blob/main/middleware/kafka/kafka/kafka/config/server.properties 



配置文件为认证相关配置server.properties 配置如下。

```toml
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# see kafka.server.KafkaConfig for additional details and defaults

############################# Server Basics #############################

# The id of the broker. This must be set to a unique integer for each broker.
broker.id=12

############################# Socket Server Settings #############################

# The address the socket server listens on. It will get the value returned from 
# java.net.InetAddress.getCanonicalHostName() if not configured.
#   FORMAT:
#     listeners = listener_name://host_name:port
#   EXAMPLE:
#     listeners = PLAINTEXT://your.host.name:9092
listeners=INTERNAL://0.0.0.0:9092,CLIENT://0.0.0.0:9093

# Hostname and port the broker will advertise to producers and consumers. If not set, 
# it uses the value for "listeners" if configured.  Otherwise, it will use the value
# returned from java.net.InetAddress.getCanonicalHostName().
advertised.listeners=INTERNAL://192.168.2.159:9092,CLIENT://192.168.2.159:9093

# Maps listener names to security protocols, the default is for them to be the same. See the config documentation for more details
listener.security.protocol.map=INTERNAL:SASL_PLAINTEXT,CLIENT:SASL_PLAINTEXT

# The number of threads that the server uses for receiving requests from the network and sending responses to the network
num.network.threads=3

# The number of threads that the server uses for processing requests, which may include disk I/O
num.io.threads=8

# The send buffer (SO_SNDBUF) used by the socket server
socket.send.buffer.bytes=102400

# The receive buffer (SO_RCVBUF) used by the socket server
socket.receive.buffer.bytes=102400

# The maximum size of a request that the socket server will accept (protection against OOM)
socket.request.max.bytes=104857600


############################# Log Basics #############################

# A comma separated list of directories under which to store log files
log.dirs=/bitnami/kafka/data

# The default number of log partitions per topic. More partitions allow greater
# parallelism for consumption, but this will also result in more files across
# the brokers.
num.partitions=1

# The number of threads per data directory to be used for log recovery at startup and flushing at shutdown.
# This value is recommended to be increased for installations with data dirs located in RAID array.
num.recovery.threads.per.data.dir=1

############################# Internal Topic Settings  #############################
# The replication factor for the group metadata internal topics "__consumer_offsets" and "__transaction_state"
# For anything other than development testing, a value greater than 1 is recommended to ensure availability such as 3.
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1

############################# Log Flush Policy #############################

# Messages are immediately written to the filesystem but by default we only fsync() to sync
# the OS cache lazily. The following configurations control the flush of data to disk.
# There are a few important trade-offs here:
#    1. Durability: Unflushed data may be lost if you are not using replication.
#    2. Latency: Very large flush intervals may lead to latency spikes when the flush does occur as there will be a lot of data to flush.
#    3. Throughput: The flush is generally the most expensive operation, and a small flush interval may lead to excessive seeks.
# The settings below allow one to configure the flush policy to flush data after a period of time or
# every N messages (or both). This can be done globally and overridden on a per-topic basis.

# The number of messages to accept before forcing a flush of data to disk
#log.flush.interval.messages=10000

# The maximum amount of time a message can sit in a log before we force a flush
#log.flush.interval.ms=1000

############################# Log Retention Policy #############################

# The following configurations control the disposal of log segments. The policy can
# be set to delete segments after a period of time, or after a given size has accumulated.
# A segment will be deleted whenever *either* of these criteria are met. Deletion always happens
# from the end of the log.

# The minimum age of a log file to be eligible for deletion due to age
log.retention.hours=168

# A size-based retention policy for logs. Segments are pruned from the log unless the remaining
# segments drop below log.retention.bytes. Functions independently of log.retention.hours.
#log.retention.bytes=1073741824

# The maximum size of a log segment file. When this size is reached a new log segment will be created.
log.segment.bytes=1073741824

# The interval at which log segments are checked to see if they can be deleted according
# to the retention policies
log.retention.check.interval.ms=600000

############################# Zookeeper #############################

# Zookeeper connection string (see zookeeper docs for details).
# This is a comma separated host:port pairs, each corresponding to a zk
# server. e.g. "127.0.0.1:3000,127.0.0.1:3001,127.0.0.1:3002".
# You can also append an optional chroot string to the urls to specify the
# root directory for all kafka znodes.
zookeeper.connect=zookeeper:2181

# Timeout in ms for connecting to zookeeper
zookeeper.connection.timeout.ms=18000


############################# Group Coordinator Settings #############################

# The following configuration specifies the time, in milliseconds, that the GroupCoordinator will delay the initial consumer rebalance.
# The rebalance will be further delayed by the value of group.initial.rebalance.delay.ms as new members join the group, up to a maximum of max.poll.interval.ms.
# The default value for this is 3 seconds.
# We override this to 0 here as it makes for a better out-of-the-box experience for development and testing.
# However, in production environments the default value of 3 seconds is more suitable as this will help to avoid unnecessary, and potentially expensive, rebalances during application startup.
group.initial.rebalance.delay.ms=0

auto.create.topics.enable=true

inter.broker.listener.name=INTERNAL

max.partition.fetch.bytes=1048576
max.request.size=1048576
sasl.enabled.mechanisms=PLAIN,SCRAM-SHA-256,SCRAM-SHA-512
sasl.mechanism.inter.broker.protocol=SCRAM-SHA-512

sasl.mechanism.controller.protocol=SCRAM-SHA-512
```

**kafka_jaas.conf**

```plain
KafkaClient {
   org.apache.kafka.common.security.plain.PlainLoginModule required
   username="clientuser1"
   password="pass123";
   };
KafkaServer {
   org.apache.kafka.common.security.scram.ScramLoginModule required
   username="interbrokeruser"
   password="interbrokerpass";
   };
Client {
   org.apache.kafka.common.security.plain.PlainLoginModule required
   username="user"
   password="pass123";
   };
```

**集群安装**

**不用太多配置，直接连接上zookepper即可。**

```yaml
# https://hub.docker.com/r/bitnami/kafka
# https://github.com/bitnami/containers/blob/main/bitnami/kafka/docker-compose-cluster.yml

version: '3'

# 定义通用配置
x-kafka-common: &kafka-common
  image: 'docker.1panel.dev/bitnami/kafka:2.8.1'
  restart: unless-stopped                                          # 指定容器退出后的重启策略为始终重启，但是不考虑在Docker守护进程启动时就已经停止了的容器
  depends_on:
    - zookepper
  links:
    - zookepper
x-kafka-common-env: &kafka-common-env
  KAFKA_CFG_INTER_BROKER_LISTENER_NAME: INTERNAL
  KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: INTERNAL:SASL_PLAINTEXT,CLIENT:SASL_PLAINTEXT
  # 用于 broker 间通信的 SASL 机制
  KAFKA_CFG_SASL_MECHANISM_INTER_BROKER_PROTOCOL: PLAIN
  # 允许使用明文监听，出于安全原因，Bitnami Apache Kafka docker 镜像禁用了 PLAINTEXT 侦听器，但可以通过下面方式开启
  ALLOW_PLAINTEXT_LISTENER: no
  KAFKA_CFG_ZOOKEEPER_CONNECT: zookeeper:2181
  KAFKA_ZOOKEEPER_PROTOCOL: SASL     # 可选值有：PLAINTEXT, SASL, SSL, and SASL_SSL. 默认值: PLAINTEXT
  KAFKA_ZOOKEEPER_USER: user
  KAFKA_ZOOKEEPER_PASSWORD: pass123
  #Inter broker credentials
  KAFKA_INTER_BROKER_USER: interbrokeruser        # Kafka 内部节点通信的用户名，默认值：user
  KAFKA_INTER_BROKER_PASSWORD: interbrokerpass    # Kafka 内部节点通信的密码，默认值：bitnami
  #Client credentials（配置 SASL 认证时，使用下面两个变量来配置用户名和密码）
  KAFKA_CLIENT_USERS: clientuser1       # 使用 SASL 模式处理客户端通信时创建的用户名，用逗号隔开。
  KAFKA_CLIENT_PASSWORDS: pass123       # 使用 SASL 模式处理客户端通信时创建的密码，用逗号隔开。
  KAFKA_CFG_OFFSETS_TOPIC_REPLICATION_FACTOR: 2 # 消费组提交的偏移量信息在几个 broker 存储，默认值：1
  KAFKA_CFG_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 2 # 事务状态 broker 存储，默认值：1
  KAFKA_CFG_DEFAULT_REPLICATION_FACTOR: 2 # 默认topic数据有几个副本。
# 网桥 -> 方便相互通讯
networks:
  kafka:
    ipam:
      driver: default
      config:
        - subnet: "172.22.7.0/24"

services:
  zookepper:
    image: docker.1panel.dev/bitnami/zookeeper:3.8                   # 原镜像`bitnami/zookeeper:latest`
    container_name: zookeeper                        # 容器名为'zookeeper-server'
    restart: unless-stopped                                  # 指定容器退出后的重启策略为始终重启，但是不考虑在Docker守护进程启动时就已经停止了的容器
    volumes:                                         # 数据卷挂载路径设置,将本机目录映射到容器目录
      - "/etc/localtime:/etc/localtime"
      # - "./kafka/zookeeper/conf:/opt/bitnami/zookeeper/conf"
      - "./kafka/zookeeper/data:/bitnami/zookeeper/data"
    environment:
      ZOO_ENABLE_AUTH: yes
      ZOO_SERVER_USERS: user
      ZOO_SERVER_PASSWORDS: pass123
      ZOO_CLIENT_USER: user
      ZOO_CLIENT_PASSWORD: pass123
    networks:
      kafka:
        ipv4_address: 172.22.7.20

  kafka-1:
    container_name: kafka-1
    <<: *kafka-common
    volumes:
      - "/etc/localtime:/etc/localtime"
      - "./kafka/kafka1/data:/bitnami/kafka/data"
    environment:
      <<: *kafka-common-env
      KAFKA_CFG_BROKER_ID: 1
      KAFKA_CFG_LISTENERS: INTERNAL://0.0.0.0:19092,CLIENT://0.0.0.0:19093
      KAFKA_CFG_ADVERTISED_LISTENERS: INTERNAL://192.168.2.159:19092,CLIENT://192.168.2.159:19093
    ports:
    - "19092:19092"
    - "19093:19093"
    networks:
      kafka:
        ipv4_address: 172.22.7.21
  kafka-2:
    container_name: kafka-2
    <<: *kafka-common
    volumes:
      - "/etc/localtime:/etc/localtime"
      - "./kafka/kafka2/data:/bitnami/kafka/data"
    environment:
      <<: *kafka-common-env
      KAFKA_CFG_BROKER_ID: 2
      KAFKA_CFG_LISTENERS: INTERNAL://0.0.0.0:19094,CLIENT://0.0.0.0:19095
      KAFKA_CFG_ADVERTISED_LISTENERS: INTERNAL://192.168.2.159:19094,CLIENT://192.168.2.159:19095
    ports:
      - "19094:19094"
      - "19095:19095"
    networks:
      kafka:
        ipv4_address: 172.22.7.22

  # kafka-map图形化管理工具
  kafka-map:
    image: registry.cn-hangzhou.aliyuncs.com/zhengqing/kafka-map     # 原镜像`dushixiang/kafka-map:latest`
    container_name: kafka-map                                        # 容器名为'kafka-map'
    restart: unless-stopped                                          # 指定容器退出后的重启策略为始终重启，但是不考虑在Docker守护进程启动时就已经停止了的容器
    volumes:
      - "./kafka/kafka-map/data:/usr/local/kafka-map/data"
    environment:
      DEFAULT_USERNAME: admin
      DEFAULT_PASSWORD: 123456
    ports:                              # 映射端口
      - "9006:8080"
    depends_on:                         # 解决容器依赖启动先后问题
      - kafka-1
      - kafka-2
    links:                              # 配置容器互相连接
      - kafka-1
      - kafka-2
    networks:
      kafka:
        ipv4_address: 172.22.7.24
```

### kafka常用命令

当有配置加密的时候，使用kafka自带的命令。要修改配置。



export KAFKA_OPTS="-Djava.security.auth.login.config=/opt/bitnami/kafka/config/kafka_jaas.conf"

```toml
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
# 
#    http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# see org.apache.kafka.clients.consumer.ConsumerConfig for more details

# list of brokers used for bootstrapping knowledge about the rest of the cluster
# format: host1:port1,host2:port2 ...
bootstrap.servers=localhost:9092

# consumer group id
group.id=test-consumer-group

# What to do when there is no initial offset in Kafka or if the current
# offset does not exist any more on the server: latest, earliest, none
#auto.offset.reset=

max.partition.fetch.bytes=1048576

#认证相关。
security.protocol=SASL_PLAINTEXT

sasl.mechanism=PLAIN
#sasl.mechanism=SCRAM-SHA-512
sasl.jaas.config= org.apache.kafka.common.security.plain.PlainLoginModule  required username="clientuser1" password="pass123";
#SCRAM-SHA-512设置
#sasl.jaas.config=org.apache.kafka.common.security.scram.ScramLoginModule required username="clientuser1" password="pass123";

#关闭消费自动提交。
enable.auto.commit=false
```

#### 查看消费组的消费情况

```
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group test-consumer-group --command-config=/bitnami/kafka/data/consumer.properties
```

#### 获取topic中的数据

要关闭自动提交，否则会提交offset

```
kafka-console-consumer.sh  --bootstrap-server localhost:9092 --topic your_topic --from-beginning  --consumer.config=/bitnami/kafka/data/consumer.properties
```

#### 消费数据

```
kafka-console-consumer.sh  --bootstrap-server localhost:9092 --topic your_topic --from-beginning  --consumer.config=/bitnami/kafka/data/consumer.properties
```

#### kafka发送数据

```
 kafka-console-producer.sh --broker-list localhost:9092  --producer.config=/bitnami/kafka/data/producer.properties  --topic your_topic  
```

#### 创建topic

```
**kafka-topics.sh --create --bootstrap-server localhost:19092 --replication-factor 2 --partitions 4 --topic testtopic** --command-config=/bitnami/kafka/data/consumer.properties
```

默认情况下，Kafka 集群的每个 topic 的副本数由 `**replication-factor**` 参数决定。如果您在创建 topic 时未指定此参数，Kafka 默认将其设置为 1。这意味着每个消息仅会在一个节点上保存，不会有副本。

如果您希望在 Kafka 集群中实现数据冗余和容错，您需要在创建 topic 时设置 `replication-factor` 为 2（或者更高）。这会导致 Kafka 将每个消息复制到多个节点上，从而提供更高的数据安全性。

**删除topic**

```
kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic testtopic --command-config=/bitnami/kafka/data/consumer.properties
```

#### 获取所有topic

```
kafka-topics.sh   --bootstrap-server localhost:9092 --list   --command-config=/bitnami/kafka/data/consumer.properties
```

### go sdk

#### 生产者

```go
package main

import (
    "crypto/sha256"
    "crypto/sha512"
    "encoding/json"
    "fmt"
    "github.com/Shopify/sarama"
    "github.com/luxun9527/zlog"
    "github.com/xdg-go/scram"
    "github.com/zeromicro/go-zero/core/stringx"
    "log"
    "testing"
)

type User struct {
    ID     string `json:"id"`
    Detail []byte `json:"detail"`
}

func TestProduce1(t *testing.T) {

    sarama.Logger = zlog.KafkaSaramaLogger
    // Kafka broker 地址
    brokers := []string{"192.168.2.159:9092"}
    // Kafka 主题
    topic := "test-topic"
    // 创建一个新的配置
    config := sarama.NewConfig()

    // 根据需求设置消息确认模式
    // "acks=0" - 生产者发送消息后立即返回
    //config.Producer.RequiredAcks = sarama.NoResponse
    // "acks=1" - 仅等待 Leader 确认
    //config.Producer.RequiredAcks = sarama.WaitForLocal
    // "acks=all" - 等待所有副本确认
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Return.Successes = true
    // 创建一个新的同步生产者
    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        log.Fatalf("Failed to create producer: %v", err)
    }
    defer producer.Close()
    u := &User{
        ID:     stringx.Randn(29),
        Detail: []byte("123456789"),
    }
    d, _ := json.Marshal(u)
    // 要发送的消息
    message := &sarama.ProducerMessage{
        Topic:     topic,
        Value:     sarama.ByteEncoder(d),
        Partition: 0, //指定分区，可选
        Key: nil,//如果指定了key，则会根据可以选择一个分区。相同的key 会被发送到同一个分区
    }

    // 发送消息
    partition, offset, err := producer.SendMessage(message)
    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }

    fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
}

func TestSaslPlainSCRAMSHA512(t *testing.T) {
    sarama.Logger = zlog.KafkaSaramaLogger

    // 创建 Kafka 配置
    config := sarama.NewConfig()
    config.Version = sarama.V2_8_0_0 // 根据你的 Kafka 版本设置
    config.Net.SASL.Enable = true
    config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
    config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
    config.Net.SASL.User = "clientuser1" // 替换为你的用户名
    config.Net.SASL.Password = "pass123" // 替换为你的密码
    config.Net.SASL.Enable = true
    config.Producer.Return.Successes = true

    // 创建生产者
    //producer, err := sarama.NewSyncProducer([]string{"192.168.2.159:19092", "192.168.2.159:19094"}, config) // 替换为你的 Kafka broker
    producer, err := sarama.NewSyncProducer([]string{"192.168.2.159:9092"}, config) // 替换为你的 Kafka broker
    if err != nil {
        log.Fatalf("Failed to create producer: %s", err)
    }
    defer producer.Close()

    // 发送消息
    message := &sarama.ProducerMessage{
        Topic: "your_topic", // 替换为你的主题
        Value: sarama.StringEncoder("Hello, Kafka!"),
    }

    partition, offset, err := producer.SendMessage(message)
    if err != nil {
        log.Fatalf("Failed to send message: %s", err)
    }

    log.Printf("Message sent to partition %d at offset %d", partition, offset)
}

func TestSaslPlain(t *testing.T) {
    // 创建 Kafka 配置
    config := sarama.NewConfig()
    config.Version = sarama.V2_8_0_0 // 根据你的 Kafka 版本设置
    config.Net.SASL.Enable = true
    config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
    config.Net.SASL.User = "clientuser1" // 替换为你的用户名
    config.Net.SASL.Password = "pass123" // 替换为你的密码
    config.Net.SASL.Enable = true
    config.Producer.Return.Successes = true

    // 创建生产者
    //producer, err := sarama.NewSyncProducer([]string{"192.168.2.159:19092", "192.168.2.159:19094"}, config) // 替换为你的 Kafka broker
    producer, err := sarama.NewSyncProducer([]string{"192.168.2.159:9092"}, config) // 替换为你的 Kafka broker
    if err != nil {
        log.Fatalf("Failed to create producer: %s", err)
    }
    defer producer.Close()

    // 发送消息
    message := &sarama.ProducerMessage{
        Topic: "your_topic", // 替换为你的主题
        Value: sarama.StringEncoder("Hello, Kafka!"),
    }

    partition, offset, err := producer.SendMessage(message)
    if err != nil {
        log.Fatalf("Failed to send message: %s", err)
    }

    log.Printf("Message sent to partition %d at offset %d", partition, offset)
}

var (
    SHA256 scram.HashGeneratorFcn = sha256.New
    SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
    *scram.Client
    *scram.ClientConversation
    scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
    x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
    if err != nil {
        return err
    }
    x.ClientConversation = x.Client.NewConversation()
    return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
    response, err = x.ClientConversation.Step(challenge)
    return
}

func (x *XDGSCRAMClient) Done() bool {
    return x.ClientConversation.Done()
}
```

#### 消费者

```go
package main

import (
    "context"
    "crypto/sha256"
    "crypto/sha512"
    "github.com/xdg-go/scram"
    "log"
    "os"
    "testing"

    "github.com/Shopify/sarama"
)

type consumerGroupHandler struct{}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for message := range claim.Messages() {

        log.Printf("Message topic:%q partition:%d offset:%d data %v\n", message.Topic, message.Partition, message.Offset, string(message.Value))
        sess.MarkMessage(message, "")
    }
    return nil
}

func TestConsumer1(t *testing.T) {
    brokers := []string{"192.168.2.159:9092"}
    group := "test-group"
    topics := []string{"you_topic"}
    sarama.Logger = log.New(os.Stdout, "", log.Ltime)
    config := sarama.NewConfig()
    config.Version = sarama.V2_1_0_0
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    config.Consumer.Offsets.Initial = sarama.OffsetNewest

    consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
    if err != nil {
        log.Fatalf("Error creating consumer group: %v", err)
    }
    defer consumerGroup.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    for {
        if err := consumerGroup.Consume(ctx, topics, consumerGroupHandler{}); err != nil {
            log.Fatalf("Error from consumer: %v", err)
        }
    }
}

func TestSaslPlainConsumerSCRAM512(t *testing.T) {
    // 创建 Kafka 配置
    config := sarama.NewConfig()
    config.Version = sarama.V2_8_0_0 // 根据你的 Kafka 版本设置
    config.Net.SASL.Enable = true
    config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
    config.Net.SASL.User = "clientuser1" // 替换为你的用户名
    config.Net.SASL.Password = "pass123" // 替换为你的密码
    config.Net.SASL.Enable = true
    config.Producer.Return.Successes = true
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    consumerGroup, err := sarama.NewConsumerGroup([]string{"192.168.2.159:9092"}, "your_topic", config)
    if err != nil {
        log.Fatalf("Error creating consumer group: %v", err)
    }
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    for {
        if err := consumerGroup.Consume(ctx, []string{"your_topic"}, consumerGroupHandler{}); err != nil {
            log.Fatalf("Error from consumer: %v", err)
        }

    }

}

func TestSaslPlainConsumer(t *testing.T) {
    // 创建 Kafka 配置
    config := sarama.NewConfig()
    config.Version = sarama.V2_8_0_0 // 根据你的 Kafka 版本设置
    config.Net.SASL.Enable = true
    config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
    config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
    config.Net.SASL.User = "clientuser1" // 替换为你的用户名
    config.Net.SASL.Password = "pass123" // 替换为你的密码
    config.Net.SASL.Enable = true
    config.Producer.Return.Successes = true
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    consumerGroup, err := sarama.NewConsumerGroup([]string{"192.168.2.159:9092"}, "test-group", config)
    if err != nil {
        log.Fatalf("Error creating consumer group: %v", err)
    }
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    for {
        if err := consumerGroup.Consume(ctx, []string{"your_topic"}, consumerGroupHandler{}); err != nil {
            log.Fatalf("Error from consumer: %v", err)
        }

    }
}

var (
    SHA256 scram.HashGeneratorFcn = sha256.New
    SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
    *scram.Client
    *scram.ClientConversation
    scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
    x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
    if err != nil {
        return err
    }
    x.ClientConversation = x.Client.NewConversation()
    return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
    response, err = x.ClientConversation.Step(challenge)
    return
}

func (x *XDGSCRAMClient) Done() bool {
    return x.ClientConversation.Done()
}
```

### 面试常见问题



[**https://github.com/IcyBiscuit/Java-Guide/blob/master/docs/system-design/distributed-system/message-queue/Kafka%E5%B8%B8%E8%A7%81%E9%9D%A2%E8%AF%95%E9%A2%98%E6%80%BB%E7%BB%93.md**](https://github.com/IcyBiscuit/Java-Guide/blob/master/docs/system-design/distributed-system/message-queue/Kafka常见面试题总结.md)

[**https://developer.aliyun.com/article/740170**](https://developer.aliyun.com/article/740170)

[**https://www.lixueduan.com/posts/kafka/09-avoid-msg-lost/**](https://www.lixueduan.com/posts/kafka/09-avoid-msg-lost/)

#### 如何避免消息丢失。

生产者：1、使用同步生产。并设置消息等到follower全部同步完，再发送下一条。

2、设置retry值，当出现错误时重试。

消费者：关闭自动提交，手动提交偏移量

kafka服务端：

服务端要集群部署，并且数据的存储相关的配置的项副本数都要设置为集群对应的节点数。

包括topic,消费者消费情况等配置。

```toml
KAFKA_CFG_OFFSETS_TOPIC_REPLICATION_FACTOR: 2 # 消费组提交的偏移量信息在几个 broker 存储，默认值：1
KAFKA_CFG_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 2 # 事务状态 broker 存储，默认值：1
KAFKA_CFG_DEFAULT_REPLICATION_FACTOR: 2 # 默认topic数据有几个副本。
```

- 4）设置 unclean.leader.election.enable = false。这是 Broker 端的参数，它控制的是哪些 Broker 有资格竞选分区的 Leader。如果一个 Broker 落后原先的 Leader 太多，那么它一旦成为新的 Leader，必然会造成消息的丢失。故一般都要将该参数设置成 false，即不允许这种情况的发生。
- 6）设置 min.insync.replicas > 1。这依然是 Broker 端参数，控制的是消息至少要被写入到多少个副本才算是“已提交”。设置成大于 1 可以提升消息持久性。在实际环境中千万不要使用默认值 1。
- 7）确保 replication.factor > min.insync.replicas。如果两者相等，那么只要有一个副本挂机，整个分区就无法正常工作了。我们不仅要改善消息的持久性，防止数据丢失，还要在不降低可用性的基础上完成。推荐设置成 replication.factor = min.insync.replicas + 1。







#### 如何避免重复消费。

https://www.modb.pro/db/73387

https://developer.aliyun.com/article/740170

##### 重复消费的原因

原因1：消费者宕机、重启或者被强行kill进程，导致消费者消费的offset没有提交。

原因2：设置`enable.auto.commit`
为true，如果在关闭消费者进程之前，取消了消费者的订阅，则有可能部分offset没提交，下次重启会重复消费。

原因3：消费后的数据，当offset还没有提交时，Partition就断开连接。比如，通常会遇到消费的数据，处理很耗时，导致超过了Kafka的`session timeout.ms`
时间，那么就会触发reblance重平衡，此时可能存在消费者offset没提交，会导致重平衡后重复消费。

 **max.poll.interval.ms 活跃检测机制简介**

出现“活锁”的情况，是它持续的发送心跳，但是没有处理。为了预防消费者在这种情况下一直持有分区， 在此基础上，如果你调用的 poll 的频率大于最大间隔，则客户端将主动地离开组，以便其他消费者接管该分区

##### 重复消费的解决方法

1. 提高消费者的处理速度。例如：对消息处理中比较耗时的步骤可通过异步的方式进行处理、利用多线程处理等。在缩短单条消息消费的同时，根据实际场景可将`max.poll.interval.ms`
   值设置大一点，避免不必要的Rebalance。可根据实际消息速率适当调小`max.poll.records`
   的值。
2. 引入消息去重机制。例如：生成消息时，在消息中加入唯一标识符如消息id等。在消费端，可以保存最近的`max.poll.records`
   条消息id到redis或mysql表中，这样在消费消息时先通过查询去重后，再进行消息的处理。
3. 保证消费者逻辑幂等。可以查看博客《一文理解如何实现接口的幂等性》