package pulsar

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	pulsarLog "github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stringx"
	"log"
	"time"
)

const Partition5Topic = "partition5"
const Partition1Topic = "my-topic"
const Retry = "test-retry"

var pulsarClient pulsar.Client

func InitPulsarCline() {
	logger := logrus.StandardLogger()
	logger.Level = logrus.ErrorLevel
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://192.168.2.99:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
		Logger:            pulsarLog.NewLoggerWithLogrus(logger),
	})
	if err != nil {
		log.Fatal(err)
	}
	pulsarClient = client
}
func PartitionSharedMode() {
	InitPulsarCline()
	for i := 0; i < 3; i++ {
		go func(i int) {
			consumer, err := pulsarClient.Subscribe(pulsar.ConsumerOptions{
				Topic:            Partition5Topic,
				SubscriptionName: "partition5-sub",
				Type:             pulsar.Shared,
			})
			if err != nil {
				log.Panicln(err)
			}
			for {
				message, err := consumer.Receive(context.Background())
				if err != nil {
					log.Println(err)
					continue
				}

				log.Printf("recieve messageID %+v partitionID=%+v data =%+v consumerID=%v key=%v", message.ID(), message.ID().PartitionIdx(), string(message.Payload()), i, message.Key())
				consumer.Ack(message)
			}
		}(i)
	}
	select {}
}

// PartitionKeySharedMode keyshare模式
func PartitionKeySharedMode() {
	InitPulsarCline()
	for i := 0; i < 3; i++ {
		go func(i int) {
			consumer, err := pulsarClient.Subscribe(pulsar.ConsumerOptions{
				Topic:            Partition5Topic,
				SubscriptionName: "partition5-keyShared",
				Type:             pulsar.KeyShared,
			})
			if err != nil {
				log.Panicln(err)
			}
			for {
				message, err := consumer.Receive(context.Background())
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("recieve messageID %+v partitionID=%+v data =%+v consumerID=%v", message.ID(), message.ID().PartitionIdx(), string(message.Payload()), i)
				consumer.Ack(message)
			}
		}(i)
	}
	select {}
}

// PartitionExclusiveSharedMode keyshare模式
func PartitionExclusiveSharedMode() {
	InitPulsarCline()
	consumer, err := pulsarClient.Subscribe(pulsar.ConsumerOptions{
		Topic:            Partition5Topic,
		SubscriptionName: "partition5-sub",
		Type:             pulsar.Exclusive,
	})
	if err != nil {
		log.Panicln(err)
	}
	for {
		message, err := consumer.Receive(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("recieve messageID %+v partitionID=%+v data =%+v", message.ID(), message.ID().PartitionIdx(), string(message.Payload()))
		consumer.Nack(message)
	}

}
func FailoverConsumer() {
	InitPulsarCline()
	for i := 0; i < 3; i++ {
		go func(i int) {
			consumer, err := pulsarClient.Subscribe(pulsar.ConsumerOptions{
				Topic:            Partition5Topic,
				SubscriptionName: "partition5-keyShared",
				Type:             pulsar.Failover,
			})
			if err != nil {
				log.Panicln(err)
			}
			for {
				message, err := consumer.Receive(context.Background())
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("recieve messageID %+v partitionID=%+v data =%+v consumerID=%v", message.ID(), message.ID().PartitionIdx(), string(message.Payload()), i)
				consumer.Ack(message)
			}
		}(i)
	}
	select {}
}

// RetryDLQ 重试和死信队列
func RetryDLQ() {
	InitPulsarCline()
	d := &pulsar.DLQPolicy{
		MaxDeliveries:    3,
		RetryLetterTopic: "persistent://public/test/xxx-RETRY",
		DeadLetterTopic:  "persistent://public/test/xxx-DLQ",
	}

	consumer, err := pulsarClient.Subscribe(pulsar.ConsumerOptions{
		Topic:               "persistent://public/test/" + Retry,
		SubscriptionName:    "RetryDLQ",
		Type:                pulsar.Shared,
		RetryEnable:         true,
		DLQ:                 d,
		NackRedeliveryDelay: time.Second * 3,
	})

	if err != nil {
		log.Panicln(err)
	}
	for {
		message, err := consumer.Receive(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("recieve messageID %+v partitionID=%+v data =%+v ", message.ID(), message.ID().PartitionIdx(), string(message.Payload()))
		consumer.Nack(message)
		//consumer.ReconsumeLater()
	}

}
func publicMessage() {
	InitPulsarCline()
	producer, err := pulsarClient.CreateProducer(pulsar.ProducerOptions{
		Topic:       Partition5Topic,
		Name:        "",
		Properties:  nil,
		SendTimeout: time.Second * 10,
		//发送队列满了是否阻塞，不阻塞为true，发送将返回错误。
		DisableBlockIfQueueFull: false,
		MaxPendingMessages:      0,
		HashingScheme:           0,
		CompressionType:         0,
		CompressionLevel:        0,
		//消息的路由方式，决定消费发送到那个分区。模式不指定key的话，消息将轮询发送到每一个分区。
		MessageRouter: nil,
		//是否开启批量
		DisableBatching:                 true,
		BatchingMaxPublishDelay:         0,
		BatchingMaxMessages:             0,
		BatchingMaxSize:                 0,
		Interceptors:                    nil,
		Schema:                          nil,
		MaxReconnectToBroker:            nil,
		BackoffPolicy:                   nil,
		BatcherBuilderType:              0,
		PartitionsAutoDiscoveryInterval: 0,
		DisableMultiSchema:              false,
		Encryption:                      nil,
		EnableChunking:                  false,
		ChunkMaxMessageSize:             0,
		ProducerAccessMode:              0,
	})
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= 10; i++ {
		data := []byte(stringx.Randn(6) + cast.ToString(i))
		key := cast.ToString(i % 5)
		messageID, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
			Key:     key,
			Payload: data,
		})
		if err != nil {
			fmt.Println("Failed to publish message", err)
		}
		log.Printf("message  partitionID=%v data = %v key =%v", messageID.PartitionIdx(), string(data), key)
		time.Sleep(time.Millisecond * 100)
	}

	defer producer.Close()

	fmt.Println("Published message")
}
