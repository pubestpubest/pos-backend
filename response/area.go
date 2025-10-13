package response

import "github.com/google/uuid"

type AreaResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type AreaWithTablesResponse struct {
	ID     uuid.UUID       `json:"id"`
	Name   string          `json:"name"`
	Tables []TableResponse `json:"tables"`
}
