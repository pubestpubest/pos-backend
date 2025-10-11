package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// Category domain - manages menu categories
type CategoryUsecase interface {
	GetAllCategories() ([]*response.CategoryResponse, error)
	GetCategoryByID(id uuid.UUID) (*response.CategoryResponse, error)
	CreateCategory(req *request.CategoryRequest) (*response.CategoryResponse, error)
	UpdateCategory(id uuid.UUID, req *request.CategoryRequest) (*response.CategoryResponse, error)
	DeleteCategory(id uuid.UUID) error
}

type CategoryRepository interface {
	GetAllCategories() ([]*models.Category, error)
	GetCategoryByID(id uuid.UUID) (*models.Category, error)
	CreateCategory(category *models.Category) error
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uuid.UUID) error
}
