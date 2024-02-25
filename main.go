package main


import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	runProducer()
}

func runProducer() {

	producer, err := createProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to connect to Kafka broker: %v", err)
	}

	eventToSend := NewTestEvent("hello world", "john")
	encodedEvent, err := EncodeTestEvent(eventToSend)
	if err != nil {
		log.Fatalf("Failed to encode test event: %v", err)
	}

	sendMessage(producer, encodedEvent, "test_events")
	closeProducer(producer)
}

func createProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka producer: %w", err)
	}

	return producer, nil
}

func closeProducer(producer sarama.SyncProducer) {
	if err := producer.Close(); err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)
	} else {
		log.Println("Kafka producer closed.")
	}
}

func sendMessage(producer sarama.SyncProducer, msgBytes []byte, topic string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}

