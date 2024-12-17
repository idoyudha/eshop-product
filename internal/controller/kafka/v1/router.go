package v1

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-product/internal/usecase"
	kafkaConSrv "github.com/idoyudha/eshop-product/pkg/kafka"
	"github.com/idoyudha/eshop-product/pkg/logger"
)

type kafkaConsumerRoutes struct {
	ucp usecase.Product
	l   logger.Interface
}

func KafkaNewRouter(
	ucp usecase.Product,
	l logger.Interface,
	c *kafkaConSrv.ConsumerServer,
) error {
	routes := &kafkaConsumerRoutes{
		ucp: ucp,
		l:   l,
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
			return nil
		default:
			ev, err := c.Consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// log.Println("CONSUME PRODUCT SERVICE!!")
				// Errors are informational and automatically handled by the consumer
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				l.Error("Error reading message: ", err)
				continue
			}

			switch *ev.TopicPartition.Topic {
			case kafkaConSrv.ProductQtyUpdateTopic:
				if err := routes.handleProductQuantityUpdated(ev); err != nil {
					l.Error("Failed to handle product update: %w", err)
				}
			default:
				l.Info("Unknown topic: %s", *ev.TopicPartition.Topic)
			}
			log.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}

	return nil
}

type KafkaProductQuantityUpdatedMessage struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductQuantity int64     `json:"product_quantity"`
}

func (r *kafkaConsumerRoutes) handleProductQuantityUpdated(msg *kafka.Message) error {
	var message KafkaProductQuantityUpdatedMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return err
	}

	log.Println("Received product quantity updated", message)
	// TODO: update product quantity in dynamo db
	return nil
}
