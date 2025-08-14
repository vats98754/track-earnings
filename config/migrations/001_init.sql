-- prices time-series
CREATE TABLE IF NOT EXISTS prices (
  ticker TEXT NOT NULL,
  time   TIMESTAMPTZ NOT NULL,
  open   DOUBLE PRECISION,
  high   DOUBLE PRECISION,
  low    DOUBLE PRECISION,
  close  DOUBLE PRECISION,
  volume DOUBLE PRECISION,
  provider TEXT NOT NULL,
  PRIMARY KEY (ticker, time)
);
CREATE INDEX IF NOT EXISTS idx_prices_ticker_time ON prices (ticker, time DESC);

-- fundamentals snapshots (by period end)
CREATE TABLE IF NOT EXISTS fundamentals (
  ticker TEXT NOT NULL,
  period_end DATE NOT NULL,
  revenue DOUBLE PRECISION,
  net_income DOUBLE PRECISION,
  eps DOUBLE PRECISION,
  ebitda DOUBLE PRECISION,
  equity DOUBLE PRECISION,
  debt DOUBLE PRECISION,
  shares DOUBLE PRECISION,
  provider TEXT NOT NULL,
  PRIMARY KEY (ticker, period_end)
);

-- macro time-series
CREATE TABLE IF NOT EXISTS macro_points (
  series TEXT NOT NULL,
  time   TIMESTAMPTZ NOT NULL,
  value  DOUBLE PRECISION,
  source TEXT NOT NULL,
  PRIMARY KEY (series, time)
);
CREATE INDEX IF NOT EXISTS idx_macro_series_time ON macro_points (series, time DESC);

-- news
CREATE TABLE IF NOT EXISTS news_items (
  source TEXT NOT NULL,
  ticker TEXT,
  time   TIMESTAMPTZ NOT NULL,
  title  TEXT,
  url    TEXT,
  sentiment DOUBLE PRECISION,
  PRIMARY KEY (source, url)
);
CREATE INDEX IF NOT EXISTS idx_news_ticker_time ON news_items (ticker, time DESC);
