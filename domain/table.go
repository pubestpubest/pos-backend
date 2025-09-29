package domain

import (
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/response"
)

type TableUsecase interface {
	GetAllTables() ([]*response.TableResponse, error)
}

type TableRepository interface {
	GetAllTables() ([]*models.DiningTable, error)
}
