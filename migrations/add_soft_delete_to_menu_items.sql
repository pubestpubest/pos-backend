-- Migration: Add soft delete support to menu_items table
-- Date: 2025-01-17
-- Description: Adds deleted_at column for soft delete functionality

-- Add deleted_at column for soft delete
ALTER TABLE menu_items 
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

-- Add index for better query performance when filtering out deleted records
CREATE INDEX IF NOT EXISTS idx_menu_items_deleted_at 
ON menu_items(deleted_at);

-- Update unique constraint on SKU to exclude soft-deleted records
-- This allows creating a new menu item with the same SKU as a deleted one
ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS menu_items_sku_key;
CREATE UNIQUE INDEX IF NOT EXISTS menu_items_sku_unique 
ON menu_items(sku) WHERE deleted_at IS NULL;

-- Verify the changes
SELECT 
    column_name, 
    data_type, 
    is_nullable
FROM 
    information_schema.columns
WHERE 
    table_name = 'menu_items' 
    AND column_name = 'deleted_at';

