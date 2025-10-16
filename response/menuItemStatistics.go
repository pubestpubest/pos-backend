package response

import (
	"time"

	"github.com/google/uuid"
)

// MenuItemStatisticsResponse - Individual menu item sales stats
type MenuItemStatisticsResponse struct {
	MenuItemID       uuid.UUID  `json:"menu_item_id"`
	MenuItemName     string     `json:"menu_item_name"`
	SKU              string     `json:"sku"`
	CategoryID       *uuid.UUID `json:"category_id"`
	CategoryName     string     `json:"category_name,omitempty"`
	CurrentPriceBaht int64      `json:"current_price_baht"`
	QuantitySold     int64      `json:"quantity_sold"`
	TotalRevenueBaht int64      `json:"total_revenue_baht"`
	AveragePriceBaht int64      `json:"average_price_baht"`
	FirstSoldAt      *time.Time `json:"first_sold_at"`
	LastSoldAt       *time.Time `json:"last_sold_at"`
	TimesOrdered     int64      `json:"times_ordered"` // Number of distinct orders
}

// AllMenuItemsStatisticsResponse - Aggregated stats for all items
type AllMenuItemsStatisticsResponse struct {
	Statistics []MenuItemStatisticsResponse `json:"statistics"`
	Summary    SalesSummary                 `json:"summary"`
}

// SalesSummary - Overall summary
type SalesSummary struct {
	TotalItemsTracked int       `json:"total_items_tracked"`
	TotalQuantitySold int64     `json:"total_quantity_sold"`
	TotalRevenueBaht  int64     `json:"total_revenue_baht"`
	DateRange         DateRange `json:"date_range"`
}

// DateRange - Time period for statistics
type DateRange struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// MenuItemDetailedStatisticsResponse - Detailed stats with time series
type MenuItemDetailedStatisticsResponse struct {
	MenuItemID       uuid.UUID         `json:"menu_item_id"`
	MenuItemName     string            `json:"menu_item_name"`
	SKU              string            `json:"sku"`
	CategoryName     string            `json:"category_name,omitempty"`
	CurrentPriceBaht int64             `json:"current_price_baht"`
	Statistics       DetailedStats     `json:"statistics"`
	TimeSeries       []TimeSeriesData  `json:"time_series,omitempty"`
	PopularModifiers []PopularModifier `json:"popular_modifiers,omitempty"`
}

// DetailedStats - Core statistics
type DetailedStats struct {
	TotalQuantitySold       int64      `json:"total_quantity_sold"`
	TotalRevenueBaht        int64      `json:"total_revenue_baht"`
	AverageQuantityPerOrder float64    `json:"average_quantity_per_order"`
	TimesOrdered            int64      `json:"times_ordered"`
	FirstSoldAt             *time.Time `json:"first_sold_at"`
	LastSoldAt              *time.Time `json:"last_sold_at"`
}

// TimeSeriesData - Sales over time
type TimeSeriesData struct {
	Period       string `json:"period"` // Date string
	QuantitySold int64  `json:"quantity_sold"`
	RevenueBaht  int64  `json:"revenue_baht"`
	OrdersCount  int64  `json:"orders_count"`
}

// PopularModifier - Most common modifiers
type PopularModifier struct {
	ModifierName string `json:"modifier_name"`
	TimesAdded   int64  `json:"times_added"`
}

// TopSellingItemResponse - Ranked item for top sellers report
type TopSellingItemResponse struct {
	Rank              int       `json:"rank"`
	MenuItemID        uuid.UUID `json:"menu_item_id"`
	MenuItemName      string    `json:"menu_item_name"`
	CategoryName      string    `json:"category_name,omitempty"`
	QuantitySold      int64     `json:"quantity_sold"`
	RevenueBaht       int64     `json:"revenue_baht"`
	PercentageOfTotal float64   `json:"percentage_of_total"`
}

// TopSellingItemsResponse - Response wrapper for top sellers
type TopSellingItemsResponse struct {
	TopItems []TopSellingItemResponse `json:"top_items"`
	Period   DateRange                `json:"period"`
}

// LowSellingItemResponse - Items with low sales
type LowSellingItemResponse struct {
	MenuItemID        uuid.UUID `json:"menu_item_id"`
	MenuItemName      string    `json:"menu_item_name"`
	CategoryName      string    `json:"category_name,omitempty"`
	QuantitySold      int64     `json:"quantity_sold"`
	RevenueBaht       int64     `json:"revenue_baht"`
	DaysSinceLastSale *int      `json:"days_since_last_sale,omitempty"`
	IsActive          bool      `json:"is_active"`
}

// NoSalesItemResponse - Items never sold
type NoSalesItemResponse struct {
	MenuItemID   uuid.UUID `json:"menu_item_id"`
	MenuItemName string    `json:"menu_item_name"`
	CategoryName string    `json:"category_name,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// LowSellingItemsResponse - Response wrapper for low/no sales items
type LowSellingItemsResponse struct {
	LowSellingItems []LowSellingItemResponse `json:"low_selling_items"`
	NoSalesItems    []NoSalesItemResponse    `json:"no_sales_items"`
	Period          DateRange                `json:"period"`
}
