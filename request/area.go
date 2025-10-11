package request

type AreaRequest struct {
	Name string `json:"name" binding:"required"`
}
