package seed

import (
	"fmt"

	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedRoles(tx *gorm.DB) error {
	for _, roleName := range SeedRoles {
		r := models.Role{Name: roleName}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&r).Error; err != nil {
			return fmt.Errorf("seedRoles: %w", err)
		}
	}
	return nil
}
