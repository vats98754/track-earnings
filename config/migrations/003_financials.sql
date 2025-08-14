-- Financial statements and related data

CREATE TABLE IF NOT EXISTS company_profile (
  ticker TEXT PRIMARY KEY,
  name TEXT,
  sector TEXT,
  industry TEXT,
  exchange TEXT,
  country TEXT,
  market_cap DOUBLE PRECISION,
  shares_outstanding DOUBLE PRECISION,
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS income_statements (
  ticker TEXT NOT NULL,
  period_end DATE NOT NULL,
  period TEXT NOT NULL, -- 'annual' or 'quarter'
  revenue DOUBLE PRECISION,
  cost_of_revenue DOUBLE PRECISION,
  gross_profit DOUBLE PRECISION,
  operating_income DOUBLE PRECISION,
  net_income DOUBLE PRECISION,
  eps DOUBLE PRECISION,
  ebitda DOUBLE PRECISION,
  provider TEXT NOT NULL,
  PRIMARY KEY (ticker, period_end, period)
);

CREATE TABLE IF NOT EXISTS balance_sheets (
  ticker TEXT NOT NULL,
  period_end DATE NOT NULL,
  period TEXT NOT NULL,
  total_assets DOUBLE PRECISION,
  total_liabilities DOUBLE PRECISION,
  total_debt DOUBLE PRECISION,
  total_equity DOUBLE PRECISION,
  cash_and_equivalents DOUBLE PRECISION,
  current_assets DOUBLE PRECISION,
  current_liabilities DOUBLE PRECISION,
  provider TEXT NOT NULL,
  PRIMARY KEY (ticker, period_end, period)
);

CREATE TABLE IF NOT EXISTS cash_flows (
  ticker TEXT NOT NULL,
  period_end DATE NOT NULL,
  period TEXT NOT NULL,
  operating_cash_flow DOUBLE PRECISION,
  investing_cash_flow DOUBLE PRECISION,
  financing_cash_flow DOUBLE PRECISION,
  free_cash_flow DOUBLE PRECISION,
  dividends_paid DOUBLE PRECISION,
  provider TEXT NOT NULL,
  PRIMARY KEY (ticker, period_end, period)
);

CREATE TABLE IF NOT EXISTS dividends (
  ticker TEXT NOT NULL,
  ex_date DATE NOT NULL,
  amount DOUBLE PRECISION NOT NULL,
  PRIMARY KEY (ticker, ex_date)
);

-- Sector peer comps: store rolling medians per sector/metric
CREATE TABLE IF NOT EXISTS sector_medians (
  as_of DATE NOT NULL,
  sector TEXT NOT NULL,
  metric TEXT NOT NULL, -- 'pe','pb','ev_ebitda','ps','fcy'
  value DOUBLE PRECISION,
  PRIMARY KEY (as_of, sector, metric)
);

-- Analyst estimates (minimal)
CREATE TABLE IF NOT EXISTS analyst_estimates (
  ticker TEXT NOT NULL,
  as_of DATE NOT NULL,
  next_year_eps DOUBLE PRECISION,
  current_year_eps DOUBLE PRECISION,
  rating_buy INTEGER,
  rating_hold INTEGER,
  rating_sell INTEGER,
  target_price DOUBLE PRECISION,
  provider TEXT NOT NULL,
  PRIMARY KEY (ticker, as_of)
);

-- Insider trades (minimal)
CREATE TABLE IF NOT EXISTS insider_trades (
  ticker TEXT NOT NULL,
  filed_date DATE NOT NULL,
  insider TEXT,
  transaction_type TEXT,
  shares BIGINT,
  price DOUBLE PRECISION,
  total_value DOUBLE PRECISION,
  PRIMARY KEY (ticker, filed_date, insider, transaction_type)
);

-- Filings (minimal)
CREATE TABLE IF NOT EXISTS filings (
  ticker TEXT,
  filed_date DATE NOT NULL,
  form TEXT NOT NULL,
  title TEXT,
  url TEXT,
  PRIMARY KEY (filed_date, form, url)
);

-- Earnings transcripts (minimal, store link and summary)
CREATE TABLE IF NOT EXISTS transcripts (
  ticker TEXT NOT NULL,
  event_date DATE NOT NULL,
  title TEXT,
  url TEXT,
  sentiment DOUBLE PRECISION,
  PRIMARY KEY (ticker, event_date, url)
);

-- Fundamentals computed metrics snapshots (TTM/point-in-time)
CREATE TABLE IF NOT EXISTS fundamentals_metrics (
  ticker TEXT NOT NULL,
  asof DATE NOT NULL,
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
