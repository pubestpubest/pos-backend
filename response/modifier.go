package response

import "github.com/google/uuid"

type ModifierResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	PriceDeltaBaht int64     `json:"price_delta_baht"`
	Note           string    `json:"note"`
}
