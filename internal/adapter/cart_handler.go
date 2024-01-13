package handler

import (
	"net/http"

	"github.com/MidnightHelix/synapsis/service"

	"github.com/gin-gonic/gin"
)

type CartHandler interface {
	AddToCart(c *gin.Context)
	GetCart(c *gin.Context)
	DeleteFromCart(c *gin.Context)
}

type cartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) CartHandler {
	return &cartHandler{service}
}

func (h *cartHandler) AddToCart(c *gin.Context) {
	var request struct {
		UserID    uint `json:"user_id" binding:"required"`
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.service.AddToCart(request.UserID, request.ProductID, request.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}

func (h *cartHandler) GetCart(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	cart, err := h.service.GetCart(request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *cartHandler) DeleteFromCart(c *gin.Context) {
	var request struct {
		UserID    uint `json:"user_id" binding:"required"`
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.service.DeleteFromCart(request.UserID, request.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart successfully"})
}
