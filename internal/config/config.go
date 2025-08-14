package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Providers struct {
		AlphaVantage Provider `yaml:"alphaVantage"`
		FMP          Provider `yaml:"fmp"`
		Finnhub      Provider `yaml:"finnhub"`
		FRED         Provider `yaml:"fred"`
		Marketaux    Provider `yaml:"marketaux"`
	} `yaml:"providers"`
	Storage struct {
		Postgres struct {
			Enabled bool   `yaml:"enabled"`
			DSN     string `yaml:"dsn"`
		} `yaml:"postgres"`
		Redis struct {
			Enabled bool   `yaml:"enabled"`
			Addr    string `yaml:"addr"`
			DB      int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"storage"`
	Schedules struct {
		PricesCron string `yaml:"pricesCron"`
		DailyCron  string `yaml:"dailyCron"`
		MacroCron  string `yaml:"macroCron"`
	} `yaml:"schedules"`
	Universe struct {
		Tickers []string `yaml:"tickers"`
	} `yaml:"universe"`
}

type Provider struct {
	Enabled bool   `yaml:"enabled"`
	APIKey  string `yaml:"apiKey"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &c, nil
}
