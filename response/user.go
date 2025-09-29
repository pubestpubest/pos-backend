package response

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName *string   `json:"full_name"`
	Email    *string   `json:"email"`
	Phone    *string   `json:"phone"`
	Status   *string   `json:"status"`
}
