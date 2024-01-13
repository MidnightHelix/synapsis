// repository/product_repository.go

package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MidnightHelix/synapsis/domain"
	"github.com/go-redis/redis/v9"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductsByCategory(category string) ([]*domain.Product, error)
}

type productRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProductRepository(db *gorm.DB, redis *redis.Client) ProductRepository {
	return &productRepository{db, redis}
}

func (r *productRepository) GetProductsByCategory(category string) ([]*domain.Product, error) {
	// Check if products are available in Redis cache
	ctx := context.Background()
	key := fmt.Sprintf("products:%s", category)
	val, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		// Cache miss, fetch products from the database
		var products []*domain.Product
		result := r.db.Where("category = ?", category).Find(&products)
		if result.Error != nil {
			return nil, result.Error
		}

		// Set products in Redis cache for future use
		val, err := json.Marshal(products)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(ctx, key, val, 0).Err()
		if err != nil {
			return nil, err
		}

		return products, nil
	} else if err != nil {
		// Redis error
		return nil, err
	}

	// Cache hit, return products from Redis
	var products []*domain.Product
	err = json.Unmarshal([]byte(val), &products)
	return products, err
}
