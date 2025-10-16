package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// MenuItem domain - manages menu items and their available modifiers
type MenuItemUsecase interface {
	GetAllMenuItems() ([]*response.MenuItemResponse, error)
	GetMenuItemByID(id uuid.UUID) (*response.MenuItemResponse, error)
	CreateMenuItem(req *request.MenuItemRequest) (*response.MenuItemResponse, error)
	UpdateMenuItem(id uuid.UUID, req *request.MenuItemRequest) (*response.MenuItemResponse, error)
	DeleteMenuItem(id uuid.UUID) error
	GetAvailableModifiers() ([]*response.ModifierResponse, error)

	// Sales statistics
	GetAllMenuItemsStatistics(startDate, endDate *time.Time, categoryID *uuid.UUID, sortBy, order string, limit int) (*response.AllMenuItemsStatisticsResponse, error)
	GetMenuItemStatistics(id uuid.UUID, startDate, endDate *time.Time, groupBy string) (*response.MenuItemDetailedStatisticsResponse, error)
	GetTopSellingItems(limit int, startDate, endDate *time.Time, metric string) (*response.TopSellingItemsResponse, error)
	GetLowSellingItems(startDate, endDate *time.Time, threshold int) (*response.LowSellingItemsResponse, error)
}

type MenuItemRepository interface {
	GetAllMenuItems() ([]*models.MenuItem, error)
	GetMenuItemByID(id uuid.UUID) (*models.MenuItem, error)
	CreateMenuItem(menuItem *models.MenuItem) error
	UpdateMenuItem(menuItem *models.MenuItem) error
	DeleteMenuItem(id uuid.UUID) error
	GetAllModifiers() ([]*models.Modifier, error)

	// Sales statistics queries
	GetMenuItemSalesStats(startDate, endDate *time.Time, categoryID *uuid.UUID) ([]MenuItemSalesData, error)
	GetMenuItemSalesDetail(menuItemID uuid.UUID, startDate, endDate *time.Time) (*MenuItemSalesData, error)
	GetMenuItemTimeSeriesData(menuItemID uuid.UUID, startDate, endDate *time.Time, groupBy string) ([]TimeSeriesSalesData, error)
	GetMenuItemPopularModifiers(menuItemID uuid.UUID, startDate, endDate *time.Time, limit int) ([]ModifierUsageData, error)
}

// Data transfer objects for repository layer
type MenuItemSalesData struct {
	MenuItemID       uuid.UUID
	MenuItemName     string
	SKU              string
	CategoryID       *uuid.UUID
	CategoryName     string
	CurrentPriceBaht int64
	QuantitySold     int64
	TotalRevenueBaht int64
	AveragePriceBaht int64
	FirstSoldAt      *time.Time
	LastSoldAt       *time.Time
	TimesOrdered     int64
}

type TimeSeriesSalesData struct {
	Period       string
	QuantitySold int64
	RevenueBaht  int64
	OrdersCount  int64
}

type ModifierUsageData struct {
	ModifierName string
	TimesAdded   int64
}
