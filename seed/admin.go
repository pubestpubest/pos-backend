package seed

import (
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedAdminUser(tx *gorm.DB) error {
	// Use env vars for credentials in real usage.
	username := SeedAdminUser.Username
	email := SeedAdminUser.Email
	hash := SeedAdminUser.PasswordHash // generate at deploy time

	var count int64
	if err := tx.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	u := models.User{
		Username:     username,
		PasswordHash: hash,
		FullName:     ptr(SeedAdminUser.FullName),
		Email:        &email,
		Status:       ptr(SeedAdminUser.Status),
	}
	return tx.Create(&u).Error
}

func seedAdminUserRoles(tx *gorm.DB) error {
	var u models.User
	if err := tx.Where("username = ?", SeedAdminUser.Username).First(&u).Error; err != nil {
		return err
	}
	var role models.Role
	if err := tx.Where("name = ?", "owner").First(&role).Error; err != nil {
		return err
	}
	ur := models.UserRole{UserID: u.ID, RoleID: role.ID}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}},
		DoNothing: true,
	}).Create(&ur).Error
}

func seedUsers(tx *gorm.DB) error {
	for _, su := range SeedUsers {
		var count int64
		if err := tx.Model(&models.User{}).Where("username = ?", su.Username).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue // User already exists
		}

		u := models.User{
			Username:     su.Username,
			PasswordHash: su.PasswordHash,
			FullName:     ptr(su.FullName),
			Email:        &su.Email,
			Phone:        &su.Phone,
			Status:       ptr(su.Status),
		}
		if err := tx.Create(&u).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedUserRoles(tx *gorm.DB) error {
	for _, su := range SeedUsers {
		var u models.User
		if err := tx.Where("username = ?", su.Username).First(&u).Error; err != nil {
			return err
		}
		var role models.Role
		if err := tx.Where("name = ?", su.RoleName).First(&role).Error; err != nil {
			return err
		}
		ur := models.UserRole{UserID: u.ID, RoleID: role.ID}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}},
			DoNothing: true,
		}).Create(&ur).Error; err != nil {
			return err
		}
	}
	return nil
}
