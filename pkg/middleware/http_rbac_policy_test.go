package middleware

import (
	"errors"
	"testing"
)

func TestStrictRBAC(t *testing.T) {
	if err := EnforceAction(RoleBuyer, ActionCreateProduct, "buyer1", "seller1"); !errors.Is(err, ErrRBACDenied) {
		t.Fatalf("buyer must not create product")
	}

	if err := EnforceAction(RoleSeller, ActionUpdateProduct, "seller1", "seller2"); !errors.Is(err, ErrRBACDenied) {
		t.Fatalf("seller must not mutate other seller's product")
	}

	if err := EnforceAction(RoleAdmin, ActionModerate, "admin1", "seller2"); err != nil {
		t.Fatalf("only admin can moderate")
	}
}
