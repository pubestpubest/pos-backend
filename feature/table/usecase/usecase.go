package usecase

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
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
			Area:   u.buildAreaResponse(table.Area),
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
		Area:   u.buildAreaResponse(table.Area),
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

// Helper function to build order response
func (u *tableUsecase) buildOrderResponse(order *models.Order) response.OrderResponse {
	var items []response.OrderItemResponse
	for _, item := range order.Items {
		modifiers := make([]response.OrderItemModifierResponse, len(item.Modifiers))
		for j, mod := range item.Modifiers {
			var modifierName string
			if mod.Modifier != nil {
				modifierName = utils.DerefString(mod.Modifier.Name)
			}
			modifiers[j] = response.OrderItemModifierResponse{
				ModifierID:     mod.ModifierID,
				ModifierName:   modifierName,
				PriceDeltaBaht: utils.DerefInt64(mod.PriceDeltaBaht),
			}
		}

		var menuItemName string
		if item.MenuItem != nil {
			menuItemName = utils.DerefString(item.MenuItem.Name)
		}

		items = append(items, response.OrderItemResponse{
			ID:            item.ID,
			MenuItemID:    item.MenuItemID,
			MenuItemName:  menuItemName,
			Quantity:      item.Quantity,
			UnitPriceBaht: item.UnitPriceBaht,
			LineTotalBaht: item.LineTotalBaht,
			Note:          utils.DerefString(item.Note),
			Modifiers:     modifiers,
		})
	}

	var tableName string
	if order.Table != nil {
		tableName = utils.DerefString(order.Table.Name)
	}

	return response.OrderResponse{
		ID:           order.ID,
		TableID:      utils.DerefUUID(order.TableID),
		TableName:    tableName,
		OpenedBy:     order.OpenedBy,
		Source:       utils.DerefString(order.Source),
		Status:       utils.DerefString(order.Status),
		SubtotalBaht: utils.DerefInt64(order.SubtotalBaht),
		DiscountBaht: utils.DerefInt64(order.DiscountBaht),
		TotalBaht:    utils.DerefInt64(order.TotalBaht),
		Note:         utils.DerefString(order.Note),
		CreatedAt:    order.CreatedAt,
		ClosedAt:     order.ClosedAt,
		Items:        items,
	}
}

func (u *tableUsecase) GetTablesWithOpenOrders() ([]*response.TableWithOrdersResponse, error) {
	tables, err := u.tableRepository.GetTablesWithOpenOrders()
	if err != nil {
		return nil, errors.Wrap(err, "[TableUsecase.GetTablesWithOpenOrders]: Error getting tables with open orders")
	}

	tableResponses := make([]*response.TableWithOrdersResponse, len(tables))
	for i, table := range tables {
		orders := make([]response.OrderResponse, len(table.Orders))
		for j, order := range table.Orders {
			orders[j] = u.buildOrderResponse(&order)
		}

		tableResponses[i] = &response.TableWithOrdersResponse{
			ID:     table.ID,
			Name:   utils.DerefString(table.Name),
			Status: utils.DerefString(table.Status),
			QRCode: utils.DerefString(table.QRSlug),
			Orders: orders,
		}
	}

	return tableResponses, nil
}

// Helper function to build area response
func (u *tableUsecase) buildAreaResponse(area *models.Area) *response.AreaResponse {
	if area == nil {
		return nil
	}
	return &response.AreaResponse{
		ID:   area.ID,
		Name: utils.DerefString(area.Name),
	}
}
