package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/idoyudha/eshop-product/internal/entity"
	awsService "github.com/idoyudha/eshop-product/pkg/aws"
)

type ProductDynamoRepo struct {
	*awsService.DynamoDB
}

func NewProductDynamoDBRepo(d *awsService.DynamoDB) *ProductDynamoRepo {
	return &ProductDynamoRepo{
		d,
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
			"category_id": &types.AttributeValueMemberS{Value: product.CategoryID},
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
		product.CategoryID = item["category_id"].(*types.AttributeValueMemberS).Value

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

	product.CategoryID = item["category_id"].(*types.AttributeValueMemberS).Value

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

func (r *ProductDynamoRepo) GetProductsByCategory(ctx context.Context, categoryID string) ([]entity.Product, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.ProductTable),
		IndexName:              aws.String("category_id-index"),
		KeyConditionExpression: aws.String("category_id = :category_id"),
		FilterExpression:       aws.String("attribute_not_exists(deleted_at)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category_id": &types.AttributeValueMemberS{Value: categoryID},
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

func (r *ProductDynamoRepo) GetProductsByCategories(ctx context.Context, categoryIDs []string) ([]entity.Product, error) {
	var wg sync.WaitGroup
	resultsChan := make(chan []entity.Product, len(categoryIDs))
	errorsChan := make(chan error, len(categoryIDs))

	for _, categoryID := range categoryIDs {
		wg.Add(1)
		go func(catID string) {
			defer wg.Done()

			products, err := r.GetProductsByCategory(ctx, catID)
			if err != nil {
				errorsChan <- fmt.Errorf("error fetching products for category %s: %w", catID, err)
				return
			}
			resultsChan <- products
		}(categoryID)
	}

	// wait for all goroutines to complete in a separate goroutine
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	var allProducts []entity.Product
	for products := range resultsChan {
		allProducts = append(allProducts, products...)
	}

	for err := range errorsChan {
		if err != nil {
			return nil, err
		}
	}

	if len(allProducts) == 0 {
		return []entity.Product{}, nil
	}

	return allProducts, nil
}

func (r *ProductDynamoRepo) Update(ctx context.Context, product *entity.Product) error {
	var updateParts []string
	expAttrNames := map[string]string{
		"#id":          "id",
		"#name":        "name",
		"#category_id": "category_id",
		"#image_url":   "image_url",
		"#description": "description",
		"#price":       "price",
		"#updated_at":  "updated_at",
	}
	expAttrValues := map[string]types.AttributeValue{}

	if product.Name != "" {
		updateParts = append(updateParts, "#name = :name")
		expAttrValues[":name"] = &types.AttributeValueMemberS{Value: product.Name}
	}
	if product.ImageURL != "" {
		updateParts = append(updateParts, "#image_url = :image_url")
		expAttrValues[":image_url"] = &types.AttributeValueMemberS{Value: product.ImageURL}
	}
	if product.Description != "" {
		updateParts = append(updateParts, "#description = :description")
		expAttrValues[":description"] = &types.AttributeValueMemberS{Value: product.Description}
	}
	if product.Price > 0 {
		updateParts = append(updateParts, "#price = :price")
		expAttrValues[":price"] = &types.AttributeValueMemberN{Value: strconv.FormatFloat(product.Price, 'f', 2, 64)}
	}

	updateParts = append(updateParts, "#updated_at = :updated_at")
	expAttrValues[":updated_at"] = &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)}

	updateExpression := "SET " + strings.Join(updateParts, ", ")

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.ProductTable),
		Key: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: product.ID},
			"category_id": &types.AttributeValueMemberS{Value: product.CategoryID},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expAttrNames,
		ExpressionAttributeValues: expAttrValues,
		ConditionExpression:       aws.String("attribute_exists(#id) AND attribute_exists(#category_id)"),
	}

	// First, let's verify if the item exists
	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.ProductTable),
		Key: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: product.ID},
			"category_id": &types.AttributeValueMemberS{Value: product.CategoryID},
		},
	}

	// Check if item exists before updating
	result, err := r.Client.GetItem(ctx, getInput)
	if err != nil {
		return fmt.Errorf("failed to check if product exists: %w", err)
	}

	if result.Item == nil {
		return fmt.Errorf("product with ID %s and category ID %s not found", product.ID, product.CategoryID)
	}

	// If item exists, proceed with update
	_, err = r.Client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *ProductDynamoRepo) GetCategoryByProductId(ctx context.Context, productID string) (*string, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.ProductTable),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: productID},
		},
		Limit: aws.Int32(1),
	}

	result, err := r.Client.Query(ctx, queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query product: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("product not found with ID: %s", productID)
	}

	item := result.Items[0]
	categoryID := item["category_id"].(*types.AttributeValueMemberS).Value
	return &categoryID, nil
}

func (r *ProductDynamoRepo) UpdateProductQty(ctx context.Context, productID, categoryID string, quantity int) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.ProductTable),
		Key: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: productID},
			"category_id": &types.AttributeValueMemberS{Value: categoryID},
		},
		UpdateExpression: aws.String("SET quantity = :quantity"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":quantity": &types.AttributeValueMemberN{Value: strconv.Itoa(quantity)},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	}

	_, err := r.Client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update product quantity: %w", err)
	}

	return nil
}

func (r *ProductDynamoRepo) Delete(ctx context.Context, productID string, categoryID string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.ProductTable),
		Key: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: productID},
			"category_id": &types.AttributeValueMemberS{Value: categoryID},
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
