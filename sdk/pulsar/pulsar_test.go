package pulsar

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"testing"
	"time"
)

func TestProduct(t *testing.T) {
	publicMessage()
}
func TestPartitionSharedMode(t *testing.T) {
	PartitionSharedMode()
}
func TestPartitionExclusiveSharedMode(t *testing.T) {
	PartitionExclusiveSharedMode()
}
func TestPartitionKeySharedMode(t *testing.T) {
	PartitionKeySharedMode()
}
func TestRetryDLQ(t *testing.T) {
	RetryDLQ()
}
func TestFailoverConsumer(t *testing.T) {
	FailoverConsumer()
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
			Type:             pulsar.Failover,
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
		URL:               "pulsar://192.168.2.231:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	for i := 0; i < 3; i++ {
		go func(i int) {
			reader, err := client.CreateReader(pulsar.ReaderOptions{
				Topic: "product5",
				//Start reading from the end topic, only getting messages published after the reader was created
				StartMessageID: pulsar.EarliestMessageID(),
				//StartMessageIDInclusive, if true, the reader will start at the `StartMessageID`, included. Default is `false` and the reader will start from the "next" message
				StartMessageIDInclusive: true,
			})
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()

			for reader.HasNext() {
				d, err := reader.Next(context.Background())
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("partitio =%v index =%v message LedgerID = %v EntryID = %v data payload %v\n", d.ID().PartitionIdx(), i, d.ID().LedgerID(), d.ID().EntryID(), string(d.Payload()))

			}

		}(i)
	}
	select {}
}
