package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// Modifier domain - manages menu item modifiers (add-ons, customizations)
type ModifierUsecase interface {
	GetAllModifiers() ([]*response.ModifierResponse, error)
	GetModifierByID(id uuid.UUID) (*response.ModifierResponse, error)
	GetModifiersByCategoryID(categoryID uuid.UUID) ([]*response.ModifierResponse, error)
	CreateModifier(req *request.ModifierRequest) (*response.ModifierResponse, error)
	UpdateModifier(id uuid.UUID, req *request.ModifierRequest) (*response.ModifierResponse, error)
	DeleteModifier(id uuid.UUID) error
}

type ModifierRepository interface {
	GetAllModifiers() ([]*models.Modifier, error)
	GetModifierByID(id uuid.UUID) (*models.Modifier, error)
	GetModifiersByCategoryID(categoryID uuid.UUID) ([]*models.Modifier, error)
	CreateModifier(modifier *models.Modifier) error
	UpdateModifier(modifier *models.Modifier) error
	DeleteModifier(id uuid.UUID) error
}
