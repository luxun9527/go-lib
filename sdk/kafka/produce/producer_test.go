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
		Partition: 0,   //指定分区，可选
		Key:       nil, //如果指定了key，则会根据可以选择一个分区。相同的key 会被发送到同一个分区
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
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 创建生产者
	producer, err := sarama.NewSyncProducer([]string{"192.168.2.159:19092", "192.168.2.159:19094"}, config) // 替换为你的 Kafka broker
	//producer, err := sarama.NewSyncProducer([]string{"192.168.2.159:9092"}, config) // 替换为你的 Kafka broker
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// 发送消息
	message := &sarama.ProducerMessage{
		Topic: "your_topic1", // 替换为你的主题
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
