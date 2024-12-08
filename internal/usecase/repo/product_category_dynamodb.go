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
	d "github.com/idoyudha/eshop-product/pkg/dynamodb"
)

type CategoryDynamoRepo struct {
	*d.DynamoDB
}

func NewCategoryDynamoRepo(db *d.DynamoDB) *CategoryDynamoRepo {
	return &CategoryDynamoRepo{
		db,
	}
}

func (r *CategoryDynamoRepo) Save(ctx context.Context, category *entity.Category) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.CategoryTable),
		Item: map[string]types.AttributeValue{
			"id":         &types.AttributeValueMemberN{Value: strconv.Itoa(category.ID)},
			"name":       &types.AttributeValueMemberS{Value: category.Name},
			"parent_id":  &types.AttributeValueMemberN{Value: strconv.Itoa(*category.ParentID)},
			"created_at": &types.AttributeValueMemberS{Value: category.CreatedAt.String()},
			"updated_at": &types.AttributeValueMemberS{Value: category.UpdatedAt.String()},
		},
	}

	_, err := r.Client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save category: %w", err)
	}
	return nil
}

func (r *CategoryDynamoRepo) GetCategories(ctx context.Context) (*[]entity.Category, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.CategoryTable),
		FilterExpression: aws.String("attribute_not_exists(deleted_at)"),
	}

	result, err := r.Client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan products: %w", err)
	}

	categories := make([]entity.Category, 0, len(result.Items))
	for _, item := range result.Items {
		category := entity.Category{}

		if id, err := strconv.Atoi(item["id"].(*types.AttributeValueMemberN).Value); err == nil {
			category.ID = id
		}
		category.Name = item["name"].(*types.AttributeValueMemberS).Value
		if parentID, err := strconv.Atoi(item["parent_id"].(*types.AttributeValueMemberN).Value); err == nil {
			category.ParentID = &parentID
		}
		if createdAt, err := time.Parse(time.RFC3339, item["created_at"].(*types.AttributeValueMemberS).Value); err == nil {
			category.CreatedAt = createdAt
		}
		if updatedAt, err := time.Parse(time.RFC3339, item["updated_at"].(*types.AttributeValueMemberS).Value); err == nil {
			category.UpdatedAt = updatedAt
		}

		categories = append(categories, category)
	}

	return &categories, nil
}

func (r *CategoryDynamoRepo) Update(ctx context.Context, category *entity.Category) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.CategoryTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: strconv.Itoa(category.ID)},
		},
		UpdateExpression: aws.String(
			"SET #name = :name, " +
				"parent_id = :parent_id, " +
				"updated_at = :updated_at",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":       &types.AttributeValueMemberS{Value: category.Name},
			":parent_id":  &types.AttributeValueMemberN{Value: strconv.Itoa(*category.ParentID)},
			":updated_at": &types.AttributeValueMemberS{Value: category.UpdatedAt.Format(time.RFC3339)},
		},
		ConditionExpression: aws.String("attribute_not_exists(deleted_at)"),
	}

	_, err := r.Client.UpdateItem(ctx, input)
	if err != nil {
		var ccf *types.ConditionalCheckFailedException
		if ok := errors.As(err, &ccf); ok {
			return fmt.Errorf("category not found or has been deleted, id: %d", category.ID)
		}
		return fmt.Errorf("failed to category product: %w", err)
	}
	return nil
}

func (r *CategoryDynamoRepo) Delete(ctx context.Context, id string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.CategoryTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: id},
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
			return fmt.Errorf("category not found or has been deleted, id: %s", id)
		}
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
