package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-product/internal/entity"
	"github.com/idoyudha/eshop-product/pkg/kafka"
)

const ProductUpdatedTopic = "product-updated"

type ProductUseCase struct {
	productRepoDynamo ProductDynamoRepo
	producer          *kafka.ProducerServer
}

func NewProductUseCase(productRepoDynamo ProductDynamoRepo, producer *kafka.ProducerServer) *ProductUseCase {
	return &ProductUseCase{
		productRepoDynamo: productRepoDynamo,
		producer:          producer,
	}
}

func (u *ProductUseCase) CreateProduct(ctx context.Context, product *entity.Product) error {
	return u.productRepoDynamo.Save(ctx, product)
}

func (u *ProductUseCase) GetProducts(ctx context.Context) (*[]entity.Product, error) {
	return u.productRepoDynamo.GetProducts(ctx)
}

func (u *ProductUseCase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return u.productRepoDynamo.GetProductByID(ctx, id)
}

func (u *ProductUseCase) GetProductsByCategory(ctx context.Context, categoryID int) ([]entity.Product, error) {
	return u.productRepoDynamo.GetProductsByCategory(ctx, categoryID)
}

type KafkaProductUpdatedMessage struct {
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductPrice float64   `json:"product_price"`
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	err := u.productRepoDynamo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	message := KafkaProductUpdatedMessage{
		ProductID:    uuid.MustParse(product.ID),
		ProductName:  product.Name,
		ProductPrice: product.Price,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		// TODO: handle error, cancel the update if failed. or try use retry mechanism
		return fmt.Errorf("failed to marshal kafka message: %w", err)
	}

	err = u.producer.Produce(
		ProductUpdatedTopic,
		[]byte(product.ID),
		messageBytes,
	)
	if err != nil {
		// TODO: handle error, cancel the update if failed. or try use retry mechanism
		return fmt.Errorf("failed to produce kafka message: %w", err)
	}

	return nil
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, productID string, categoryID int) error {
	return u.productRepoDynamo.Delete(ctx, productID, categoryID)
}
