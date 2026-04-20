package catalog

import (
	"errors"
	"fmt"
)

var (
	ErrForbidden          = errors.New("forbidden")
	ErrProductUnavailable = errors.New("product unavailable")
)

type ProductStatus string

const (
	ProductStatusPublished  ProductStatus = "published"
	ProductStatusHidden     ProductStatus = "hidden"
	ProductStatusArchived   ProductStatus = "archived"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
)

type Product struct {
	ID         string
	SellerID   string
	Title      string
	PriceCents int64
	Status     ProductStatus
	Stock      uint32
}

func (p Product) IsAvailable() bool {
	return p.Status == ProductStatusPublished && p.Stock > 0
}

type Service struct{}

func (s Service) ValidateForPurchase(product Product) error {
	if !product.IsAvailable() {
		return fmt.Errorf("%w: stauts=%s stock=%d", ErrProductUnavailable, product.Status, product.Stock)
	}
	return nil
}

func (s Service) AuthorizeCreateProduct(role string) error {
	if role != "seller" && role != "admin" {
		return ErrForbidden
	}
	return nil
}

func (s Service) AuthorizeMutateProduct(role, actorID string, product Product) error {
	if role == "admin" {
		return nil
	}

	if role != "seller" || actorID != product.SellerID {
		return ErrForbidden
	}

	return nil
}
