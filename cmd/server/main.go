package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/vats98754/track-earnings/internal/httpapi"
	"github.com/vats98754/track-earnings/internal/ingest"
	"github.com/vats98754/track-earnings/internal/providers"
	"github.com/vats98754/track-earnings/internal/providers/alphavantage"
	"github.com/vats98754/track-earnings/internal/providers/fmp"
	"github.com/vats98754/track-earnings/internal/providers/fred"
	"github.com/vats98754/track-earnings/internal/providers/mock"
	"github.com/vats98754/track-earnings/internal/repo"
	"github.com/vats98754/track-earnings/internal/scheduler"
	"github.com/vats98754/track-earnings/internal/services"
	"github.com/vats98754/track-earnings/internal/storage"
	"github.com/vats98754/track-earnings/internal/universe"
)

func main() {
	_ = godotenv.Load()

	zerolog.TimeFieldFormat = time.RFC3339
	if os.Getenv("ENV") == "dev" {
		zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822})
	}

	var priceProv providers.PriceProvider
	var fundProv providers.FundamentalsProvider
	var newsProv providers.NewsProvider
	var analystProv ingest.AnalystProvider
	var finProv ingest.FinancialsProvider
	var transcriptsProv ingest.TranscriptsProvider
	var fmpClient *fmp.Client
	providerName := "mock"

	if key := os.Getenv("ALPHA_VANTAGE_KEY"); key != "" {
		av := alphavantage.New(key)
		priceProv = av
		providerName = "alphavantage"
		zlog.Info().Msg("using Alpha Vantage for prices")
	}
	if key := os.Getenv("FMP_KEY"); key != "" {
		fmpClient = fmp.New(key)
		fundProv = fmpClient
		newsProv = fmpClient
		analystProv = fmpClient
		finProv = fmpClient
		transcriptsProv = fmpClient
		if providerName == "mock" {
			providerName = "fmp"
		}
		zlog.Info().Msg("using FMP for fundamentals, news, analyst, financials, transcripts")
	}

	if priceProv == nil || fundProv == nil {
		m := mock.New()
		if priceProv == nil {
			priceProv = m
		}
		if fundProv == nil {
			fundProv = m
		}
		if providerName == "" {
			providerName = "mock"
		}
		zlog.Info().Msg("using mock provider for missing services")
	}

	// Optional Postgres
	var rpo *repo.Repo
	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		pg, err := storage.NewPostgres(dsn)
		if err != nil {
			zlog.Fatal().Err(err).Msg("postgres connect failed")
		}
		rpo = repo.New(pg.DB)
		zlog.Info().Msg("postgres connected")
	}

	// Optional Redis
	var red *storage.Redis
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		db := 0
		if s := os.Getenv("REDIS_DB"); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				db = v
			}
		}
		red = storage.NewRedis(addr, db)
		zlog.Info().Msg("redis client initialized")
	}

	// Universe registry and selection
	reg := universe.NewRegistry()
	cwd, _ := os.Getwd()
	sp, _ := universe.LoadFromFile("sp500", filepath.Join(cwd, "config/universe/sp500.txt"))
	r1k, _ := universe.LoadFromFile("russell1000", filepath.Join(cwd, "config/universe/russell1000.txt"))
	r2k, _ := universe.LoadFromFile("russell2000", filepath.Join(cwd, "config/universe/russell2000.txt"))
	bm, _ := universe.LoadFromFile("benchmarks", filepath.Join(cwd, "config/universe/benchmarks.txt"))
	reg.Add(sp)
	reg.Add(r1k)
	reg.Add(r2k)
	reg.Add(bm)

	// Scheduler
	if rpo != nil {
		// Prices
		var tickers []string
		group := strings.ToLower(os.Getenv("UNIVERSE_GROUP"))
		if s := os.Getenv("TICKERS"); s != "" {
			tickers = strings.Split(s, ",")
		} else if u, ok := reg.Get(group); ok {
			tickers = u.Tickers
		} else {
			tickers = []string{"AAPL", "MSFT", "GOOGL"}
		}
		if b := strings.ToLower(os.Getenv("INGEST_BENCHMARKS")); b == "1" || b == "true" || b == "yes" {
			if u, ok := reg.Get("benchmarks"); ok {
				tickers = append(tickers, u.Tickers...)
			}
		}
		interval := time.Minute
		if s := os.Getenv("PRICES_INTERVAL"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v > 0 {
				interval = time.Duration(v) * time.Second
			}
		}
		sch := scheduler.New()
		sch.Every("prices", interval, ingest.PricesJob{Tickers: tickers, Provider: priceProv, ProviderName: providerName, Repo: rpo, Redis: red}.Run)

		// Fundamentals daily
		if fundProv != nil {
			sch.Every("fundamentals-daily", 24*time.Hour, ingest.FundamentalsJob{Tickers: tickers, Provider: fundProv, Repo: rpo, ProviderName: providerName}.Run)
		}
		// News hourly
		if newsProv != nil {
			sch.Every("news-hourly", time.Hour, ingest.NewsJob{Tickers: tickers, Provider: newsProv, Repo: rpo, Source: "fmp"}.Run)
		}
		// Analyst hourly
		if analystProv != nil {
			sch.Every("analyst-hourly", time.Hour, ingest.AnalystJob{Tickers: tickers, Provider: analystProv, Repo: rpo}.Run)
		}
		// Transcripts daily
		if transcriptsProv != nil {
			sch.Every("transcripts-daily", 24*time.Hour, ingest.TranscriptsJob{Tickers: tickers, Provider: transcriptsProv, Repo: rpo}.Run)
		}
		// Financials backfill daily (rate-limited inside)
		if finProv != nil {
			sch.Every("financials-backfill-daily", 24*time.Hour, ingest.FinancialsHistoryJob{Tickers: tickers, Provider: finProv, Repo: rpo}.Run)
		}
		// Sector medians daily
		sch.Every("sector-medians-daily", 24*time.Hour, ingest.SectorMediansJob{Repo: rpo}.Run)
		// Breadth hourly
		sch.Every("breadth-hourly", time.Hour, ingest.BreadthJob{Repo: rpo, Universe: tickers}.Run)
		// Corporate actions daily
		if fmpClient != nil {
			sch.Every("corpactions-daily", 24*time.Hour, ingest.CorporateActionsJob{Tickers: tickers, Provider: fmpClient, Repo: rpo}.Run)
		}
		// Insider & ownership daily
		if fmpClient != nil {
			sch.Every("insider-ownership-daily", 24*time.Hour, ingest.InsiderOwnershipJob{Tickers: tickers, Provider: fmpClient, Repo: rpo}.Run)
		}
		// Alerts hourly
		sch.Every("alerts-hourly", time.Hour, ingest.AlertsJob{Repo: rpo}.Run)
		// Macro daily
		if fredKey := os.Getenv("FRED_KEY"); fredKey != "" {
			f := fred.New(fredKey)
			mseries := []string{"DGS10", "CPIAUCSL", "UNRATE", "GDP", "DCOILWTICO", "GOLDAMGBD228NLBM", "DEXUSEU", "DEXJPUS", "VIXCLS"}
			sch.Every("macro-daily", 24*time.Hour, ingest.MacroJob{
				Series:  mseries,
				Provider: f,
				Sink: func(ctx context.Context, series string, at time.Time, value float64) error {
					return rpo.UpsertMacroPoint(ctx, series, at, value, "fred")
				},
			}.Run)
		}

		ctx := context.Background()
		go sch.Start(ctx)
	}

	vh := httpapi.ValuationHandler{Service: services.ValuationService{Prices: priceProv, Fund: fundProv, Repo: rpo}}
	uh := httpapi.UniverseHandler{Reg: reg}
	th := httpapi.TechnicalsHandler{Repo: rpo}
	nh := httpapi.NewsHandler{Repo: rpo}
	r := httpapi.NewRouter(httpapi.Services{Valuation: vh, Universes: uh, Technicals: th, News: nh})

	srv := &http.Server{Addr: getEnv("ADDR", ":8080"), Handler: r}
	zlog.Info().Str("addr", srv.Addr).Msg("starting server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zlog.Fatal().Err(err).Msg("server error")
	}
	_ = srv.Shutdown(context.Background())
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
