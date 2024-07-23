package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/zeromicro/go-zero/core/stringx"
	"log"
	"testing"
)

type User struct {
	ID     string `json:"id"`
	Detail []byte `json:"detail"`
}

func TestProduce1(t *testing.T) {
	// Kafka broker 地址
	brokers := []string{"192.168.2.159:9092"}
	// Kafka 主题
	topic := "test-topic"

	// 创建一个新的配置
	config := sarama.NewConfig()
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
		Topic: topic,
		Value: sarama.ByteEncoder(d),
	}

	// 发送消息
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
}
