package seed

import (
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedRolePermissions(tx *gorm.DB) error {
	for roleName, codes := range SeedRolePermissions {
		var role models.Role
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			return err
		}
		for _, c := range codes {
			var perm models.Permission
			if err := tx.Where("code = ?", c).First(&perm).Error; err != nil {
				return err
			}
			rp := models.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "role_id"}, {Name: "permission_id"}},
				DoNothing: true,
			}).Create(&rp).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
