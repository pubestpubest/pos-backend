package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	authRepository domain.AuthRepository
}

func NewAuthUsecase(authRepository domain.AuthRepository) domain.AuthUsecase {
	return &authUsecase{authRepository: authRepository}
}

func (u *authUsecase) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	// Get user by username
	user, err := u.authRepository.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "[AuthUsecase.Login]: Invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.Wrap(err, "[AuthUsecase.Login]: Invalid username or password")
	}

	// Check user status
	if user.Status != nil && *user.Status == "locked" {
		return nil, errors.New("[AuthUsecase.Login]: User account is locked")
	}

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, errors.Wrap(err, "[AuthUsecase.Login]: Error generating token")
	}

	// Create session (24 hours expiry)
	expiresAt := time.Now().Add(24 * time.Hour)
	session := &models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := u.authRepository.CreateSession(session); err != nil {
		return nil, errors.Wrap(err, "[AuthUsecase.Login]: Error creating session")
	}

	// Get user permissions
	permissions, err := u.authRepository.GetUserPermissions(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "[AuthUsecase.Login]: Error getting permissions")
	}

	// Build response
	authResponse := &response.AuthResponse{
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
			Phone:    user.Phone,
			Status:   user.Status,
		},
		Token:       token,
		ExpiresAt:   expiresAt,
		Permissions: permissions,
	}

	return authResponse, nil
}

func (u *authUsecase) Logout(token string) error {
	if err := u.authRepository.DeleteSession(token); err != nil {
		return errors.Wrap(err, "[AuthUsecase.Logout]: Error deleting session")
	}
	return nil
}

func (u *authUsecase) ChangePassword(userID uuid.UUID, req *request.ChangePasswordRequest) error {
	// Get user
	user, err := u.authRepository.GetUserWithRolesAndPermissions(userID)
	if err != nil {
		return errors.Wrap(err, "[AuthUsecase.ChangePassword]: Error getting user")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("[AuthUsecase.ChangePassword]: Old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "[AuthUsecase.ChangePassword]: Error hashing password")
	}

	// Update password
	if err := u.authRepository.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return errors.Wrap(err, "[AuthUsecase.ChangePassword]: Error updating password")
	}

	return nil
}

func (u *authUsecase) VerifyPermission(userID uuid.UUID, permissionCode string) (bool, error) {
	permissions, err := u.authRepository.GetUserPermissions(userID)
	if err != nil {
		return false, errors.Wrap(err, "[AuthUsecase.VerifyPermission]: Error getting permissions")
	}

	for _, p := range permissions {
		if p == permissionCode {
			return true, nil
		}
	}

	return false, nil
}

func (u *authUsecase) GetUserPermissions(userID uuid.UUID) ([]string, error) {
	permissions, err := u.authRepository.GetUserPermissions(userID)
	if err != nil {
		return nil, errors.Wrap(err, "[AuthUsecase.GetUserPermissions]: Error getting permissions")
	}
	return permissions, nil
}

func (u *authUsecase) GetUserByToken(token string) (*models.User, error) {
	session, err := u.authRepository.GetSessionByToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "[AuthUsecase.GetUserByToken]: Invalid or expired token")
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("[AuthUsecase.GetUserByToken]: Session expired")
	}

	// Get user with full details
	user, err := u.authRepository.GetUserWithRolesAndPermissions(session.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "[AuthUsecase.GetUserByToken]: Error getting user")
	}

	return user, nil
}

// generateToken creates a random token for session
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
