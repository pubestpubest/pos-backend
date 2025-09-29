package repository

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type tableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) domain.TableRepository {
	return &tableRepository{db: db}
}

func (r *tableRepository) GetAllTables() ([]*models.DiningTable, error) {
	var tablesList []*models.DiningTable
	if err := r.db.Find(&tablesList).Error; err != nil {
		err = errors.Wrap(err, "[TableRepository.GetAllTables]: Error getting tables")
		return nil, err
	}

	return tablesList, nil
}
