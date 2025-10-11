package request

type ModifierRequest struct {
	Name           string  `json:"name" binding:"required"`
	PriceDeltaBaht *int64  `json:"price_delta_baht"`
	Note           *string `json:"note"`
}
