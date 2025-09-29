package domain

import (
	"github.com/pubestpubest/go-clean-arch-template/models"
	"github.com/pubestpubest/go-clean-arch-template/response"
)

type UserUsecase interface {
	GetUser(id uint32) (*response.UserResponse, error)
}

type UserRepository interface {
	GetUser(id uint32) (*models.User, error)
}
