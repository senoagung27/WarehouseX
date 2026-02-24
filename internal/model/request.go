package model

import (
	"time"

	"github.com/google/uuid"
)

// Request types
const (
	RequestTypeInbound  = "INBOUND"
	RequestTypeOutbound = "OUTBOUND"
)

// Request statuses (state machine)
const (
	StatusPending   = "PENDING"
	StatusApproved  = "APPROVED"
	StatusRejected  = "REJECTED"
	StatusCompleted = "COMPLETED"
)

type Request struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Type       string     `gorm:"size:20;not null" json:"type"`
	Status     string     `gorm:"size:20;not null;default:'PENDING'" json:"status"`
	ItemID     uuid.UUID  `gorm:"type:uuid;not null" json:"item_id"`
	Quantity   int        `gorm:"not null" json:"quantity"`
	Notes      string     `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy  uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	ApprovedBy *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations (for preloading)
	Item     Inventory `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	Creator  User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Approver *User     `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

func (Request) TableName() string {
	return "requests"
}

// ValidTransition checks if state transition is allowed
func ValidTransition(from, to string) bool {
	transitions := map[string][]string{
		StatusPending:  {StatusApproved, StatusRejected},
		StatusApproved: {StatusCompleted},
	}
	allowed, ok := transitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
