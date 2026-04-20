package middleware

import "errors"

var ErrRBACDenied = errors.New("rbac denied")

type Action string
type ActionRole string

const (
	RoleBuyer  ActionRole = "buyer"
	RoleSeller ActionRole = "seller"
	RoleAdmin  ActionRole = "admin"
)

const (
	ActionCreateProduct Action = "create_product"
	ActionUpdateProduct Action = "update_product"
	ActionModerate      Action = "moderate"
)

func EnforceAction(role ActionRole, action Action, actorID, resourceOwnerID string) error {
	switch role {
	case RoleAdmin:
		return nil
	case RoleBuyer:
		if action == ActionCreateProduct || action == ActionUpdateProduct || action == ActionModerate {
			return ErrRBACDenied
		}
		return nil
	case RoleSeller:
		if action == ActionModerate {
			return ErrRBACDenied
		}

		if (action == ActionCreateProduct || action == ActionUpdateProduct) && actorID != resourceOwnerID {
			return ErrRBACDenied
		}
		return nil
	default:
		return ErrRBACDenied
	}
}
