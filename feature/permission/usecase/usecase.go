package usecase

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/response"
	"github.com/pubestpubest/pos-backend/utils"
)

type permissionUsecase struct {
	permissionRepository domain.PermissionRepository
}

func NewPermissionUsecase(permissionRepository domain.PermissionRepository) domain.PermissionUsecase {
	return &permissionUsecase{permissionRepository: permissionRepository}
}

func (u *permissionUsecase) GetAllPermissions() ([]*response.PermissionResponse, error) {
	permissions, err := u.permissionRepository.GetAllPermissions()
	if err != nil {
		return nil, errors.Wrap(err, "[PermissionUsecase.GetAllPermissions]: Error getting permissions")
	}

	permissionResponses := make([]*response.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionResponses[i] = &response.PermissionResponse{
			ID:          permission.ID,
			Code:        permission.Code,
			Description: utils.DerefString(permission.Description),
		}
	}

	return permissionResponses, nil
}
