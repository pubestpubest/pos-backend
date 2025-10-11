package response

import (
	"time"

	"github.com/google/uuid"
)

type OrderResponse struct {
	ID           uuid.UUID           `json:"id"`
	TableID      uuid.UUID           `json:"table_id"`
	TableName    string              `json:"table_name"`
	OpenedBy     uuid.UUID           `json:"opened_by"`
	Source       string              `json:"source"`
	Status       string              `json:"status"`
	SubtotalBaht int64               `json:"subtotal_baht"`
	DiscountBaht int64               `json:"discount_baht"`
	TotalBaht    int64               `json:"total_baht"`
	Note         string              `json:"note"`
	CreatedAt    time.Time           `json:"created_at"`
	ClosedAt     *time.Time          `json:"closed_at"`
	Items        []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ID            uuid.UUID                   `json:"id"`
	MenuItemID    uuid.UUID                   `json:"menu_item_id"`
	MenuItemName  string                      `json:"menu_item_name"`
	Quantity      int                         `json:"quantity"`
	UnitPriceBaht int64                       `json:"unit_price_baht"`
	LineTotalBaht int64                       `json:"line_total_baht"`
	Note          string                      `json:"note"`
	Modifiers     []OrderItemModifierResponse `json:"modifiers"`
}

type OrderItemModifierResponse struct {
	ModifierID     uuid.UUID `json:"modifier_id"`
	ModifierName   string    `json:"modifier_name"`
	PriceDeltaBaht int64     `json:"price_delta_baht"`
}
