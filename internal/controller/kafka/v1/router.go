package v1

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-product/internal/usecase"
	"github.com/idoyudha/eshop-product/pkg/kafka"
	"github.com/idoyudha/eshop-product/pkg/logger"
)

func KafkaNewRouter(
	ucp usecase.Product,
	l logger.Interface,
	c *kafka.ConsumerServer,
) {
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
		default:
			ev, err := c.Consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// log.Println("CONSUME PRODUCT SERVICE!!")
				// Errors are informational and automatically handled by the consumer
				continue
			}
			log.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}
}

type KafkaProductQuantityUpdatedMessage struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductQuantity int64     `json:"product_quantity"`
}

// func (r *KafkaRouter) HandleProductQuantityUpdated(contex context.Context, msg *kafka.Message) error {
// 	var message KafkaProductQuantityUpdatedMessage

// 	if err := json.Unmarshal(msg.Value, &message); err != nil {
// 		return err
// 	}

// 	log.Println("product quantity updated", message)
// 	// TODO: update product quantity in dynamo db
// 	return nil
// }
