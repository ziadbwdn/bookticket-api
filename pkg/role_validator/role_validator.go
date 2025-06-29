package role

import (
	"root-app/internal/entities"
)

// IsValid checks if a role is one of the predefined valid roles.
func IsValid(role entities.UserRole) bool {
	switch role {
	case entities.RoleAdmin, entities.RoleUser:
		return true
	default:
		return false
	}
}

// CanUpdateDrillingStatus checks if a role has permission to update drilling status.
func UserFeatureValidation(role entities.UserRole) bool {
	return role == entities.RoleUser
}

// CanManageUsers checks if a role has permission to manage users.
func CanManageFullAccess(role entities.UserRole) bool {
	return role == entities.RoleAdmin
}
