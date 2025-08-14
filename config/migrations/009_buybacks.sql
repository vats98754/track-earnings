-- Buybacks table
CREATE TABLE IF NOT EXISTS buybacks (
  ticker TEXT NOT NULL,
  period_end DATE NOT NULL,
  period TEXT NOT NULL,
  amount DOUBLE PRECISION,
  PRIMARY KEY (ticker, period_end, period)
);
