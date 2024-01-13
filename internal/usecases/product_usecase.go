package service

import "github.com/MidnightHelix/synapsis/domain"

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
