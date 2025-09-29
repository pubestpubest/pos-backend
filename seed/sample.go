package seed

import (
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedAreas(tx *gorm.DB) error {
	for _, n := range SeedAreas {
		a := models.Area{Name: &n}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&a).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedTables(tx *gorm.DB) error {
	for _, r := range SeedTables {
		var area models.Area
		if err := tx.Where("name = ?", r.AreaName).First(&area).Error; err != nil {
			return err
		}
		t := models.DiningTable{
			AreaID: &area.ID,
			Name:   &r.Name,
			Seats:  &r.Seats,
			Status: ptr("free"),
			QRSlug: &r.Slug,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "qr_slug"}},
			DoNothing: true,
		}).Create(&t).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedCategories(tx *gorm.DB) error {
	for _, c := range SeedCategories {
		rec := models.Category{Name: ptr(c.Name), DisplayOrder: ptrInt(c.DisplayOrder)}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&rec).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedMenuItems(tx *gorm.DB) error {
	for _, it := range SeedMenuItems {
		var cat models.Category
		if err := tx.Where("name = ?", it.CategoryName).First(&cat).Error; err != nil {
			return err
		}
		i := models.MenuItem{
			CategoryID: &cat.ID,
			Name:       &it.Name,
			SKU:        &it.SKU,
			PriceBaht:  ptrI64(int64(it.PriceBaht)),
			Active:     ptrBool(true),
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "sku"}},
			DoNothing: true,
		}).Create(&i).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedModifiers(tx *gorm.DB) error {
	for _, m := range SeedModifiers {
		rec := models.Modifier{Name: &m.Name, PriceDeltaBaht: ptrI64(int64(m.DeltaBaht))}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&rec).Error; err != nil {
			return err
		}
	}
	return nil
}
