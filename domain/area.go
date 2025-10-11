package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// Area domain - manages dining areas/sections
type AreaUsecase interface {
	GetAllAreas() ([]*response.AreaResponse, error)
	GetAreaByID(id uuid.UUID) (*response.AreaResponse, error)
	CreateArea(req *request.AreaRequest) (*response.AreaResponse, error)
	UpdateArea(id uuid.UUID, req *request.AreaRequest) (*response.AreaResponse, error)
	DeleteArea(id uuid.UUID) error
}

type AreaRepository interface {
	GetAllAreas() ([]*models.Area, error)
	GetAreaByID(id uuid.UUID) (*models.Area, error)
	CreateArea(area *models.Area) error
	UpdateArea(area *models.Area) error
	DeleteArea(id uuid.UUID) error
}
