package request

type UserCreateRequest struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required,min=6"`
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Status   *string `json:"status"`
}

type UserUpdateRequest struct {
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Status   *string `json:"status"`
}

type AssignRoleRequest struct {
	RoleID int `json:"role_id" binding:"required"`
}
