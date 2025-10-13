package request

import "github.com/google/uuid"

type OrderCreateRequest struct {
	TableID  uuid.UUID  `json:"table_id" binding:"required"`
	OpenedBy *uuid.UUID `json:"opened_by"`
	Source   string     `json:"source" binding:"required,oneof=staff customer"`
	Note     *string    `json:"note"`
}

type AddOrderItemRequest struct {
	MenuItemID  uuid.UUID   `json:"menu_item_id" binding:"required"`
	Quantity    int         `json:"quantity" binding:"required,min=1"`
	Note        *string     `json:"note"`
	ModifierIDs []uuid.UUID `json:"modifier_ids"`
}

type UpdateOrderItemQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
