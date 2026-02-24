package model

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ItemName  string    `gorm:"size:255;not null" json:"item_name"`
	SKU       string    `gorm:"size:100;uniqueIndex" json:"sku"`
	Quantity  int       `gorm:"not null;default:0" json:"quantity"`
	Unit      string    `gorm:"size:50;not null;default:'pcs'" json:"unit"`
	Version   int       `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Inventory) TableName() string {
	return "inventory"
}
