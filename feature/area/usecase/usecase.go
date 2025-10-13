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

type areaUsecase struct {
	areaRepository domain.AreaRepository
}

func NewAreaUsecase(areaRepository domain.AreaRepository) domain.AreaUsecase {
	return &areaUsecase{areaRepository: areaRepository}
}

func (u *areaUsecase) GetAllAreas() ([]*response.AreaResponse, error) {
	areas, err := u.areaRepository.GetAllAreas()
	if err != nil {
		return nil, errors.Wrap(err, "[AreaUsecase.GetAllAreas]: Error getting areas")
	}

	areaResponses := make([]*response.AreaResponse, len(areas))
	for i, area := range areas {
		areaResponses[i] = &response.AreaResponse{
			ID:   area.ID,
			Name: utils.DerefString(area.Name),
		}
	}

	return areaResponses, nil
}

func (u *areaUsecase) GetAreaByID(id uuid.UUID) (*response.AreaResponse, error) {
	area, err := u.areaRepository.GetAreaByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[AreaUsecase.GetAreaByID]: Error getting area")
	}

	return &response.AreaResponse{
		ID:   area.ID,
		Name: utils.DerefString(area.Name),
	}, nil
}

func (u *areaUsecase) CreateArea(req *request.AreaRequest) (*response.AreaResponse, error) {
	area := &models.Area{
		Name: &req.Name,
	}

	if err := u.areaRepository.CreateArea(area); err != nil {
		return nil, errors.Wrap(err, "[AreaUsecase.CreateArea]: Error creating area")
	}

	return &response.AreaResponse{
		ID:   area.ID,
		Name: req.Name,
	}, nil
}

func (u *areaUsecase) UpdateArea(id uuid.UUID, req *request.AreaRequest) (*response.AreaResponse, error) {
	// Get existing area
	area, err := u.areaRepository.GetAreaByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[AreaUsecase.UpdateArea]: Area not found")
	}

	// Update fields
	area.Name = &req.Name

	if err := u.areaRepository.UpdateArea(area); err != nil {
		return nil, errors.Wrap(err, "[AreaUsecase.UpdateArea]: Error updating area")
	}

	return &response.AreaResponse{
		ID:   area.ID,
		Name: req.Name,
	}, nil
}

func (u *areaUsecase) DeleteArea(id uuid.UUID) error {
	// Check if area exists
	_, err := u.areaRepository.GetAreaByID(id)
	if err != nil {
		return errors.Wrap(err, "[AreaUsecase.DeleteArea]: Area not found")
	}

	if err := u.areaRepository.DeleteArea(id); err != nil {
		return errors.Wrap(err, "[AreaUsecase.DeleteArea]: Error deleting area")
	}

	return nil
}

// Helper function to build table response
func (u *areaUsecase) buildTableResponse(table *models.DiningTable) response.TableResponse {
	var areaResponse *response.AreaResponse
	if table.Area != nil {
		areaResponse = &response.AreaResponse{
			ID:   table.Area.ID,
			Name: utils.DerefString(table.Area.Name),
		}
	}

	return response.TableResponse{
		ID:     table.ID,
		Name:   utils.DerefString(table.Name),
		Seats:  utils.DerefInt(table.Seats),
		Status: utils.DerefString(table.Status),
		QRCode: utils.DerefString(table.QRSlug),
		Area:   areaResponse,
	}
}

func (u *areaUsecase) GetAreasWithTables() ([]*response.AreaWithTablesResponse, error) {
	areas, err := u.areaRepository.GetAreasWithTables()
	if err != nil {
		return nil, errors.Wrap(err, "[AreaUsecase.GetAreasWithTables]: Error getting areas with tables")
	}

	areaResponses := make([]*response.AreaWithTablesResponse, len(areas))
	for i, area := range areas {
		tables := make([]response.TableResponse, len(area.Tables))
		for j, table := range area.Tables {
			tables[j] = u.buildTableResponse(&table)
		}

		areaResponses[i] = &response.AreaWithTablesResponse{
			ID:     area.ID,
			Name:   utils.DerefString(area.Name),
			Tables: tables,
		}
	}

	return areaResponses, nil
}
