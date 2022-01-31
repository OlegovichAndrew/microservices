package transport

import (
	"github.com/Shopify/sarama"
	"log"
)

var kafkaVersion = sarama.V3_0_0_0

func CreateProducer(brokerList []string, clientID string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.ClientID = clientID

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func SendMessage(producer sarama.SyncProducer, topic, message string) error {
	_, _, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})

	return err
}

func CreateTopic(brokerList []string, topicName string, nPartitions int32, replicas int16) error {
	config := sarama.NewConfig()
	config.Version = kafkaVersion

	admin, err := sarama.NewClusterAdmin(brokerList, config)
	if err != nil {
		return err
	}
	defer func() { _ = admin.Close() }()

	err = admin.CreateTopic(topicName, &sarama.TopicDetail{
		NumPartitions:     nPartitions,
		ReplicationFactor: replicas,
	}, false)

	return err
}
