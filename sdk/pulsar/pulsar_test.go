package d

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cast"
	"log"
	"testing"
	"time"
)

func TestProduct(t *testing.T) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://192.168.179.99:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "trade-5",
	})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Key:     cast.ToString(i),
			Payload: []byte("hello" + cast.ToString(i)),
		})
		time.Sleep(time.Millisecond * 100)
	}

	defer producer.Close()
	if err != nil {
		fmt.Println("Failed to publish message", err)
	}
	fmt.Println("Published message")
}
func TestConsumer(t *testing.T) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://192.168.2.231:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "persistent://public/default/exchange-simulator",
		SubscriptionName: "sub-2",       // 订阅名称
		Type:             pulsar.Shared, // 订阅类型: 独占模式
	})
	if err != nil {
		t.Fatal(err)
	}
	defer consumer.Close()
	//if err := consumer.Seek(msgID); err != nil {
	//	log.Println(err)
	//	return
	//}
	//for {
	//	msg, err := consumer.Receive(context.Background())
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	var d pb.ExecuteReportMessage
	//	if err := proto.Unmarshal(msg.Payload(), &d); err != nil {
	//		t.Log("err", err)
	//		return
	//	}
	//	log.Printf("%+v", d)
	//
	//}
}

// 消费shared模式多个消费者
func TestMutliConsumer(t *testing.T) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://192.168.179.99:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	for i := 0; i < 3; i++ {
		consumer, err := client.Subscribe(pulsar.ConsumerOptions{
			Topic:            "my-topic",
			SubscriptionName: "sub-2",       // 订阅名称
			Type:             pulsar.Shared, // 订阅类型
		})

		if err != nil {
			t.Fatal(err)
		}
		defer consumer.Close()
		/*2023/06/07 17:53:43 consumer 0 receive message
		consume: hello0
		2023/06/07 17:53:43 consumer 1 receive message
		consume: hello1
		2023/06/07 17:53:43 consumer 2 receive message
		consume: hello2
		2023/06/07 17:53:43 consumer 0 receive message
		consume: hello3
		2023/06/07 17:53:43 consumer 1 receive message
		consume: hello4
		2023/06/07 17:53:43 consumer 2 receive message
		consume: hello5
		2023/06/07 17:53:43 consumer 0 receive message
		consume: hello6
		2023/06/07 17:53:43 consumer 1 receive message
		consume: hello7
		2023/06/07 17:53:44 consumer 2 receive message
		consume: hello8
		2023/06/07 17:53:44 consumer 0 receive message
		consume: hello9 */
		go func(i int) {
			for {
				msg, err := consumer.Receive(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				log.Printf("consumer %d receive message", i)
				if err := processMsg(msg); err != nil {
					consumer.Nack(msg)
				} else {
					consumer.Ack(msg)
				}
			}
		}(i)
	}
	select {}
}
func TestPartsWithFailover(t *testing.T) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://192.168.179.99:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})

	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	for i := 0; i < 3; i++ {
		consumer, err := client.Subscribe(pulsar.ConsumerOptions{
			Topic:            "trade-5",
			SubscriptionName: "test1",
			Type:             pulsar.Failover, // 订阅类型: 独占模式
		})
		if err != nil {
			t.Fatal(err)
		}
		defer consumer.Close()
		go func(i int) {
			for {
				msg, err := consumer.Receive(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				log.Printf("customer %d receivce message data %v \n", i, string(msg.Payload()))
				consumer.Ack(msg)
			}
		}(i)

	}
	select {}

}

func processMsg(msg pulsar.Message) error {
	fmt.Printf("consume: %s \n", msg.Payload())
	return nil
}
func TestReader(t *testing.T) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://192.168.179.99:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	reader, err := client.CreateReader(pulsar.ReaderOptions{
		Topic: "my-topic",
		//Start reading from the end topic, only getting messages published after the reader was created
		StartMessageID: pulsar.LatestMessageID(),
		//StartMessageIDInclusive, if true, the reader will start at the `StartMessageID`, included. Default is `false` and the reader will start from the "next" message
		StartMessageIDInclusive: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	d, err := reader.Next(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("message LedgerID = %v EntryID = %v data payload %v\n", d.ID().LedgerID(), d.ID().EntryID(), string(d.Payload()))

}
