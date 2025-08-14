-- Corporate actions and ownership tables

CREATE TABLE IF NOT EXISTS splits (
  ticker TEXT NOT NULL,
  ex_date DATE NOT NULL,
  ratio TEXT,
  numerator DOUBLE PRECISION,
  denominator DOUBLE PRECISION,
  PRIMARY KEY (ticker, ex_date)
);

CREATE TABLE IF NOT EXISTS institutional_ownership (
  ticker TEXT NOT NULL,
  as_of DATE NOT NULL,
  holders INT,
  percent_owned DOUBLE PRECISION,
  PRIMARY KEY (ticker, as_of)
);
