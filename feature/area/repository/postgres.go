package repository

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type areaRepository struct {
	db *gorm.DB
}

func NewAreaRepository(db *gorm.DB) domain.AreaRepository {
	return &areaRepository{db: db}
}

func (r *areaRepository) GetAllAreas() ([]*models.Area, error) {
	var areas []*models.Area
	if err := r.db.Order("name ASC").Find(&areas).Error; err != nil {
		return nil, errors.Wrap(err, "[AreaRepository.GetAllAreas]: Error querying database")
	}
	return areas, nil
}

func (r *areaRepository) GetAreaByID(id uuid.UUID) (*models.Area, error) {
	var area models.Area
	if err := r.db.Where("id = ?", id).First(&area).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[AreaRepository.GetAreaByID]: Area not found")
		}
		return nil, errors.Wrap(err, "[AreaRepository.GetAreaByID]: Error querying database")
	}
	return &area, nil
}

func (r *areaRepository) CreateArea(area *models.Area) error {
	if err := r.db.Create(area).Error; err != nil {
		return errors.Wrap(err, "[AreaRepository.CreateArea]: Error creating area")
	}
	return nil
}

func (r *areaRepository) UpdateArea(area *models.Area) error {
	if err := r.db.Save(area).Error; err != nil {
		return errors.Wrap(err, "[AreaRepository.UpdateArea]: Error updating area")
	}
	return nil
}

func (r *areaRepository) DeleteArea(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Area{}).Error; err != nil {
		return errors.Wrap(err, "[AreaRepository.DeleteArea]: Error deleting area")
	}
	return nil
}
