package dto

type CreateInventoryInput struct {
	ItemName string `json:"item_name" binding:"required"`
	SKU      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"min=0"`
	Unit     string `json:"unit" binding:"required"`
}

type UpdateInventoryInput struct {
	ItemName string `json:"item_name"`
	SKU      string `json:"sku"`
	Unit     string `json:"unit"`
}
