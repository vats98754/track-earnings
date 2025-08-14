-- Key fundamentals metrics (TTM or latest)
CREATE TABLE IF NOT EXISTS fundamentals_metrics (
  ticker TEXT NOT NULL,
  asof   TIMESTAMPTZ NOT NULL,
  eps DOUBLE PRECISION,
  book_value_ps DOUBLE PRECISION,
  sales_ps DOUBLE PRECISION,
  ebitda DOUBLE PRECISION,
  growth_eps DOUBLE PRECISION,
  roe DOUBLE PRECISION,
  roic DOUBLE PRECISION,
  fcf DOUBLE PRECISION,
  market_cap DOUBLE PRECISION,
  debt DOUBLE PRECISION,
  equity DOUBLE PRECISION,
  provider TEXT NOT NULL,
  PRIMARY KEY (ticker, asof)
);
CREATE INDEX IF NOT EXISTS idx_fm_ticker_asof ON fundamentals_metrics (ticker, asof DESC);
