package handler

import (
	"net/http"
	"your_project/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	GetProductsByCategory(c *gin.Context)
}

type productHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) ProductHandler {
	return &productHandler{service}
}

func (h *productHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	products, err := h.service.GetProductsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, products)
}
