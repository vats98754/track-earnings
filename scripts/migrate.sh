#!/bin/zsh
set -euo pipefail

: ${POSTGRES_DSN:="postgres://user:pass@localhost:5432/valuation?sslmode=disable"}

if ! command -v psql >/dev/null 2>&1; then
  echo "psql not found. Install PostgreSQL client first." >&2
  exit 1
fi

psql "$POSTGRES_DSN" -f config/migrations/001_init.sql
psql "$POSTGRES_DSN" -f config/migrations/002_timescale.sql
psql "$POSTGRES_DSN" -f config/migrations/003_financials.sql
psql "$POSTGRES_DSN" -f config/migrations/003_fundamentals_metrics.sql
psql "$POSTGRES_DSN" -f config/migrations/004_alter_financials.sql
psql "$POSTGRES_DSN" -f config/migrations/005_earnings_quality.sql
psql "$POSTGRES_DSN" -f config/migrations/006_sector_medians.sql
psql "$POSTGRES_DSN" -f config/migrations/007_corporate_actions.sql
psql "$POSTGRES_DSN" -f config/migrations/008_alter_fm_asof.sql
psql "$POSTGRES_DSN" -f config/migrations/009_buybacks.sql

echo "Migrations applied."
