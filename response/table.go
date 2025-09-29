package response

import "github.com/google/uuid"

type TableResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Seats  int       `json:"seats"`
	Status string    `json:"status"`
	QRCode string    `json:"qr_code"`
	AreaID uuid.UUID `json:"area_id"`
}
