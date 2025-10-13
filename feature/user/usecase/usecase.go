package usecase

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
	"github.com/pubestpubest/pos-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepository: userRepository}
}

func (u *userUsecase) GetAllUsers() ([]*response.UserResponse, error) {
	users, err := u.userRepository.GetAllUsers()
	if err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.GetAllUsers]: Error getting users")
	}

	userResponses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = u.buildUserResponse(user)
	}

	return userResponses, nil
}

func (u *userUsecase) GetUserByID(id uuid.UUID) (*response.UserResponse, error) {
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.GetUserByID]: Error getting user")
	}

	return u.buildUserResponse(user), nil
}

func (u *userUsecase) CreateUser(req *request.UserCreateRequest) (*response.UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.CreateUser]: Error hashing password")
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       req.Status,
	}

	if err := u.userRepository.CreateUser(user); err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.CreateUser]: Error creating user")
	}

	// Get the created user with roles
	createdUser, err := u.userRepository.GetUserByID(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.CreateUser]: Error retrieving created user")
	}

	return u.buildUserResponse(createdUser), nil
}

func (u *userUsecase) UpdateUser(id uuid.UUID, req *request.UserUpdateRequest) (*response.UserResponse, error) {
	// Get existing user
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.UpdateUser]: User not found")
	}

	// Update fields
	if req.FullName != nil {
		user.FullName = req.FullName
	}
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Status != nil {
		user.Status = req.Status
	}

	if err := u.userRepository.UpdateUser(user); err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.UpdateUser]: Error updating user")
	}

	// Get the updated user with roles
	updatedUser, err := u.userRepository.GetUserByID(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.UpdateUser]: Error retrieving updated user")
	}

	return u.buildUserResponse(updatedUser), nil
}

func (u *userUsecase) AssignRoleToUser(userID uuid.UUID, roleID int) error {
	// Check if user exists
	_, err := u.userRepository.GetUserByID(userID)
	if err != nil {
		return errors.Wrap(err, "[UserUsecase.AssignRoleToUser]: User not found")
	}

	userRole := &models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := u.userRepository.AssignRole(userRole); err != nil {
		return errors.Wrap(err, "[UserUsecase.AssignRoleToUser]: Error assigning role")
	}

	return nil
}

// Helper function to build role responses
func (u *userUsecase) buildRoleResponses(roles []models.Role) []response.RoleResponse {
	roleResponses := make([]response.RoleResponse, len(roles))
	for i, role := range roles {
		permissions := make([]response.PermissionResponse, len(role.Permissions))
		for j, permission := range role.Permissions {
			permissions[j] = response.PermissionResponse{
				ID:          permission.ID,
				Code:        permission.Code,
				Description: utils.DerefString(permission.Description),
			}
		}
		roleResponses[i] = response.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissions,
		}
	}
	return roleResponses
}

// Helper function to build user response
func (u *userUsecase) buildUserResponse(user *models.User) *response.UserResponse {
	return &response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    user.Phone,
		Status:   user.Status,
		Roles:    u.buildRoleResponses(user.Roles),
	}
}
