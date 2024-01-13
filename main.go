// main.go

package main

import (
	"github.com/MidnightHelix/synapsis/repository"

	"github.com/MidnightHelix/synapsis/handler"

	"github.com/MidnightHelix/synapsis/domain"
	"github.com/MidnightHelix/synapsis/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("synapsis_store.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&domain.Product{}, &domain.Cart{}, &domain.User{})

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server details
		Password: "",               // No password
		DB:       0,                // Default DB
	})

	productRepo := repository.NewProductRepository(db, redisClient)
	cartRepo := repository.NewCartRepository(db, redisClient)
	userRepo := repository.NewUserRepository(db, redisClient)

	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo)
	userService := service.NewUserService(userRepo)

	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/products/:category", productHandler.GetProductsByCategory)
		v1.POST("/cart/add", cartHandler.AddToCart)
		v1.GET("/cart", cartHandler.GetCart)
		v1.DELETE("/cart/:productID", cartHandler.DeleteFromCart)
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
	}

	router.Run(":8080")
}
