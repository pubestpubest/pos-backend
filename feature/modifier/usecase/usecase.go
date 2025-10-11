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
}

func NewModifierUsecase(modifierRepository domain.ModifierRepository) domain.ModifierUsecase {
	return &modifierUsecase{modifierRepository: modifierRepository}
}

func (u *modifierUsecase) GetAllModifiers() ([]*response.ModifierResponse, error) {
	modifiers, err := u.modifierRepository.GetAllModifiers()
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.GetAllModifiers]: Error getting modifiers")
	}

	modifierResponses := make([]*response.ModifierResponse, len(modifiers))
	for i, modifier := range modifiers {
		modifierResponses[i] = &response.ModifierResponse{
			ID:             modifier.ID,
			Name:           utils.DerefString(modifier.Name),
			PriceDeltaBaht: utils.DerefInt64(modifier.PriceDeltaBaht),
			Note:           utils.DerefString(modifier.Note),
		}
	}

	return modifierResponses, nil
}

func (u *modifierUsecase) GetModifierByID(id uuid.UUID) (*response.ModifierResponse, error) {
	modifier, err := u.modifierRepository.GetModifierByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.GetModifierByID]: Error getting modifier")
	}

	return &response.ModifierResponse{
		ID:             modifier.ID,
		Name:           utils.DerefString(modifier.Name),
		PriceDeltaBaht: utils.DerefInt64(modifier.PriceDeltaBaht),
		Note:           utils.DerefString(modifier.Note),
	}, nil
}

func (u *modifierUsecase) CreateModifier(req *request.ModifierRequest) (*response.ModifierResponse, error) {
	// Set default price delta to 0 if not provided
	priceDelta := int64(0)
	if req.PriceDeltaBaht != nil {
		priceDelta = *req.PriceDeltaBaht
	}

	modifier := &models.Modifier{
		Name:           &req.Name,
		PriceDeltaBaht: &priceDelta,
		Note:           req.Note,
	}

	if err := u.modifierRepository.CreateModifier(modifier); err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.CreateModifier]: Error creating modifier")
	}

	return &response.ModifierResponse{
		ID:             modifier.ID,
		Name:           req.Name,
		PriceDeltaBaht: priceDelta,
		Note:           utils.DerefString(req.Note),
	}, nil
}

func (u *modifierUsecase) UpdateModifier(id uuid.UUID, req *request.ModifierRequest) (*response.ModifierResponse, error) {
	// Get existing modifier
	modifier, err := u.modifierRepository.GetModifierByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.UpdateModifier]: Modifier not found")
	}

	// Update fields
	modifier.Name = &req.Name
	if req.PriceDeltaBaht != nil {
		modifier.PriceDeltaBaht = req.PriceDeltaBaht
	}
	modifier.Note = req.Note

	if err := u.modifierRepository.UpdateModifier(modifier); err != nil {
		return nil, errors.Wrap(err, "[ModifierUsecase.UpdateModifier]: Error updating modifier")
	}

	return &response.ModifierResponse{
		ID:             modifier.ID,
		Name:           utils.DerefString(modifier.Name),
		PriceDeltaBaht: utils.DerefInt64(modifier.PriceDeltaBaht),
		Note:           utils.DerefString(modifier.Note),
	}, nil
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
