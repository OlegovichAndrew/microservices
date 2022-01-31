package transport

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

var kafkaVersion = sarama.V3_0_0_0

func CreateConsumerGroup(brokerList []string, clientID string, groupName string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.ClientID = clientID

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupName, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v\n", err)
	}
	return consumerGroup
}

func ConsumeMessages(ctx context.Context, group sarama.ConsumerGroup, topic string) {
	consumer := Consumer{
		ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := group.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Fatalf("Error from consumer: %v\n", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	wg.Wait()
}

type Consumer struct {
	ready chan bool
}

func (consumer *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Kafka: value=%s, time=%v, topic=%s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
