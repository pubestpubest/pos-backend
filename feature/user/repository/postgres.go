package repository

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/go-clean-arch-template/domain"
	"github.com/pubestpubest/go-clean-arch-template/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(id uint32) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		err = errors.Wrap(err, "[UserRepository.GetUser]: Error getting user")
		return nil, err
	}

	return &user, nil
}
