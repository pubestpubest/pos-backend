package response

import "github.com/google/uuid"

type ModifierResponse struct {
	ID             uuid.UUID         `json:"id"`
	CategoryID     *uuid.UUID        `json:"category_id,omitempty"`
	Name           string            `json:"name"`
	PriceDeltaBaht int64             `json:"price_delta_baht"`
	Note           string            `json:"note"`
	Category       *CategoryResponse `json:"category,omitempty"`
}
