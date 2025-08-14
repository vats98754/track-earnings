# Track Earnings - Comprehensive Stock Valuation Framework (Go)

A modular Go service to fetch, store, and analyze fundamentals, prices, macro factors, and qualitative signals for over/undervaluation.

Features
- Pluggable data providers (Alpha Vantage, FMP, Finnhub, SEC EDGAR, FRED, Marketaux, etc.) via interfaces
- Schedulers for intraday/daily/weekly jobs
- Storage abstractions for Postgres/TimescaleDB and Redis/RedisTimeSeries
- Compute valuation ratios (P/E, P/B, EV/EBITDA, FCF Yield, PEG, ROE, ROIC)
- Technical indicators (RSI, SMA/EMA, MACD) using TA library or computed in-house
- Sentiment ingestion (news/social) and transcript parsing hooks
- Alerting and Valuation Engine with rules (e.g., undervalued, overbought)
- REST API for valuation reports and live signals

Quick start
1. Install Go 1.22+
2. Copy `.env.example` to `.env` and set your API keys (optional). By default, the server uses a mock provider.
3. Start dependencies (TimescaleDB + Redis):
   - docker compose up -d
   - ./scripts/migrate.sh
4. Run server:
   - make dev
5. Try endpoints:
   - curl http://localhost:8080/health
   - curl http://localhost:8080/api/valuation/AAPL
   - curl http://localhost:8080/api/technicals/AAPL
   - curl http://localhost:8080/api/universes
   - curl http://localhost:8080/api/universe/sp500

Universes and Benchmarks
- Predefined groups in `config/universe/`:
  - sp500.txt, russell1000.txt, russell2000.txt, benchmarks.txt (SPY/QQQ/IWM and sector ETFs XLK, XLE, XLF, etc.)
- Select ingest group via env `UNIVERSE_GROUP` (e.g., `sp500`, `russell1000`).
- Append benchmarks to ingestion with `INGEST_BENCHMARKS=true`.
- Add new groups by adding a new `.txt` file and registering via code or replacing an existing file.

Configuration (.env)
- ADDR=:8080, ENV=dev
- ALPHA_VANTAGE_KEY= (prices)
- FMP_KEY= (fundamentals)
- POSTGRES_DSN=postgres://user:pass@localhost:5432/valuation?sslmode=disable
- REDIS_ADDR=localhost:6379, REDIS_DB=0
- UNIVERSE_GROUP=sp500
- TICKERS=AAPL,MSFT,GOOGL (overrides group)
- PRICES_INTERVAL=60
- INGEST_BENCHMARKS=true

Notes
- Docker uses TimescaleDB image; migrations enable Timescale hypertables for prices and macro.
- Macro ingestion from FRED enabled when FRED_KEY is set (series: DGS10, CPIAUCSL, UNRATE, GDP).
- SEC EDGAR client stub included for future filings ingestion.
- Keep API keys in `.env` (gitignored). Do not commit secrets to `.env.example`.
