package cart

import (
	"errors"
)

var (
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrOutOfStock      = errors.New("out of stock")
)

type Item struct {
	ProductID string
	Quantity  uint32
}

type ActiveCart struct {
	UserID string
	Items  map[string]Item
}

type Service struct {
	carts map[string]*ActiveCart
}

func NewService() *Service {
	return &Service{carts: map[string]*ActiveCart{}}
}

func (s *Service) GetOrCreateActiveCart(userID string) *ActiveCart {
	if cart, ok := s.carts[userID]; ok {
		return cart
	}

	cart := &ActiveCart{UserID: userID, Items: map[string]Item{}}
	s.carts[userID] = cart

	return cart
}

func (s *Service) AddItem(userID, productID string, quantity, stock uint32) error {
	if quantity < 1 {
		return ErrInvalidQuantity
	}

	if quantity > stock {
		return ErrOutOfStock
	}

	cart := s.GetOrCreateActiveCart(userID)
	item := Item{ProductID: productID, Quantity: quantity}
	cart.Items[productID] = item

	return nil
}
