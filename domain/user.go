package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// User domain - manages staff/employee users
type UserUsecase interface {
	GetAllUsers() ([]*response.UserResponse, error)
	GetUserByID(id uuid.UUID) (*response.UserResponse, error)
	CreateUser(req *request.UserCreateRequest) (*response.UserResponse, error)
	UpdateUser(id uuid.UUID, req *request.UserUpdateRequest) (*response.UserResponse, error)
	AssignRoleToUser(userID uuid.UUID, roleID int) error
}

type UserRepository interface {
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetUserWithRoles(id uuid.UUID) (*models.User, error)
	AssignRole(userRole *models.UserRole) error
}
