-- Add missing fields to financials
ALTER TABLE IF EXISTS income_statements ADD COLUMN IF NOT EXISTS ebitda DOUBLE PRECISION;
