package request

import "github.com/google/uuid"

type MenuItemRequest struct {
	CategoryID *uuid.UUID `json:"category_id"`
	Name       string     `json:"name" binding:"required"`
	SKU        string     `json:"sku" binding:"required"`
	PriceBaht  int64      `json:"price_baht" binding:"required"`
	Active     *bool      `json:"active"`
	ImageURL   *string    `json:"image_url"`
}
