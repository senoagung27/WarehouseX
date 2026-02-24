package dto

type CreateRequestInput struct {
	ItemID   string `json:"item_id" binding:"required,uuid"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Notes    string `json:"notes"`
}
