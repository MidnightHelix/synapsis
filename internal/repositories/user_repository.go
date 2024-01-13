// repository/user_repository.go

package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/domain"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(username, password string) error
	Login(username, password string) (bool, error)
}

type userRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepository {
	return &userRepository{db, redis}
}

func (r *userRepository) Register(username, password string) error {
	user := &domain.User{Username: username, Password: password}
	result := r.db.Create(user)

	// Invalidate Redis cache for this user
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", username)
	_ = r.redis.Del(ctx, key)

	return result.Error
}

func (r *userRepository) Login(username, password string) (bool, error) {
	// Check if user credentials are available in Redis cache
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", username)
	val, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		// Cache miss, fetch user credentials from the database
		var user domain.User
		result := r.db.Where("username = ? AND password = ?", username, password).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// User not found
				return false, nil
			}
			return false, result.Error
		}

		// Set user credentials in Redis cache for future use
		val, err := json.Marshal(user)
		if err != nil {
			return false, err
		}
		err = r.redis.Set(ctx, key, val, 0).Err()
		if err != nil {
			return false, err
		}

		return true, nil
	} else if err != nil {
		// Redis error
		return false, err
	}

	// Cache hit, return user credentials from Redis
	var user domain.User
	err = json.Unmarshal([]byte(val), &user)
	return true, err
}
