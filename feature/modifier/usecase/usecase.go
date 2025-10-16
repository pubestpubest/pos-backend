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

type modifierUsecase struct {
	modifierRepository domain.ModifierRepository
	categoryRepository domain.CategoryRepository
}

func NewModifierUsecase(modifierRepository domain.ModifierRepository, categoryRepository domain.CategoryRepository) domain.ModifierUsecase {
	return &modifierUsecase{
		modifierRepository: modifierRepository,
		categoryRepository: categoryRepository,
	}
}

func (u *modifierUsecase) GetAllModifiers() ([]*response.ModifierResponse, error) {
	modifiers, err := u.modifierRepository.GetAllModifiers()
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.GetAllModifiers]: Error getting modifiers")
	}

	modifierResponses := make([]*response.ModifierResponse, len(modifiers))
	for i, modifier := range modifiers {
		modifierResponses[i] = u.mapModifierToResponse(modifier)
	}

	return modifierResponses, nil
}

func (u *modifierUsecase) GetModifierByID(id uuid.UUID) (*response.ModifierResponse, error) {
	modifier, err := u.modifierRepository.GetModifierByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.GetModifierByID]: Error getting modifier")
	}

	return u.mapModifierToResponse(modifier), nil
}

func (u *modifierUsecase) GetModifiersByCategoryID(categoryID uuid.UUID) ([]*response.ModifierResponse, error) {
	modifiers, err := u.modifierRepository.GetModifiersByCategoryID(categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.GetModifiersByCategoryID]: Error getting modifiers")
	}

	modifierResponses := make([]*response.ModifierResponse, len(modifiers))
	for i, modifier := range modifiers {
		modifierResponses[i] = u.mapModifierToResponse(modifier)
	}

	return modifierResponses, nil
}

func (u *modifierUsecase) CreateModifier(req *request.ModifierRequest) (*response.ModifierResponse, error) {
	// Validate category exists if provided
	if req.CategoryID != nil {
		_, err := u.categoryRepository.GetCategoryByID(*req.CategoryID)
		if err != nil {
			return nil, errors.Wrap(err, "[ModifierUsecase.CreateModifier]: Category not found")
		}
	}

	// Set default price delta to 0 if not provided
	priceDelta := int64(0)
	if req.PriceDeltaBaht != nil {
		priceDelta = *req.PriceDeltaBaht
	}

	modifier := &models.Modifier{
		CategoryID:     req.CategoryID,
		Name:           &req.Name,
		PriceDeltaBaht: &priceDelta,
		Note:           req.Note,
	}

	if err := u.modifierRepository.CreateModifier(modifier); err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.CreateModifier]: Error creating modifier")
	}

	// Fetch the created modifier with category relationship
	createdModifier, err := u.modifierRepository.GetModifierByID(modifier.ID)
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.CreateModifier]: Error fetching created modifier")
	}

	return u.mapModifierToResponse(createdModifier), nil
}

func (u *modifierUsecase) UpdateModifier(id uuid.UUID, req *request.ModifierRequest) (*response.ModifierResponse, error) {
	// Get existing modifier
	modifier, err := u.modifierRepository.GetModifierByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.UpdateModifier]: Modifier not found")
	}

	// Validate category exists if provided
	if req.CategoryID != nil {
		_, err := u.categoryRepository.GetCategoryByID(*req.CategoryID)
		if err != nil {
			return nil, errors.Wrap(err, "[ModifierUsecase.UpdateModifier]: Category not found")
		}
	}

	// Update fields
	modifier.CategoryID = req.CategoryID
	modifier.Name = &req.Name
	if req.PriceDeltaBaht != nil {
		modifier.PriceDeltaBaht = req.PriceDeltaBaht
	}
	modifier.Note = req.Note

	if err := u.modifierRepository.UpdateModifier(modifier); err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.UpdateModifier]: Error updating modifier")
	}

	// Fetch the updated modifier with category relationship
	updatedModifier, err := u.modifierRepository.GetModifierByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.UpdateModifier]: Error fetching updated modifier")
	}

	return u.mapModifierToResponse(updatedModifier), nil
}

func (u *modifierUsecase) DeleteModifier(id uuid.UUID) error {
	// Check if modifier exists
	_, err := u.modifierRepository.GetModifierByID(id)
	if err != nil {
		return errors.Wrap(err, "[ModifierUsecase.DeleteModifier]: Modifier not found")
	}

	if err := u.modifierRepository.DeleteModifier(id); err != nil {
		return errors.Wrap(err, "[ModifierUsecase.DeleteModifier]: Error deleting modifier")
	}

	return nil
}

// Helper function to map modifier model to response
func (u *modifierUsecase) mapModifierToResponse(modifier *models.Modifier) *response.ModifierResponse {
	resp := &response.ModifierResponse{
		ID:             modifier.ID,
		CategoryID:     modifier.CategoryID,
		Name:           utils.DerefString(modifier.Name),
		PriceDeltaBaht: utils.DerefInt64(modifier.PriceDeltaBaht),
		Note:           utils.DerefString(modifier.Note),
	}

	// Include category data if present
	if modifier.Category != nil {
		resp.Category = &response.CategoryResponse{
			ID:           modifier.Category.ID,
			Name:         utils.DerefString(modifier.Category.Name),
			DisplayOrder: utils.DerefInt(modifier.Category.DisplayOrder),
		}
	}

	return resp
}
