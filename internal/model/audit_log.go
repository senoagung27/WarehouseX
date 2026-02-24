package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Entity      string          `gorm:"size:100;not null" json:"entity"`
	EntityID    uuid.UUID       `gorm:"type:uuid;not null" json:"entity_id"`
	Action      string          `gorm:"size:50;not null" json:"action"`
	UserID      uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	BeforeValue json.RawMessage `gorm:"type:jsonb" json:"before_value,omitempty"`
	AfterValue  json.RawMessage `gorm:"type:jsonb" json:"after_value,omitempty"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`

	// Relation
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
