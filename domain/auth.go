package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// Auth domain - manages authentication and authorization
type AuthUsecase interface {
	Login(req *request.LoginRequest) (*response.AuthResponse, error)
	Logout(token string) error
	ChangePassword(userID uuid.UUID, req *request.ChangePasswordRequest) error
	VerifyPermission(userID uuid.UUID, permissionCode string) (bool, error)
	GetUserPermissions(userID uuid.UUID) ([]string, error)
	GetUserByToken(token string) (*models.User, error)
}

type AuthRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	GetUserWithRolesAndPermissions(id uuid.UUID) (*models.User, error)
	GetUserPermissions(userID uuid.UUID) ([]string, error)
	UpdatePassword(userID uuid.UUID, passwordHash string) error
	CreateSession(session *models.Session) error
	GetSessionByToken(token string) (*models.Session, error)
	DeleteSession(token string) error
	CleanupExpiredSessions() error
}
