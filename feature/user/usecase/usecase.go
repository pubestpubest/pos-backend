package usecase

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/go-clean-arch-template/constant"
	"github.com/pubestpubest/go-clean-arch-template/domain"
	"github.com/pubestpubest/go-clean-arch-template/response"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepository: userRepository}
}

func (u *userUsecase) GetUser(id uint32) (*response.UserResponse, error) {
	user, err := u.userRepository.GetUser(id)
	if err != nil {
		err = errors.Wrap(err, "[UserUsecase.GetUser]: Error getting user")
		return nil, err
	}

	if user == nil {
		err = errors.New(constant.UserNotFound)
		return nil, err
	}

	return &response.UserResponse{
		ID:        uint32(user.ID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Age:       user.Age,
	}, nil
}
