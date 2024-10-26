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
	config.Consumer.Offsets.AutoCommit.Enable = false
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
	config.Consumer.Offsets.AutoCommit.Enable = false
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
	config.Consumer.Offsets.AutoCommit.Enable = false

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
