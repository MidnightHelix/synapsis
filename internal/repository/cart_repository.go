package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MidnightHelix/synapsis/domain"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CartRepository interface {
	AddToCart(userID, productID uint, quantity int) error
	GetCart(userID uint) ([]*domain.Product, error)
	DeleteFromCart(userID, productID uint) error
}

type cartRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCartRepository(db *gorm.DB, redis *redis.Client) CartRepository {
	return &cartRepository{db, redis}
}

func (r *cartRepository) AddToCart(userID, productID uint, quantity int) error {
	// Add the product to the cart in the database
	// ...

	// Update or create the cart in Redis
	ctx := context.Background()
	key := fmt.Sprintf("cart:%d", userID)

	cart, err := r.GetCartFromRedis(ctx, key)
	if err != nil {
		return err
	}

	// Check if the product is already in the cart
	var found bool
	for _, item := range cart {
		if item.ID == productID {
			item.Quantity += quantity
			found = true
			break
		}
	}

	// If the product is not in the cart, add it
	if !found {
		cart = append(cart, &domain.Product{ID: productID, Quantity: quantity})
	}

	return r.SetCartInRedis(ctx, key, cart)
}

func (r *cartRepository) GetCart(userID uint) ([]*domain.Product, error) {
	// Retrieve the cart from the database
	// ...

	// Retrieve the cart from Redis
	ctx := context.Background()
	key := fmt.Sprintf("cart:%d", userID)
	return r.GetCartFromRedis(ctx, key)
}

func (r *cartRepository) DeleteFromCart(userID, productID uint) error {
	// Remove the product from the cart in the database
	// ...

	// Remove the product from the cart in Redis
	ctx := context.Background()
	key := fmt.Sprintf("cart:%d", userID)

	cart, err := r.GetCartFromRedis(ctx, key)
	if err != nil {
		return err
	}

	// Filter out the deleted product
	var updatedCart []*domain.Product
	for _, item := range cart {
		if item.ID != productID {
			updatedCart = append(updatedCart, item)
		}
	}

	return r.SetCartInRedis(ctx, key, updatedCart)
}

func (r *cartRepository) GetCartFromRedis(ctx context.Context, key string) ([]*domain.Product, error) {
	val, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var cart []*domain.Product
	err = json.Unmarshal([]byte(val), &cart)
	return cart, err
}

func (r *cartRepository) SetCartInRedis(ctx context.Context, key string, cart []*domain.Product) error {
	val, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, key, val, 0).Err()
}
