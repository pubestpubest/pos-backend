package response

type RoleResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
}

type PermissionResponse struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}
