package cart

import (
	"errors"
	"testing"
)

func TestAddItemQuantityBounds(t *testing.T) {
	svc := NewService()

	if err := svc.AddItem("u1", "p1", 0, 5); !errors.Is(err, ErrInvalidQuantity) {
		t.Fatalf("expected ErrInvalidQuality, got %v", err)
	}

	if err := svc.AddItem("u1", "p1", 2, 1); !errors.Is(err, ErrOutOfStock) {
		t.Fatalf("expected ErrOutOfStock, got %v", err)
	}

	if err := svc.AddItem("u1", "p1", 2, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSingleActiveCartPerUser(t *testing.T) {
	svc := NewService()
	first := svc.GetOrCreateActiveCart("u1")
	second := svc.GetOrCreateActiveCart("u1")

	if first != second {
		t.Fatalf("expected same active cart instance for user")
	}
}
