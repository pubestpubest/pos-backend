package repository

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAllCategories() ([]*models.Category, error) {
	var categories []*models.Category
	if err := r.db.Order("display_order ASC, name ASC").Find(&categories).Error; err != nil {
		return nil, errors.Wrap(err, "[CategoryRepository.GetAllCategories]: Error querying database")
	}
	return categories, nil
}

func (r *categoryRepository) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[CategoryRepository.GetCategoryByID]: Category not found")
		}
		return nil, errors.Wrap(err, "[CategoryRepository.GetCategoryByID]: Error querying database")
	}
	return &category, nil
}

func (r *categoryRepository) CreateCategory(category *models.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		return errors.Wrap(err, "[CategoryRepository.CreateCategory]: Error creating category")
	}
	return nil
}

func (r *categoryRepository) UpdateCategory(category *models.Category) error {
	if err := r.db.Save(category).Error; err != nil {
		return errors.Wrap(err, "[CategoryRepository.UpdateCategory]: Error updating category")
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Category{}).Error; err != nil {
		return errors.Wrap(err, "[CategoryRepository.DeleteCategory]: Error deleting category")
	}
	return nil
}
