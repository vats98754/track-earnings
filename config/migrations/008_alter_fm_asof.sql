-- Align fundamentals_metrics.asof to TIMESTAMPTZ if needed
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='fundamentals_metrics' AND column_name='asof' AND data_type='date'
    ) THEN
        ALTER TABLE fundamentals_metrics ALTER COLUMN asof TYPE TIMESTAMPTZ USING (asof::timestamptz);
    END IF;
END$$;
