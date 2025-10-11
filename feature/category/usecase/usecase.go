package usecase

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
	"github.com/pubestpubest/pos-backend/utils"
)

type categoryUsecase struct {
	categoryRepository domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepository domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{categoryRepository: categoryRepository}
}

func (u *categoryUsecase) GetAllCategories() ([]*response.CategoryResponse, error) {
	categories, err := u.categoryRepository.GetAllCategories()
	if err != nil {
		return nil, errors.Wrap(err, "[CategoryUsecase.GetAllCategories]: Error getting categories")
	}

	categoryResponses := make([]*response.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = &response.CategoryResponse{
			ID:           category.ID,
			Name:         utils.DerefString(category.Name),
			DisplayOrder: utils.DerefInt(category.DisplayOrder),
		}
	}

	return categoryResponses, nil
}

func (u *categoryUsecase) GetCategoryByID(id uuid.UUID) (*response.CategoryResponse, error) {
	category, err := u.categoryRepository.GetCategoryByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[CategoryUsecase.GetCategoryByID]: Error getting category")
	}

	return &response.CategoryResponse{
		ID:           category.ID,
		Name:         utils.DerefString(category.Name),
		DisplayOrder: utils.DerefInt(category.DisplayOrder),
	}, nil
}

func (u *categoryUsecase) CreateCategory(req *request.CategoryRequest) (*response.CategoryResponse, error) {
	// Set default display order if not provided
	displayOrder := 0
	if req.DisplayOrder != nil {
		displayOrder = *req.DisplayOrder
	}

	category := &models.Category{
		Name:         &req.Name,
		DisplayOrder: &displayOrder,
	}

	if err := u.categoryRepository.CreateCategory(category); err != nil {
		return nil, errors.Wrap(err, "[CategoryUsecase.CreateCategory]: Error creating category")
	}

	return &response.CategoryResponse{
		ID:           category.ID,
		Name:         req.Name,
		DisplayOrder: displayOrder,
	}, nil
}

func (u *categoryUsecase) UpdateCategory(id uuid.UUID, req *request.CategoryRequest) (*response.CategoryResponse, error) {
	// Get existing category
	category, err := u.categoryRepository.GetCategoryByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[CategoryUsecase.UpdateCategory]: Category not found")
	}

	// Update fields
	category.Name = &req.Name
	if req.DisplayOrder != nil {
		category.DisplayOrder = req.DisplayOrder
	}

	if err := u.categoryRepository.UpdateCategory(category); err != nil {
		return nil, errors.Wrap(err, "[CategoryUsecase.UpdateCategory]: Error updating category")
	}

	return &response.CategoryResponse{
		ID:           category.ID,
		Name:         utils.DerefString(category.Name),
		DisplayOrder: utils.DerefInt(category.DisplayOrder),
	}, nil
}

func (u *categoryUsecase) DeleteCategory(id uuid.UUID) error {
	// Check if category exists
	_, err := u.categoryRepository.GetCategoryByID(id)
	if err != nil {
		return errors.Wrap(err, "[CategoryUsecase.DeleteCategory]: Category not found")
	}

	if err := u.categoryRepository.DeleteCategory(id); err != nil {
		return errors.Wrap(err, "[CategoryUsecase.DeleteCategory]: Error deleting category")
	}

	return nil
}
