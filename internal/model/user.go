package model

import (
	"time"

	"github.com/google/uuid"
)

// User roles
const (
	RoleStaff      = "staff"
	RoleSupervisor = "supervisor"
	RoleAdmin      = "admin"
	RoleAuditor    = "auditor"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:50;not null;default:'staff'" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// ValidRoles returns all valid roles
func ValidRoles() []string {
	return []string{RoleStaff, RoleSupervisor, RoleAdmin, RoleAuditor}
}

// CanApprove checks if the role has approval permission
func CanApprove(role string) bool {
	return role == RoleSupervisor || role == RoleAdmin
}
