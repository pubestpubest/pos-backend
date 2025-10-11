package request

import "github.com/google/uuid"

type PaymentRequest struct {
	OrderID     uuid.UUID `json:"order_id" binding:"required"`
	Method      string    `json:"method" binding:"required,oneof=cash card promptpay"`
	AmountBaht  int64     `json:"amount_baht" binding:"required,min=1"`
	Provider    *string   `json:"provider"`
	ProviderRef *string   `json:"provider_ref"`
}
