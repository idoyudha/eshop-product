package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/idoyudha/eshop-product/config"
)

const (
	ProductGroup          = "product-group"
	ProductQtyUpdateTopic = "product-quantity-updated"
	maxRetries            = 5
	retryDelay            = 2 * time.Second
)

type ConsumerServer struct {
	Consumer *kafka.Consumer
}

func NewKafkaConsumer(kafkaCfg config.Kafka) (*ConsumerServer, error) {
	log.Printf("Creating Kafka consumer with broker URL: %s", kafkaCfg.Broker)

	config := &kafka.ConfigMap{
		"bootstrap.servers":         kafkaCfg.Broker,
		"group.id":                  ProductGroup,
		"auto.offset.reset":         "earliest",
		"session.timeout.ms":        45000,
		"heartbeat.interval.ms":     15000,
		"metadata.max.age.ms":       300000,
		"enable.auto.commit":        true,
		"auto.commit.interval.ms":   5000,
		"enable.partition.eof":      false,
		"allow.auto.create.topics":  true,
		"max.poll.interval.ms":      300000,
		"max.partition.fetch.bytes": 1048576,
		"fetch.max.bytes":           52428800,
	}

	log.Printf("Kafka configuration: broker=%s, group=%s, auto.offset.reset=earliest",
		kafkaCfg.Broker, ProductGroup)

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	topics := []string{
		ProductQtyUpdateTopic,
	}

	var subscribeErr error
	for i := 0; i < maxRetries; i++ {
		subscribeErr = c.SubscribeTopics(topics, nil)
		if subscribeErr == nil {
			log.Printf("successfully subscribed to topics")
			break
		}
		log.Printf("attempt %d: failed to subscribe to topics: %v. retrying in %v...",
			i+1, subscribeErr, retryDelay)
		time.Sleep(retryDelay)
	}

	if subscribeErr != nil {
		c.Close()
		return nil, fmt.Errorf("failed to subscribe to topics after %d attempts: %v",
			maxRetries, subscribeErr)
	}

	return &ConsumerServer{
		Consumer: c,
	}, nil
}

func (c *ConsumerServer) Close() error {
	if c.Consumer != nil {
		return c.Consumer.Close()
	}
	return nil
}
