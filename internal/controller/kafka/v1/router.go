package v1

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-product/internal/usecase"
)

type KafkaRouter struct {
	uc usecase.Product
}

func NewKafkaRouter(uc usecase.Product) *KafkaRouter {
	return &KafkaRouter{
		uc: uc,
	}
}

type KafkaProductQuantityUpdatedMessage struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductQuantity int64     `json:"product_quantity"`
}

func (r *KafkaRouter) HandleProductQuantityUpdated(msg *kafka.Message) error {
	var message KafkaProductQuantityUpdatedMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return err
	}

	log.Println("product quantity updated", message)
	// TODO: update product quantity in dynamo db
	return nil
}
