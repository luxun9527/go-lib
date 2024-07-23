package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Shopify/sarama"
)

type consumerGroupHandler struct{}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		var u User
		if err := json.Unmarshal(message.Value, &u); err != nil {
			return err
		}
		log.Printf("User: %+v\n", u)
		sess.MarkMessage(message, "")
	}
	return nil
}

type User struct {
	ID     string  `json:"id"`
	Detail *[]byte `json:"detail"`
}

func TestConsumer1(t *testing.T) {
	brokers := []string{"192.168.2.159:9092"}
	group := "test-group"
	topics := []string{"test-topic"}

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

	// 监听系统信号以优雅地关闭
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, topics, consumerGroupHandler{}); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
			// 检查上下文是否被取消以避免重新消费
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-sigterm
	fmt.Println("Terminating: via signal")
}
