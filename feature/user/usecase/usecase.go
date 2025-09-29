package usecase

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/response"
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
		err = errors.Wrap(err, "[UserUsecase.GetUser]: Error getting user")
		return nil, err
	}
	userResponses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
			Phone:    user.Phone,
			Status:   user.Status,
		}
	}

	return userResponses, nil
}
