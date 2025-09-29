package seed

import (
	"time"

	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

func seedSampleOrders(tx *gorm.DB) error {
	var t models.DiningTable
	if err := tx.Where("name = ?", SeedDevSampleOrder.TableName).First(&t).Error; err != nil {
		return err
	}

	now := time.Now()
	o := models.Order{
		TableID:      &t.ID,
		Source:       ptr("staff"),
		Status:       ptr("open"),
		SubtotalBaht: ptrI64(12000),
		DiscountBaht: ptrI64(0),
		TotalBaht:    ptrI64(12000),
		CreatedAt:    now,
	}
	// idempotent by unique (table_id, created_at) is awkward; just skip if one open exists
	var exists int64
	if err := tx.Model(&models.Order{}).
		Where("table_id = ? AND status = 'open'", t.ID).Count(&exists).Error; err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}

	if err := tx.Create(&o).Error; err != nil {
		return err
	}

	// items
	var latte models.MenuItem
	if err := tx.Where("sku = ?", SeedDevSampleOrder.ItemSKU).First(&latte).Error; err != nil {
		return err
	}
	item := models.OrderItem{
		OrderID:       o.ID,
		MenuItemID:    latte.ID,
		Quantity:      SeedDevSampleOrder.Quantity,
		UnitPriceBaht: *latte.PriceBaht,
		LineTotalBaht: int64(SeedDevSampleOrder.Quantity) * *latte.PriceBaht,
	}
	if err := tx.Create(&item).Error; err != nil {
		return err
	}

	// recompute order totals in baht
	o.SubtotalBaht = ptrI64(item.LineTotalBaht)
	o.DiscountBaht = ptrI64(0)
	o.TotalBaht = ptrI64(item.LineTotalBaht)
	if err := tx.Model(&o).Updates(map[string]any{
		"subtotal_baht": *o.SubtotalBaht,
		"discount_baht": *o.DiscountBaht,
		"total_baht":    *o.TotalBaht,
	}).Error; err != nil {
		return err
	}

	// payment in baht
	p := models.Payment{OrderID: o.ID, Method: ptr("cash"), AmountBaht: *o.TotalBaht, Status: ptr("succeeded"), CreatedAt: now}
	return tx.Create(&p).Error
}
