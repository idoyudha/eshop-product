package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/idoyudha/eshop-product/internal/entity"
	d "github.com/idoyudha/eshop-product/pkg/aws"
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
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.ProductTable),
		FilterExpression: aws.String("attribute_not_exists(deleted_at)"),
	}

	result, err := r.Client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan products: %w", err)
	}

	products := make([]entity.Product, 0, len(result.Items))
	for _, item := range result.Items {
		product := entity.Product{}

		product.ID = item["id"].(*types.AttributeValueMemberS).Value
		product.SKU = item["sku"].(*types.AttributeValueMemberS).Value
		product.Name = item["name"].(*types.AttributeValueMemberS).Value

		if imgURL, ok := item["image_url"]; ok {
			product.ImageURL = imgURL.(*types.AttributeValueMemberS).Value
		}

		product.Description = item["description"].(*types.AttributeValueMemberS).Value

		if categoryID, err := strconv.Atoi(item["category_id"].(*types.AttributeValueMemberN).Value); err == nil {
			product.CategoryID = categoryID
		}
		if price, err := strconv.ParseFloat(item["price"].(*types.AttributeValueMemberN).Value, 64); err == nil {
			product.Price = price
		}
		if quantity, err := strconv.Atoi(item["quantity"].(*types.AttributeValueMemberN).Value); err == nil {
			product.Quantity = quantity
		}

		if createdAt, err := time.Parse(time.RFC3339, item["created_at"].(*types.AttributeValueMemberS).Value); err == nil {
			product.CreatedAt = createdAt
		}
		if updatedAt, err := time.Parse(time.RFC3339, item["updated_at"].(*types.AttributeValueMemberS).Value); err == nil {
			product.UpdatedAt = updatedAt
		}

		products = append(products, product)
	}
	return &products, nil
}

func (r *ProductDynamoRepo) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.ProductTable),
		KeyConditionExpression: aws.String("id = :id"),
		FilterExpression:       aws.String("attribute_not_exists(deleted_at)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := r.Client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("product not found with id: %s", id)
	}

	item := result.Items[0]

	product := &entity.Product{}
	product.ID = id
	product.SKU = item["sku"].(*types.AttributeValueMemberS).Value
	product.Name = item["name"].(*types.AttributeValueMemberS).Value
	product.Description = item["description"].(*types.AttributeValueMemberS).Value

	if imgURL, ok := item["image_url"]; ok {
		product.ImageURL = imgURL.(*types.AttributeValueMemberS).Value
	}

	if categoryID, err := strconv.Atoi(item["category_id"].(*types.AttributeValueMemberN).Value); err == nil {
		product.CategoryID = categoryID
	}
	if price, err := strconv.ParseFloat(item["price"].(*types.AttributeValueMemberN).Value, 64); err == nil {
		product.Price = price
	}
	if quantity, err := strconv.Atoi(item["quantity"].(*types.AttributeValueMemberN).Value); err == nil {
		product.Quantity = quantity
	}

	if createdAt, err := time.Parse(time.RFC3339, item["created_at"].(*types.AttributeValueMemberS).Value); err == nil {
		product.CreatedAt = createdAt
	}
	if updatedAt, err := time.Parse(time.RFC3339, item["updated_at"].(*types.AttributeValueMemberS).Value); err == nil {
		product.UpdatedAt = updatedAt
	}

	return product, nil
}

func (r *ProductDynamoRepo) GetProductsByCategory(ctx context.Context, categoryID int) ([]entity.Product, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.ProductTable),
		IndexName:              aws.String("category_id-index"),
		KeyConditionExpression: aws.String("category_id = :category_id"),
		FilterExpression:       aws.String("attribute_not_exists(deleted_at)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category_id": &types.AttributeValueMemberN{Value: strconv.Itoa(categoryID)},
		},
	}

	result, err := r.Client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by category: %w", err)
	}

	if len(result.Items) == 0 {
		return []entity.Product{}, nil
	}

	products := make([]entity.Product, 0, len(result.Items))
	for _, item := range result.Items {
		product := entity.Product{}

		product.ID = item["id"].(*types.AttributeValueMemberS).Value
		product.CategoryID = categoryID
		product.SKU = item["sku"].(*types.AttributeValueMemberS).Value
		product.Name = item["name"].(*types.AttributeValueMemberS).Value
		product.Description = item["description"].(*types.AttributeValueMemberS).Value

		if price, err := strconv.ParseFloat(item["price"].(*types.AttributeValueMemberN).Value, 64); err == nil {
			product.Price = price
		}
		if quantity, err := strconv.Atoi(item["quantity"].(*types.AttributeValueMemberN).Value); err == nil {
			product.Quantity = quantity
		}

		if createdAt, err := time.Parse(time.RFC3339, item["created_at"].(*types.AttributeValueMemberS).Value); err == nil {
			product.CreatedAt = createdAt
		}
		if updatedAt, err := time.Parse(time.RFC3339, item["updated_at"].(*types.AttributeValueMemberS).Value); err == nil {
			product.UpdatedAt = updatedAt
		}

		if imgURL, ok := item["image_url"]; ok {
			product.ImageURL = imgURL.(*types.AttributeValueMemberS).Value
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductDynamoRepo) Update(ctx context.Context, product *entity.Product) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.ProductTable),
		Key: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: product.ID},
			"category_id": &types.AttributeValueMemberN{Value: strconv.Itoa(product.CategoryID)},
		},
		UpdateExpression: aws.String(
			"SET #name = :name, " +
				"image_url = :image_url, " +
				"description = :description, " +
				"price = :price, " +
				"quantity = :quantity, " +
				"updated_at = :updated_at",
		),
		ExpressionAttributeNames: map[string]string{
			"#name": "name", // name is reserved word in dynamodb
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":        &types.AttributeValueMemberS{Value: product.Name},
			":image_url":   &types.AttributeValueMemberS{Value: product.ImageURL},
			":description": &types.AttributeValueMemberS{Value: product.Description},
			":price":       &types.AttributeValueMemberN{Value: strconv.FormatFloat(product.Price, 'f', 2, 64)},
			":quantity":    &types.AttributeValueMemberN{Value: strconv.Itoa(product.Quantity)},
			":updated_at":  &types.AttributeValueMemberS{Value: product.UpdatedAt.Format(time.RFC3339)},
		},
		ConditionExpression: aws.String("attribute_not_exists(deleted_at)"),
	}

	_, err := r.Client.UpdateItem(ctx, input)
	if err != nil {
		var ccf *types.ConditionalCheckFailedException
		if ok := errors.As(err, &ccf); ok {
			return fmt.Errorf("product not found or has been deleted, id: %s", product.ID)
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *ProductDynamoRepo) Delete(ctx context.Context, productID string, categoryID int) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.ProductTable),
		Key: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: productID},
			"category_id": &types.AttributeValueMemberN{Value: strconv.Itoa(categoryID)},
		},
		UpdateExpression: aws.String("SET deleted_at = :deleted_at"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":deleted_at": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
		ConditionExpression: aws.String("attribute_not_exists(deleted_at)"),
	}

	_, err := r.Client.UpdateItem(ctx, input)
	if err != nil {
		var ccf *types.ConditionalCheckFailedException
		if ok := errors.As(err, &ccf); ok {
			return fmt.Errorf("product not found or already deleted, id: %s", productID)
		}
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
