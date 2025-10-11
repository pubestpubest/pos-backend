package usecase

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/response"
	"github.com/pubestpubest/pos-backend/utils"
)

type tableUsecase struct {
	tableRepository domain.TableRepository
}

func NewTableUsecase(tableRepository domain.TableRepository) domain.TableUsecase {
	return &tableUsecase{tableRepository: tableRepository}
}

func (u *tableUsecase) GetAllTables() ([]*response.TableResponse, error) {
	tables, err := u.tableRepository.GetAllTables()
	if err != nil {
		return nil, errors.Wrap(err, "[TableUsecase.GetAllTables]: Error getting tables")
	}

	tableResponses := make([]*response.TableResponse, len(tables))
	for i, table := range tables {
		tableResponses[i] = &response.TableResponse{
			ID:     table.ID,
			Name:   utils.DerefString(table.Name),
			Seats:  utils.DerefInt(table.Seats),
			Status: utils.DerefString(table.Status),
			QRCode: utils.DerefString(table.QRSlug),
			AreaID: utils.DerefUUID(table.AreaID),
		}
	}

	return tableResponses, nil
}

func (u *tableUsecase) GetTableByID(id uuid.UUID) (*response.TableResponse, error) {
	table, err := u.tableRepository.GetTableByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[TableUsecase.GetTableByID]: Error getting table")
	}

	return &response.TableResponse{
		ID:     table.ID,
		Name:   utils.DerefString(table.Name),
		Seats:  utils.DerefInt(table.Seats),
		Status: utils.DerefString(table.Status),
		QRCode: utils.DerefString(table.QRSlug),
		AreaID: utils.DerefUUID(table.AreaID),
	}, nil
}

func (u *tableUsecase) UpdateTableStatus(id uuid.UUID, status string) error {
	// Get existing table
	table, err := u.tableRepository.GetTableByID(id)
	if err != nil {
		return errors.Wrap(err, "[TableUsecase.UpdateTableStatus]: Table not found")
	}

	// Update status
	table.Status = &status

	if err := u.tableRepository.UpdateTable(table); err != nil {
		return errors.Wrap(err, "[TableUsecase.UpdateTableStatus]: Error updating table status")
	}

	return nil
}
