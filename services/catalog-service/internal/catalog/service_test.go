package catalog

import (
	"errors"
	"testing"
)

func TestValidateForPurchase(t *testing.T) {
	svc := Service{}

	tests := []struct {
		name    string
		product Product
		wantErr bool
	}{
		{"published", Product{Status: ProductStatusPublished, Stock: 5}, false},
		{"hidden", Product{Status: ProductStatusHidden, Stock: 5}, true},
		{"archived", Product{Status: ProductStatusArchived, Stock: 5}, true},
		{"out_of_stock_status", Product{Status: ProductStatusOutOfStock, Stock: 5}, true},
		{"zero_stock", Product{Status: ProductStatusPublished, Stock: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateForPurchase(tt.product)
			if tt.wantErr && !errors.Is(err, ErrProductUnavailable) {
				t.Fatalf("expected ErrProductUnavailable, got %v", err)
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAuthorizeCreateProduct(t *testing.T) {
	svc := Service{}
	if err := svc.AuthorizeCreateProduct("buyer"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("buyer should be forbidden: %v", err)
	}

	if err := svc.AuthorizeCreateProduct("seller"); err != nil {
		t.Fatalf("seller should be allowed: %v", err)
	}
}

func TestAuthorizeMutateProduct(t *testing.T) {
	svc := Service{}
	product := Product{SellerID: "seller1"}

	if err := svc.AuthorizeMutateProduct("seller", "seller2", product); !errors.Is(err, ErrForbidden) {
		t.Fatalf("different seller should be forbidden")
	}

	if err := svc.AuthorizeMutateProduct("admin", "admin1", product); err != nil {
		t.Fatalf("admin should be allowed")
	}
}
