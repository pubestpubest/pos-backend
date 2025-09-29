package domain

import (
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/response"
)

type UserUsecase interface {
	GetAllUsers() ([]*response.UserResponse, error)
}

type UserRepository interface {
	GetAllUsers() ([]*models.User, error)
}
