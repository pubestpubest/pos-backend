package response

import "github.com/google/uuid"

type TableResponse struct {
	ID     uuid.UUID     `json:"id"`
	Name   string        `json:"name"`
	Seats  int           `json:"seats"`
	Status string        `json:"status"`
	QRCode string        `json:"qr_code"`
	Area   *AreaResponse `json:"area"`
}

type TableWithOrdersResponse struct {
	ID     uuid.UUID       `json:"id"`
	Name   string          `json:"name"`
	Status string          `json:"status"`
	QRCode string          `json:"qr_code"`
	Orders []OrderResponse `json:"orders"`
}
