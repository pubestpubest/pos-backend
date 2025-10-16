package usecase

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
	"github.com/pubestpubest/pos-backend/utils"
)

type menuItemUsecase struct {
	menuItemRepository domain.MenuItemRepository
}

func NewMenuItemUsecase(menuItemRepository domain.MenuItemRepository) domain.MenuItemUsecase {
	return &menuItemUsecase{menuItemRepository: menuItemRepository}
}

func (u *menuItemUsecase) GetAllMenuItems() ([]*response.MenuItemResponse, error) {
	menuItems, err := u.menuItemRepository.GetAllMenuItems()
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetAllMenuItems]: Error getting menu items")
	}

	menuItemResponses := make([]*response.MenuItemResponse, len(menuItems))
	for i, menuItem := range menuItems {
		menuItemResponses[i] = &response.MenuItemResponse{
			ID:        menuItem.ID,
			Name:      utils.DerefString(menuItem.Name),
			SKU:       utils.DerefString(menuItem.SKU),
			PriceBaht: utils.DerefInt64(menuItem.PriceBaht),
			Active:    utils.DerefBool(menuItem.Active),
			ImageURL:  utils.DerefString(menuItem.ImageURL),
			Category:  response.CategoryResponse{ID: utils.DerefUUID(menuItem.CategoryID), Name: utils.DerefString(menuItem.Category.Name), DisplayOrder: utils.DerefInt(menuItem.Category.DisplayOrder)},
		}
	}

	return menuItemResponses, nil
}

func (u *menuItemUsecase) GetMenuItemByID(id uuid.UUID) (*response.MenuItemResponse, error) {
	menuItem, err := u.menuItemRepository.GetMenuItemByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetMenuItemByID]: Error getting menu item")
	}

	return &response.MenuItemResponse{
		ID:        menuItem.ID,
		Name:      utils.DerefString(menuItem.Name),
		SKU:       utils.DerefString(menuItem.SKU),
		PriceBaht: utils.DerefInt64(menuItem.PriceBaht),
		Active:    utils.DerefBool(menuItem.Active),
		ImageURL:  utils.DerefString(menuItem.ImageURL),
		Category:  response.CategoryResponse{ID: utils.DerefUUID(menuItem.CategoryID), Name: utils.DerefString(menuItem.Category.Name), DisplayOrder: utils.DerefInt(menuItem.Category.DisplayOrder)},
	}, nil
}

func (u *menuItemUsecase) CreateMenuItem(req *request.MenuItemRequest) (*response.MenuItemResponse, error) {
	active := true
	if req.Active != nil {
		active = *req.Active
	}

	menuItem := &models.MenuItem{
		CategoryID: req.CategoryID,
		Name:       &req.Name,
		SKU:        &req.SKU,
		PriceBaht:  &req.PriceBaht,
		Active:     &active,
		ImageURL:   req.ImageURL,
	}

	if err := u.menuItemRepository.CreateMenuItem(menuItem); err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.CreateMenuItem]: Error creating menu item")
	}

	return &response.MenuItemResponse{
		ID:        menuItem.ID,
		Name:      req.Name,
		SKU:       req.SKU,
		PriceBaht: req.PriceBaht,
		Active:    active,
		ImageURL:  utils.DerefString(req.ImageURL),
		Category:  response.CategoryResponse{ID: utils.DerefUUID(req.CategoryID), Name: "", DisplayOrder: 0},
	}, nil
}

func (u *menuItemUsecase) UpdateMenuItem(id uuid.UUID, req *request.MenuItemRequest) (*response.MenuItemResponse, error) {
	// Get existing menu item
	menuItem, err := u.menuItemRepository.GetMenuItemByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.UpdateMenuItem]: Menu item not found")
	}

	// Update fields
	menuItem.CategoryID = req.CategoryID
	menuItem.Name = &req.Name
	menuItem.SKU = &req.SKU
	menuItem.PriceBaht = &req.PriceBaht
	menuItem.ImageURL = req.ImageURL
	if req.Active != nil {
		menuItem.Active = req.Active
	}

	if err := u.menuItemRepository.UpdateMenuItem(menuItem); err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.UpdateMenuItem]: Error updating menu item")
	}

	return &response.MenuItemResponse{
		ID:        menuItem.ID,
		Name:      utils.DerefString(menuItem.Name),
		SKU:       utils.DerefString(menuItem.SKU),
		PriceBaht: utils.DerefInt64(menuItem.PriceBaht),
		Active:    utils.DerefBool(menuItem.Active),
		ImageURL:  utils.DerefString(menuItem.ImageURL),
		Category:  response.CategoryResponse{ID: utils.DerefUUID(menuItem.CategoryID), Name: utils.DerefString(menuItem.Category.Name), DisplayOrder: utils.DerefInt(menuItem.Category.DisplayOrder)},
	}, nil
}

func (u *menuItemUsecase) DeleteMenuItem(id uuid.UUID) error {
	// Check if menu item exists
	_, err := u.menuItemRepository.GetMenuItemByID(id)
	if err != nil {
		return errors.Wrap(err, "[MenuItemUsecase.DeleteMenuItem]: Menu item not found")
	}

	if err := u.menuItemRepository.DeleteMenuItem(id); err != nil {
		return errors.Wrap(err, "[MenuItemUsecase.DeleteMenuItem]: Error deleting menu item")
	}

	return nil
}

func (u *menuItemUsecase) GetAvailableModifiers() ([]*response.ModifierResponse, error) {
	modifiers, err := u.menuItemRepository.GetAllModifiers()
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetAvailableModifiers]: Error getting modifiers")
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

// GetAllMenuItemsStatistics returns sales statistics for all menu items
func (u *menuItemUsecase) GetAllMenuItemsStatistics(startDate, endDate *time.Time, categoryID *uuid.UUID, sortBy, order string, limit int) (*response.AllMenuItemsStatisticsResponse, error) {
	// Get sales data from repository
	salesData, err := u.menuItemRepository.GetMenuItemSalesStats(startDate, endDate, categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetAllMenuItemsStatistics]: Error getting sales statistics")
	}

	// Convert to response format
	statistics := make([]response.MenuItemStatisticsResponse, len(salesData))
	totalQuantitySold := int64(0)
	totalRevenueBaht := int64(0)

	for i, data := range salesData {
		statistics[i] = response.MenuItemStatisticsResponse{
			MenuItemID:       data.MenuItemID,
			MenuItemName:     data.MenuItemName,
			SKU:              data.SKU,
			CategoryID:       data.CategoryID,
			CategoryName:     data.CategoryName,
			CurrentPriceBaht: data.CurrentPriceBaht,
			QuantitySold:     data.QuantitySold,
			TotalRevenueBaht: data.TotalRevenueBaht,
			AveragePriceBaht: data.AveragePriceBaht,
			FirstSoldAt:      data.FirstSoldAt,
			LastSoldAt:       data.LastSoldAt,
			TimesOrdered:     data.TimesOrdered,
		}

		totalQuantitySold += data.QuantitySold
		totalRevenueBaht += data.TotalRevenueBaht
	}

	// Sort the results
	switch sortBy {
	case "revenue":
		sort.Slice(statistics, func(i, j int) bool {
			if order == "asc" {
				return statistics[i].TotalRevenueBaht < statistics[j].TotalRevenueBaht
			}
			return statistics[i].TotalRevenueBaht > statistics[j].TotalRevenueBaht
		})
	case "name":
		sort.Slice(statistics, func(i, j int) bool {
			if order == "asc" {
				return statistics[i].MenuItemName < statistics[j].MenuItemName
			}
			return statistics[i].MenuItemName > statistics[j].MenuItemName
		})
	default: // quantity_sold
		sort.Slice(statistics, func(i, j int) bool {
			if order == "asc" {
				return statistics[i].QuantitySold < statistics[j].QuantitySold
			}
			return statistics[i].QuantitySold > statistics[j].QuantitySold
		})
	}

	// Apply limit if specified
	if limit > 0 && limit < len(statistics) {
		statistics = statistics[:limit]
	}

	// Build response
	return &response.AllMenuItemsStatisticsResponse{
		Statistics: statistics,
		Summary: response.SalesSummary{
			TotalItemsTracked: len(salesData),
			TotalQuantitySold: totalQuantitySold,
			TotalRevenueBaht:  totalRevenueBaht,
			DateRange: response.DateRange{
				Start: startDate,
				End:   endDate,
			},
		},
	}, nil
}

// GetMenuItemStatistics returns detailed statistics for a specific menu item
func (u *menuItemUsecase) GetMenuItemStatistics(id uuid.UUID, startDate, endDate *time.Time, groupBy string) (*response.MenuItemDetailedStatisticsResponse, error) {
	// Get sales detail
	salesDetail, err := u.menuItemRepository.GetMenuItemSalesDetail(id, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetMenuItemStatistics]: Error getting sales detail")
	}

	// Get time series data if groupBy is specified
	var timeSeries []response.TimeSeriesData
	if groupBy != "" {
		tsData, err := u.menuItemRepository.GetMenuItemTimeSeriesData(id, startDate, endDate, groupBy)
		if err != nil {
			return nil, errors.Wrap(err, "[MenuItemUsecase.GetMenuItemStatistics]: Error getting time series data")
		}

		timeSeries = make([]response.TimeSeriesData, len(tsData))
		for i, ts := range tsData {
			timeSeries[i] = response.TimeSeriesData{
				Period:       ts.Period,
				QuantitySold: ts.QuantitySold,
				RevenueBaht:  ts.RevenueBaht,
				OrdersCount:  ts.OrdersCount,
			}
		}
	}

	// Get popular modifiers
	modifierData, err := u.menuItemRepository.GetMenuItemPopularModifiers(id, startDate, endDate, 5)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetMenuItemStatistics]: Error getting popular modifiers")
	}

	popularModifiers := make([]response.PopularModifier, len(modifierData))
	for i, mod := range modifierData {
		popularModifiers[i] = response.PopularModifier{
			ModifierName: mod.ModifierName,
			TimesAdded:   mod.TimesAdded,
		}
	}

	// Calculate average quantity per order
	avgQuantityPerOrder := float64(0)
	if salesDetail.TimesOrdered > 0 {
		avgQuantityPerOrder = float64(salesDetail.QuantitySold) / float64(salesDetail.TimesOrdered)
	}

	// Build response
	return &response.MenuItemDetailedStatisticsResponse{
		MenuItemID:       salesDetail.MenuItemID,
		MenuItemName:     salesDetail.MenuItemName,
		SKU:              salesDetail.SKU,
		CategoryName:     salesDetail.CategoryName,
		CurrentPriceBaht: salesDetail.CurrentPriceBaht,
		Statistics: response.DetailedStats{
			TotalQuantitySold:       salesDetail.QuantitySold,
			TotalRevenueBaht:        salesDetail.TotalRevenueBaht,
			AverageQuantityPerOrder: avgQuantityPerOrder,
			TimesOrdered:            salesDetail.TimesOrdered,
			FirstSoldAt:             salesDetail.FirstSoldAt,
			LastSoldAt:              salesDetail.LastSoldAt,
		},
		TimeSeries:       timeSeries,
		PopularModifiers: popularModifiers,
	}, nil
}

// GetTopSellingItems returns the top selling menu items
func (u *menuItemUsecase) GetTopSellingItems(limit int, startDate, endDate *time.Time, metric string) (*response.TopSellingItemsResponse, error) {
	// Get all sales data
	salesData, err := u.menuItemRepository.GetMenuItemSalesStats(startDate, endDate, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetTopSellingItems]: Error getting sales statistics")
	}

	// Calculate total for percentage
	totalMetricValue := int64(0)
	for _, data := range salesData {
		if metric == "revenue" {
			totalMetricValue += data.TotalRevenueBaht
		} else {
			totalMetricValue += data.QuantitySold
		}
	}

	// Sort by metric
	sort.Slice(salesData, func(i, j int) bool {
		if metric == "revenue" {
			return salesData[i].TotalRevenueBaht > salesData[j].TotalRevenueBaht
		}
		return salesData[i].QuantitySold > salesData[j].QuantitySold
	})

	// Apply limit
	if limit > 0 && limit < len(salesData) {
		salesData = salesData[:limit]
	}

	// Build response
	topItems := make([]response.TopSellingItemResponse, len(salesData))
	for i, data := range salesData {
		percentage := float64(0)
		if totalMetricValue > 0 {
			if metric == "revenue" {
				percentage = float64(data.TotalRevenueBaht) / float64(totalMetricValue) * 100
			} else {
				percentage = float64(data.QuantitySold) / float64(totalMetricValue) * 100
			}
		}

		topItems[i] = response.TopSellingItemResponse{
			Rank:              i + 1,
			MenuItemID:        data.MenuItemID,
			MenuItemName:      data.MenuItemName,
			CategoryName:      data.CategoryName,
			QuantitySold:      data.QuantitySold,
			RevenueBaht:       data.TotalRevenueBaht,
			PercentageOfTotal: percentage,
		}
	}

	return &response.TopSellingItemsResponse{
		TopItems: topItems,
		Period: response.DateRange{
			Start: startDate,
			End:   endDate,
		},
	}, nil
}

// GetLowSellingItems returns items with low or no sales
func (u *menuItemUsecase) GetLowSellingItems(startDate, endDate *time.Time, threshold int) (*response.LowSellingItemsResponse, error) {
	// Get all sales data
	salesData, err := u.menuItemRepository.GetMenuItemSalesStats(startDate, endDate, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetLowSellingItems]: Error getting sales statistics")
	}

	lowSellingItems := []response.LowSellingItemResponse{}
	noSalesItems := []response.NoSalesItemResponse{}

	for _, data := range salesData {
		if data.QuantitySold == 0 {
			// No sales at all
			noSalesItems = append(noSalesItems, response.NoSalesItemResponse{
				MenuItemID:   data.MenuItemID,
				MenuItemName: data.MenuItemName,
				CategoryName: data.CategoryName,
				IsActive:     true, // Assuming true, would need to query if needed
				CreatedAt:    time.Now(),
			})
		} else if data.QuantitySold < int64(threshold) {
			// Low sales
			daysSinceLastSale := 0
			if data.LastSoldAt != nil {
				daysSinceLastSale = int(time.Since(*data.LastSoldAt).Hours() / 24)
			}

			lowSellingItems = append(lowSellingItems, response.LowSellingItemResponse{
				MenuItemID:        data.MenuItemID,
				MenuItemName:      data.MenuItemName,
				CategoryName:      data.CategoryName,
				QuantitySold:      data.QuantitySold,
				RevenueBaht:       data.TotalRevenueBaht,
				DaysSinceLastSale: &daysSinceLastSale,
				IsActive:          true,
			})
		}
	}

	return &response.LowSellingItemsResponse{
		LowSellingItems: lowSellingItems,
		NoSalesItems:    noSalesItems,
		Period: response.DateRange{
			Start: startDate,
			End:   endDate,
		},
	}, nil
}
