// service/product_service.go

package usecase

import "your_project/domain"

type ProductService interface {
	GetProductsByCategory(category string) ([]*domain.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) GetProductsByCategory(category string) ([]*domain.Product, error) {
	return s.repo.GetProductsByCategory(category)
}
