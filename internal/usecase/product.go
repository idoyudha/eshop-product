package usecase

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-product/internal/entity"
	"github.com/idoyudha/eshop-product/pkg/kafka"
)

const (
	productCreatedTopic = "product-created"
	productUpdatedTopic = "product-updated"
)

type ProductUseCase struct {
	productRepoImage  ProductS3Repo
	productRepoDynamo ProductDynamoRepo
	producer          *kafka.ProducerServer
}

func NewProductUseCase(
	productRepoImage ProductS3Repo,
	productRepoDynamo ProductDynamoRepo,
	producer *kafka.ProducerServer,
) *ProductUseCase {
	return &ProductUseCase{
		productRepoImage:  productRepoImage,
		productRepoDynamo: productRepoDynamo,
		producer:          producer,
	}
}

type kafkaProductCreatedMessage struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  string  `json:"category_id"`
}

func (u *ProductUseCase) CreateProduct(ctx context.Context, product *entity.Product, imageFile *multipart.FileHeader) (*entity.Product, error) {
	err := product.GenerateProductID()
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	product.GenerateSKU()

	// save image to s3
	imageURL, err := u.productRepoImage.UploadImage(ctx, imageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	product.SetImageURL(imageURL)

	// save product to dynamo
	err = u.productRepoDynamo.Save(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	message := kafkaProductCreatedMessage{
		ID:          product.ID,
		SKU:         product.SKU,
		Name:        product.Name,
		ImageURL:    product.ImageURL,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CategoryID:  product.CategoryID,
	}

	err = u.producer.Produce(
		productCreatedTopic,
		[]byte(product.ID),
		message,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send product data to kafka: %w", err)
	}

	// sent product data to main warehouse
	return product, nil
}

func (u *ProductUseCase) GetProducts(ctx context.Context) (*[]entity.Product, error) {
	return u.productRepoDynamo.GetProducts(ctx)
}

func (u *ProductUseCase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return u.productRepoDynamo.GetProductByID(ctx, id)
}

func (u *ProductUseCase) GetProductsByCategory(ctx context.Context, categoryID string) ([]entity.Product, error) {
	return u.productRepoDynamo.GetProductsByCategory(ctx, categoryID)
}

type kafkaProductUpdatedMessage struct {
	ProductID          uuid.UUID `json:"product_id"`
	ProductName        string    `json:"product_name"`
	ProductImageURL    string    `json:"product_image_url"`
	ProductDescription string    `json:"product_description"`
	ProductPrice       float64   `json:"product_price"`
	ProductCategoryID  uuid.UUID `json:"product_category_id"`
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, product *entity.Product, imageFile *multipart.FileHeader) error {
	imageURL, err := u.productRepoImage.UploadImage(ctx, imageFile)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	product.SetImageURL(imageURL)

	err = u.productRepoDynamo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	message := kafkaProductUpdatedMessage{
		ProductID:          uuid.MustParse(product.ID),
		ProductName:        product.Name,
		ProductImageURL:    product.ImageURL,
		ProductDescription: product.Description,
		ProductPrice:       product.Price,
		ProductCategoryID:  uuid.MustParse(product.CategoryID),
	}

	err = u.producer.Produce(
		productUpdatedTopic,
		[]byte(product.ID),
		message,
	)
	if err != nil {
		// TODO: handle error, cancel the update if failed. or try use retry mechanism
		return fmt.Errorf("failed to produce kafka message: %w", err)
	}

	return nil
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, productID string, categoryID string) error {
	return u.productRepoDynamo.Delete(ctx, productID, categoryID)
}
