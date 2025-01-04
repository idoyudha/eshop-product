package repo

import (
	"context"
	"fmt"

	"github.com/idoyudha/eshop-product/internal/entity"
	rClient "github.com/idoyudha/eshop-product/pkg/redis"
	"github.com/redis/go-redis/v9"
)

const (
	categoryKeyPrefix = "category:"         // hash storing category data
	categorySetKey    = "categories"        // set of all category IDs
	categoryParentKey = "category_parents:" // set of child IDs for each parent
)

type CategoryRedisRepo struct {
	*rClient.RedisClient
}

func NewCategoryRedisRepo(redis *rClient.RedisClient) *CategoryRedisRepo {
	return &CategoryRedisRepo{
		redis,
	}
}

// save multiple categories at once, getting from dynamo db
func (r *CategoryRedisRepo) SaveAll(ctx context.Context, categories *[]entity.Category) error {
	pipe := r.Client.Pipeline()

	pipe.Del(ctx, categorySetKey)

	for _, category := range *categories {
		// store category data in hash
		categoryKey := categoryKeyPrefix + category.ID

		categoryData := map[string]interface{}{
			"name": category.Name,
		}

		if category.ParentID != nil {
			categoryData["parent_id"] = *category.ParentID
			// add to parent's children set
			pipe.SAdd(ctx, categoryParentKey+*category.ParentID, category.ID)
		}

		pipe.HSet(ctx, categoryKey, categoryData)
		pipe.SAdd(ctx, categorySetKey, category.ID)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save categories: %w", err)
	}

	return nil
}

func (r *CategoryRedisRepo) GetAll(ctx context.Context) (*[]entity.Category, error) {
	categoryIDs, err := r.Client.SMembers(ctx, categorySetKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get category IDs: %w", err)
	}

	if len(categoryIDs) == 0 {
		return nil, nil
	}

	pipe := r.Client.Pipeline()
	categoryDataCmds := make(map[string]*redis.MapStringStringCmd)

	for _, id := range categoryIDs {
		categoryDataCmds[id] = pipe.HGetAll(ctx, categoryKeyPrefix+id)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories data: %w", err)
	}

	categories := make([]entity.Category, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		data := categoryDataCmds[id].Val()
		if len(data) == 0 {
			continue
		}

		category := entity.Category{
			ID:   id,
			Name: data["name"],
		}

		if parentID, exists := data["parent_id"]; exists {
			category.ParentID = &parentID
		}

		categories = append(categories, category)
	}

	return &categories, nil
}

func (r *CategoryRedisRepo) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	data, err := r.Client.HGetAll(ctx, categoryKeyPrefix+id).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get category data: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	category := &entity.Category{
		ID:   id,
		Name: data["name"],
	}

	if parentID, ok := data["parent_id"]; ok && parentID != "" {
		category.ParentID = &parentID
	}

	return category, nil
}

func (r *CategoryRedisRepo) GetByParentID(ctx context.Context, parentID string) (*[]entity.Category, error) {
	childIDs, err := r.Client.SMembers(ctx, categoryParentKey+parentID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get child IDs: %w", err)
	}

	pipe := r.Client.Pipeline()
	categoryDataCmds := make(map[string]*redis.MapStringStringCmd)

	for _, id := range childIDs {
		categoryDataCmds[id] = pipe.HGetAll(ctx, categoryKeyPrefix+id)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get child categories data: %w", err)
	}

	categories := make([]entity.Category, 0, len(childIDs))
	for _, id := range childIDs {
		data := categoryDataCmds[id].Val()
		if len(data) == 0 {
			continue
		}

		category := entity.Category{
			ID:   id,
			Name: data["name"],
		}
		category.ParentID = &parentID

		categories = append(categories, category)
	}

	return &categories, nil
}

func (r *CategoryRedisRepo) Add(ctx context.Context, category *entity.Category) error {
	pipe := r.Client.Pipeline()

	categoryKey := categoryKeyPrefix + category.ID
	categoryData := map[string]interface{}{
		"name": category.Name,
	}

	if category.ParentID != nil {
		categoryData["parent_id"] = *category.ParentID
		pipe.SAdd(ctx, categoryParentKey+*category.ParentID, category.ID)
	}

	pipe.HSet(ctx, categoryKey, categoryData)
	pipe.SAdd(ctx, categorySetKey, category.ID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add category: %w", err)
	}

	return nil
}

func (r *CategoryRedisRepo) Update(ctx context.Context, id string, newName string) error {
	categoryKey := categoryKeyPrefix + id

	pipe := r.Client.Pipeline()
	pipe.HSet(ctx, categoryKey, "name", newName)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update category name: %w", err)
	}

	return nil
}

func (r *CategoryRedisRepo) Delete(ctx context.Context, id string) error {
	// get category data to check for parent
	categoryKey := categoryKeyPrefix + id
	categoryData, err := r.Client.HGetAll(ctx, categoryKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get category data: %w", err)
	}

	pipe := r.Client.Pipeline()

	// remove from parent's children set if parent exists
	if parentID, exists := categoryData["parent_id"]; exists {
		pipe.SRem(ctx, categoryParentKey+parentID, id)
	}

	// remove category data and from main set
	pipe.Del(ctx, categoryKey)
	pipe.SRem(ctx, categorySetKey, id)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
