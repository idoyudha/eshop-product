package repo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/idoyudha/eshop-product/internal/entity"
	d "github.com/idoyudha/eshop-product/pkg/dynamodb"
)

type ProductDynamoRepo struct {
	*d.DynamoDB
}

func NewProductRepo(db *d.DynamoDB) *ProductDynamoRepo {
	return &ProductDynamoRepo{
		db,
	}
}

func (r *ProductDynamoRepo) Save(ctx context.Context, product *entity.Product) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.ProductTable),
		Item: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: product.ID},
			"sku":         &types.AttributeValueMemberS{Value: product.SKU},
			"name":        &types.AttributeValueMemberS{Value: product.Name},
			"image_url":   &types.AttributeValueMemberS{Value: product.ImageURL},
			"description": &types.AttributeValueMemberS{Value: product.Description},
			"price":       &types.AttributeValueMemberN{Value: strconv.FormatFloat(product.Price, 'f', -1, 64)},
			"quantity":    &types.AttributeValueMemberN{Value: strconv.Itoa(product.Quantity)},
			"category_id": &types.AttributeValueMemberN{Value: strconv.Itoa(product.CategoryID)},
			"created_at":  &types.AttributeValueMemberS{Value: product.CreatedAt.String()},
			"updated_at":  &types.AttributeValueMemberS{Value: product.UpdatedAt.String()},
		},
	}

	_, err := r.Client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}

	return nil
}

func (r *ProductDynamoRepo) GetProducts(ctx context.Context) (*[]entity.Product, error) {
	// TODO: implement scan all of products
	return nil, nil
}

func (r *ProductDynamoRepo) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	// TODO: implement to get product by id
	return nil, nil
}

func (r *ProductDynamoRepo) GetProductsByCategory(ctx context.Context, categoryID int) ([]entity.Product, error) {
	// TODO: implement to query products by category_id
	return nil, nil
}

func (r *ProductDynamoRepo) Update(ctx context.Context, product *entity.Product) error {
	// TODO: implement update product
	return nil
}

func (r *ProductDynamoRepo) Delete(ctx context.Context, id string) error {
	// TODO: implement delete product
	return nil
}
