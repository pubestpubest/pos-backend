package seed

import (
	"fmt"

	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedPermissions(tx *gorm.DB) error {
	for _, sp := range SeedPermissions {
		p := models.Permission{Code: sp.Code, Description: ptr(sp.Description)}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoNothing: true,
		}).Create(&p).Error; err != nil {
			return fmt.Errorf("seedPermissions: %w", err)
		}
	}
	return nil
}
