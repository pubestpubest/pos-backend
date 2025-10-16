package request

import "github.com/google/uuid"

type ModifierRequest struct {
	Name           string     `json:"name" binding:"required"`
	CategoryID     *uuid.UUID `json:"category_id"`
	PriceDeltaBaht *int64     `json:"price_delta_baht"`
	Note           *string    `json:"note"`
}
