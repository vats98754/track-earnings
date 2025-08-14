-- Earnings quality table
CREATE TABLE IF NOT EXISTS earnings_quality (
  ticker TEXT NOT NULL,
  as_of DATE NOT NULL,
  reported_eps DOUBLE PRECISION,
  adjusted_eps DOUBLE PRECISION,
  delta DOUBLE PRECISION,
  notes TEXT,
  provider TEXT NOT NULL,
  PRIMARY KEY (ticker, as_of)
);
