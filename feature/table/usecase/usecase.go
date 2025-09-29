package usecase

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/response"
)

type tableUsecase struct {
	tableRepository domain.TableRepository
}

func NewTableUsecase(tableRepository domain.TableRepository) domain.TableUsecase {
	return &tableUsecase{tableRepository: tableRepository}
}

func (u *tableUsecase) GetAllTables() ([]*response.TableResponse, error) {
	table, err := u.tableRepository.GetAllTables()
	if err != nil {
		err = errors.Wrap(err, "[TableUsecase.GetAllTables]: Error getting table")
		return nil, err
	}

	tableResponses := make([]*response.TableResponse, len(table))
	for i, table := range table {
		tableResponses[i] = &response.TableResponse{
			ID:     table.ID,
			Name:   *table.Name,
			Seats:  *table.Seats,
			Status: *table.Status,
			QRCode: *table.QRSlug,
			AreaID: *table.AreaID,
		}
	}

	return tableResponses, nil
}
