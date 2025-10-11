package repository

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type modifierRepository struct {
	db *gorm.DB
}

func NewModifierRepository(db *gorm.DB) domain.ModifierRepository {
	return &modifierRepository{db: db}
}

func (r *modifierRepository) GetAllModifiers() ([]*models.Modifier, error) {
	var modifiers []*models.Modifier
	if err := r.db.Order("name ASC").Find(&modifiers).Error; err != nil {
		return nil, errors.Wrap(err, "[ModifierRepository.GetAllModifiers]: Error querying database")
	}
	return modifiers, nil
}

func (r *modifierRepository) GetModifierByID(id uuid.UUID) (*models.Modifier, error) {
	var modifier models.Modifier
	if err := r.db.Where("id = ?", id).First(&modifier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[ModifierRepository.GetModifierByID]: Modifier not found")
		}
		return nil, errors.Wrap(err, "[ModifierRepository.GetModifierByID]: Error querying database")
	}
	return &modifier, nil
}

func (r *modifierRepository) CreateModifier(modifier *models.Modifier) error {
	if err := r.db.Create(modifier).Error; err != nil {
		return errors.Wrap(err, "[ModifierRepository.CreateModifier]: Error creating modifier")
	}
	return nil
}

func (r *modifierRepository) UpdateModifier(modifier *models.Modifier) error {
	if err := r.db.Save(modifier).Error; err != nil {
		return errors.Wrap(err, "[ModifierRepository.UpdateModifier]: Error updating modifier")
	}
	return nil
}

func (r *modifierRepository) DeleteModifier(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Modifier{}).Error; err != nil {
		return errors.Wrap(err, "[ModifierRepository.DeleteModifier]: Error deleting modifier")
	}
	return nil
}
