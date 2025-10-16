package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type menuItemRepository struct {
	db *gorm.DB
}

func NewMenuItemRepository(db *gorm.DB) domain.MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r *menuItemRepository) GetAllMenuItems() ([]*models.MenuItem, error) {
	var menuItemsList []*models.MenuItem
	if err := r.db.Preload("Category").Order("name ASC").Find(&menuItemsList).Error; err != nil {
		return nil, errors.Wrap(err, "[MenuItemRepository.GetAllMenuItems]: Error getting menu items")
	}
	return menuItemsList, nil
}

func (r *menuItemRepository) GetMenuItemByID(id uuid.UUID) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Preload("Category").Where("id = ?", id).First(&menuItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[MenuItemRepository.GetMenuItemByID]: Menu item not found")
		}
		return nil, errors.Wrap(err, "[MenuItemRepository.GetMenuItemByID]: Error querying database")
	}
	return &menuItem, nil
}

func (r *menuItemRepository) CreateMenuItem(menuItem *models.MenuItem) error {
	if err := r.db.Create(menuItem).Error; err != nil {
		return errors.Wrap(err, "[MenuItemRepository.CreateMenuItem]: Error creating menu item")
	}
	return nil
}

func (r *menuItemRepository) UpdateMenuItem(menuItem *models.MenuItem) error {
	if err := r.db.Save(menuItem).Error; err != nil {
		return errors.Wrap(err, "[MenuItemRepository.UpdateMenuItem]: Error updating menu item")
	}
	return nil
}

func (r *menuItemRepository) DeleteMenuItem(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.MenuItem{}).Error; err != nil {
		return errors.Wrap(err, "[MenuItemRepository.DeleteMenuItem]: Error deleting menu item")
	}
	return nil
}

func (r *menuItemRepository) GetAllModifiers() ([]*models.Modifier, error) {
	var modifiers []*models.Modifier
	if err := r.db.Order("name ASC").Find(&modifiers).Error; err != nil {
		return nil, errors.Wrap(err, "[MenuItemRepository.GetAllModifiers]: Error getting modifiers")
	}
	return modifiers, nil
}

// GetMenuItemSalesStats returns sales statistics for all menu items
func (r *menuItemRepository) GetMenuItemSalesStats(startDate, endDate *time.Time, categoryID *uuid.UUID) ([]domain.MenuItemSalesData, error) {
	var results []domain.MenuItemSalesData

	query := r.db.Table("menu_items mi").
		Select(`
			mi.id as menu_item_id,
			COALESCE(mi.name, '') as menu_item_name,
			COALESCE(mi.sku, '') as sku,
			mi.category_id,
			COALESCE(c.name, '') as category_name,
			COALESCE(mi.price_baht, 0) as current_price_baht,
			COALESCE(SUM(oi.quantity), 0) as quantity_sold,
			COALESCE(SUM(oi.line_total_baht), 0) as total_revenue_baht,
			CASE 
				WHEN SUM(oi.quantity) > 0 THEN CAST(COALESCE(SUM(oi.line_total_baht) / SUM(oi.quantity), 0) AS BIGINT)
				ELSE 0
			END as average_price_baht,
			MIN(o.created_at) as first_sold_at,
			MAX(o.created_at) as last_sold_at,
			COUNT(DISTINCT o.id) as times_ordered
		`).
		Joins("LEFT JOIN categories c ON mi.category_id = c.id").
		Joins("LEFT JOIN order_items oi ON mi.id = oi.menu_item_id").
		Joins("LEFT JOIN orders o ON oi.order_id = o.id AND o.status = ?", "paid").
		Where("mi.deleted_at IS NULL")

	// Apply filters
	if startDate != nil {
		query = query.Where("o.created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("o.created_at <= ?", endDate)
	}
	if categoryID != nil {
		query = query.Where("mi.category_id = ?", categoryID)
	}

	query = query.Group("mi.id, mi.name, mi.sku, mi.category_id, c.name, mi.price_baht")

	if err := query.Scan(&results).Error; err != nil {
		return nil, errors.Wrap(err, "[MenuItemRepository.GetMenuItemSalesStats]: Error querying sales statistics")
	}

	return results, nil
}

// GetMenuItemSalesDetail returns detailed sales statistics for a specific menu item
func (r *menuItemRepository) GetMenuItemSalesDetail(menuItemID uuid.UUID, startDate, endDate *time.Time) (*domain.MenuItemSalesData, error) {
	var result domain.MenuItemSalesData

	query := r.db.Table("menu_items mi").
		Select(`
			mi.id as menu_item_id,
			COALESCE(mi.name, '') as menu_item_name,
			COALESCE(mi.sku, '') as sku,
			mi.category_id,
			COALESCE(c.name, '') as category_name,
			COALESCE(mi.price_baht, 0) as current_price_baht,
			COALESCE(SUM(oi.quantity), 0) as quantity_sold,
			COALESCE(SUM(oi.line_total_baht), 0) as total_revenue_baht,
			CASE 
				WHEN SUM(oi.quantity) > 0 THEN CAST(COALESCE(SUM(oi.line_total_baht) / SUM(oi.quantity), 0) AS BIGINT)
				ELSE 0
			END as average_price_baht,
			MIN(o.created_at) as first_sold_at,
			MAX(o.created_at) as last_sold_at,
			COUNT(DISTINCT o.id) as times_ordered
		`).
		Joins("LEFT JOIN categories c ON mi.category_id = c.id").
		Joins("LEFT JOIN order_items oi ON mi.id = oi.menu_item_id").
		Joins("LEFT JOIN orders o ON oi.order_id = o.id AND o.status = ?", "paid").
		Where("mi.id = ?", menuItemID).
		Where("mi.deleted_at IS NULL")

	if startDate != nil {
		query = query.Where("o.created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("o.created_at <= ?", endDate)
	}

	query = query.Group("mi.id, mi.name, mi.sku, mi.category_id, c.name, mi.price_baht")

	if err := query.Scan(&result).Error; err != nil {
		return nil, errors.Wrap(err, "[MenuItemRepository.GetMenuItemSalesDetail]: Error querying sales detail")
	}

	return &result, nil
}

// GetMenuItemTimeSeriesData returns sales data grouped by time period
func (r *menuItemRepository) GetMenuItemTimeSeriesData(menuItemID uuid.UUID, startDate, endDate *time.Time, groupBy string) ([]domain.TimeSeriesSalesData, error) {
	var results []domain.TimeSeriesSalesData

	// Determine date format based on groupBy
	dateFormat := "TO_CHAR(o.created_at, 'YYYY-MM-DD')"
	switch groupBy {
	case "week":
		dateFormat = "TO_CHAR(DATE_TRUNC('week', o.created_at), 'YYYY-MM-DD')"
	case "month":
		dateFormat = "TO_CHAR(DATE_TRUNC('month', o.created_at), 'YYYY-MM-DD')"
	}

	query := r.db.Table("orders o").
		Select(`
			`+dateFormat+` as period,
			COALESCE(SUM(oi.quantity), 0) as quantity_sold,
			COALESCE(SUM(oi.line_total_baht), 0) as revenue_baht,
			COUNT(DISTINCT o.id) as orders_count
		`).
		Joins("JOIN order_items oi ON o.id = oi.order_id").
		Where("oi.menu_item_id = ?", menuItemID).
		Where("o.status = ?", "paid")

	if startDate != nil {
		query = query.Where("o.created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("o.created_at <= ?", endDate)
	}

	query = query.Group("period").Order("period ASC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, errors.Wrap(err, "[MenuItemRepository.GetMenuItemTimeSeriesData]: Error querying time series data")
	}

	return results, nil
}

// GetMenuItemPopularModifiers returns the most popular modifiers for a menu item
func (r *menuItemRepository) GetMenuItemPopularModifiers(menuItemID uuid.UUID, startDate, endDate *time.Time, limit int) ([]domain.ModifierUsageData, error) {
	var results []domain.ModifierUsageData

	query := r.db.Table("order_item_modifiers oim").
		Select(`
			COALESCE(m.name, '') as modifier_name,
			COUNT(*) as times_added
		`).
		Joins("JOIN modifiers m ON oim.modifier_id = m.id").
		Joins("JOIN order_items oi ON oim.order_item_id = oi.id").
		Joins("JOIN orders o ON oi.order_id = o.id").
		Where("oi.menu_item_id = ?", menuItemID).
		Where("o.status = ?", "paid")

	if startDate != nil {
		query = query.Where("o.created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("o.created_at <= ?", endDate)
	}

	query = query.Group("m.name").
		Order("times_added DESC").
		Limit(limit)

	if err := query.Scan(&results).Error; err != nil {
		return nil, errors.Wrap(err, "[MenuItemRepository.GetMenuItemPopularModifiers]: Error querying popular modifiers")
	}

	return results, nil
}
