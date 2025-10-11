package repository

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[AuthRepository.GetUserByUsername]: User not found")
		}
		return nil, errors.Wrap(err, "[AuthRepository.GetUserByUsername]: Error querying database")
	}
	return &user, nil
}

func (r *authRepository) GetUserWithRolesAndPermissions(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Roles.Permissions").Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[AuthRepository.GetUserWithRolesAndPermissions]: User not found")
		}
		return nil, errors.Wrap(err, "[AuthRepository.GetUserWithRolesAndPermissions]: Error querying database")
	}
	return &user, nil
}

func (r *authRepository) GetUserPermissions(userID uuid.UUID) ([]string, error) {
	var permissions []string

	err := r.db.Table("permissions").
		Select("DISTINCT permissions.code").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permissions.code", &permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "[AuthRepository.GetUserPermissions]: Error querying database")
	}

	return permissions, nil
}

func (r *authRepository) UpdatePassword(userID uuid.UUID, passwordHash string) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error; err != nil {
		return errors.Wrap(err, "[AuthRepository.UpdatePassword]: Error updating password")
	}
	return nil
}

func (r *authRepository) CreateSession(session *models.Session) error {
	if err := r.db.Create(session).Error; err != nil {
		return errors.Wrap(err, "[AuthRepository.CreateSession]: Error creating session")
	}
	return nil
}

func (r *authRepository) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	if err := r.db.Preload("User").Where("token = ?", token).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[AuthRepository.GetSessionByToken]: Session not found")
		}
		return nil, errors.Wrap(err, "[AuthRepository.GetSessionByToken]: Error querying database")
	}
	return &session, nil
}

func (r *authRepository) DeleteSession(token string) error {
	if err := r.db.Where("token = ?", token).Delete(&models.Session{}).Error; err != nil {
		return errors.Wrap(err, "[AuthRepository.DeleteSession]: Error deleting session")
	}
	return nil
}

func (r *authRepository) CleanupExpiredSessions() error {
	if err := r.db.Where("expires_at < NOW()").Delete(&models.Session{}).Error; err != nil {
		return errors.Wrap(err, "[AuthRepository.CleanupExpiredSessions]: Error cleaning up sessions")
	}
	return nil
}
