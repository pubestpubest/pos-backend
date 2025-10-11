package usecase

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/response"
	"github.com/pubestpubest/pos-backend/utils"
)

type roleUsecase struct {
	roleRepository domain.RoleRepository
}

func NewRoleUsecase(roleRepository domain.RoleRepository) domain.RoleUsecase {
	return &roleUsecase{roleRepository: roleRepository}
}

func (u *roleUsecase) GetAllRoles() ([]*response.RoleResponse, error) {
	roles, err := u.roleRepository.GetAllRoles()
	if err != nil {
		return nil, errors.Wrap(err, "[RoleUsecase.GetAllRoles]: Error getting roles")
	}

	roleResponses := make([]*response.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = &response.RoleResponse{
			ID:   role.ID,
			Name: role.Name,
		}
	}

	return roleResponses, nil
}

func (u *roleUsecase) GetRoleWithPermissions(id int) (*response.RoleResponse, error) {
	role, err := u.roleRepository.GetRoleWithPermissions(id)
	if err != nil {
		return nil, errors.Wrap(err, "[RoleUsecase.GetRoleWithPermissions]: Error getting role")
	}

	permissions := make([]response.PermissionResponse, len(role.Permissions))
	for i, perm := range role.Permissions {
		permissions[i] = response.PermissionResponse{
			ID:          perm.ID,
			Code:        perm.Code,
			Description: utils.DerefString(perm.Description),
		}
	}

	return &response.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}, nil
}
