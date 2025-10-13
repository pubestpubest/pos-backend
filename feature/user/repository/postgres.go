package repository

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Preload("Roles.Permissions").Order("username ASC").Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "[UserRepository.GetAllUsers]: Error getting users")
	}
	return users, nil
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Roles.Permissions").Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[UserRepository.GetUserByID]: User not found")
		}
		return nil, errors.Wrap(err, "[UserRepository.GetUserByID]: Error querying database")
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.Wrap(err, "[UserRepository.CreateUser]: Error creating user")
	}
	return nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.Wrap(err, "[UserRepository.UpdateUser]: Error updating user")
	}
	return nil
}

func (r *userRepository) GetUserWithRoles(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Roles.Permissions").Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[UserRepository.GetUserWithRoles]: User not found")
		}
		return nil, errors.Wrap(err, "[UserRepository.GetUserWithRoles]: Error querying database")
	}
	return &user, nil
}

func (r *userRepository) AssignRole(userRole *models.UserRole) error {
	if err := r.db.Create(userRole).Error; err != nil {
		return errors.Wrap(err, "[UserRepository.AssignRole]: Error assigning role")
	}
	return nil
}
