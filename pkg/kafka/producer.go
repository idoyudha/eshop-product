package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/idoyudha/eshop-product/pkg/logger"
)

type ProducerServer struct {
	producer *kafka.Producer
	l        logger.Interface
}

func NewKafkaProducer(brokerURL string, l logger.Interface) (*ProducerServer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &ProducerServer{
		producer: p,
		l:        l,
	}, nil
}

func (s *ProducerServer) Close() {
	s.producer.Close()
}

func (s *ProducerServer) Produce(topic string, key []byte, value []byte) error {
	return s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, nil)
}
