package response

import (
	"time"

	"github.com/google/uuid"
)

type PaymentResponse struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	Method      string    `json:"method"`
	AmountBaht  int64     `json:"amount_baht"`
	Currency    string    `json:"currency"`
	Provider    string    `json:"provider"`
	ProviderRef string    `json:"provider_ref"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaymentMethodResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
