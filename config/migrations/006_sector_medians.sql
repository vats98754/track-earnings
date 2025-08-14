-- Ensure sector_medians has required columns
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='sector_medians' AND column_name='sector'
    ) THEN
        ALTER TABLE IF EXISTS sector_medians ADD COLUMN sector TEXT;
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='sector_medians' AND column_name='metric'
    ) THEN
        ALTER TABLE IF EXISTS sector_medians ADD COLUMN metric TEXT;
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='sector_medians' AND column_name='value'
    ) THEN
        ALTER TABLE IF EXISTS sector_medians ADD COLUMN value DOUBLE PRECISION;
    END IF;
END$$;

-- If table does not exist at all (older setups), create it
CREATE TABLE IF NOT EXISTS sector_medians (
  as_of DATE NOT NULL,
  sector TEXT NOT NULL,
  metric TEXT NOT NULL,
  value DOUBLE PRECISION,
  PRIMARY KEY (as_of, sector, metric)
);
