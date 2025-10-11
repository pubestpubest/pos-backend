package request

type CategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	DisplayOrder *int   `json:"display_order"`
}
