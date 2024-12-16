package v1

import (
	"encoding/json"

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

type KafkaUpdateProductQuantityMessage struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductQuantity int64     `json:"product_quantity"`
}

func (r *KafkaRouter) HandleProductAmountUpdated(msg *kafka.Message) error {
	var message KafkaUpdateProductQuantityMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return err
	}

	// TODO: update product quantity in dynamo db
	return nil
}
