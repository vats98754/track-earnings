-- Enable TimescaleDB and convert tables to hypertables
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Prices hypertable
SELECT create_hypertable('prices', 'time', if_not_exists => TRUE);

-- Macro hypertable
SELECT create_hypertable('macro_points', 'time', if_not_exists => TRUE);

-- Optional compression policies (tune as needed)
-- ALTER TABLE prices SET (timescaledb.compress, timescaledb.compress_segmentby = 'ticker');
-- SELECT add_compression_policy('prices', interval '7 days');
-- SELECT add_retention_policy('prices', interval '5 years');
