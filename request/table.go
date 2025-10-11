package request

type UpdateTableStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=free occupied needs_pay"`
}
