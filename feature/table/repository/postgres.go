package repository

import (
	"github.com/google/uuid"
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
	if err := r.db.Preload("Area").Order("name ASC").Find(&tablesList).Error; err != nil {
		return nil, errors.Wrap(err, "[TableRepository.GetAllTables]: Error getting tables")
	}
	return tablesList, nil
}

func (r *tableRepository) GetTableByID(id uuid.UUID) (*models.DiningTable, error) {
	var table models.DiningTable
	if err := r.db.Preload("Area").Where("id = ?", id).First(&table).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[TableRepository.GetTableByID]: Table not found")
		}
		return nil, errors.Wrap(err, "[TableRepository.GetTableByID]: Error querying database")
	}
	return &table, nil
}

func (r *tableRepository) UpdateTable(table *models.DiningTable) error {
	if err := r.db.Save(table).Error; err != nil {
		return errors.Wrap(err, "[TableRepository.UpdateTable]: Error updating table")
	}
	return nil
}

func (r *tableRepository) GetTablesWithOpenOrders() ([]*models.DiningTable, error) {
	var tables []*models.DiningTable
	if err := r.db.Preload("Area").Preload("Orders", "status = ?", "open").Preload("Orders.Table").Preload("Orders.Items").Preload("Orders.Items.MenuItem").Preload("Orders.Items.Modifiers").Preload("Orders.Items.Modifiers.Modifier").Where("id IN (SELECT DISTINCT table_id FROM orders WHERE status = ? AND table_id IS NOT NULL)", "open").Order("name ASC").Find(&tables).Error; err != nil {
		return nil, errors.Wrap(err, "[TableRepository.GetTablesWithOpenOrders]: Error getting tables with open orders")
	}
	return tables, nil
}
