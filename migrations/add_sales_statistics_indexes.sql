-- Migration: Add indexes for sales statistics queries
-- Purpose: Improve performance of menu item sales tracking queries
-- Created: 2024-10-16

-- Index on order_items.menu_item_id for faster joins
-- Used by: All sales statistics queries
CREATE INDEX IF NOT EXISTS idx_order_items_menu_item_id ON order_items(menu_item_id);

-- Composite index on orders (status, created_at) for filtering paid orders with date ranges
-- Used by: All sales statistics queries that filter by paid status and date
CREATE INDEX IF NOT EXISTS idx_orders_status_created_at ON orders(status, created_at);

-- Index on order_items.order_id for faster joins
-- Used by: Queries joining order_items to orders
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);

-- Index on order_item_modifiers.order_item_id for modifier popularity queries
-- Used by: GetMenuItemPopularModifiers query
CREATE INDEX IF NOT EXISTS idx_order_item_modifiers_order_item_id ON order_item_modifiers(order_item_id);

-- Index on order_item_modifiers.modifier_id for modifier joins
-- Used by: GetMenuItemPopularModifiers query
CREATE INDEX IF NOT EXISTS idx_order_item_modifiers_modifier_id ON order_item_modifiers(modifier_id);

-- Note: These indexes will speed up queries but will slightly slow down INSERT operations.
-- The trade-off is worth it for a reporting/analytics feature that reads more than writes.

