package response

import "github.com/google/uuid"

type AreaResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
