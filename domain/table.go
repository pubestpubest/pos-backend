package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/response"
)

// Table domain - manages dining tables
type TableUsecase interface {
	GetAllTables() ([]*response.TableResponse, error)
	GetTableByID(id uuid.UUID) (*response.TableResponse, error)
	UpdateTableStatus(id uuid.UUID, status string) error
}

type TableRepository interface {
	GetAllTables() ([]*models.DiningTable, error)
	GetTableByID(id uuid.UUID) (*models.DiningTable, error)
	UpdateTable(table *models.DiningTable) error
}
