package response

import "github.com/google/uuid"

type MenuItemResponse struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	SKU       string           `json:"sku"`
	PriceBaht int64            `json:"price_baht"`
	Active    bool             `json:"active"`
	ImageURL  string           `json:"image_url"`
	Category  CategoryResponse `json:"category"`
}
