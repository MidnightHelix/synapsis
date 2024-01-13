package usecase

import "your_project/domain"

type CartService interface {
	AddToCart(userID, productID uint, quantity int) error
	GetCart(userID uint) ([]*domain.Product, error)
	DeleteFromCart(userID, productID uint) error
}

type cartService struct {
	repo CartRepository
}

func NewCartService(repo CartRepository) CartService {
	return &cartService{repo}
}

func (s *cartService) AddToCart(userID, productID uint, quantity int) error {
	return s.repo.AddToCart(userID, productID, quantity)
}

func (s *cartService) GetCart(userID uint) ([]*domain.Product, error) {
	return s.repo.GetCart(userID)
}

func (s *cartService) DeleteFromCart(userID, productID uint) error {
	return s.repo.DeleteFromCart(userID, productID)
}
