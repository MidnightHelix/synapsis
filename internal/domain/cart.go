// domain/cart.go

package domain

type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	ProductID uint
	Quantity  int
}
