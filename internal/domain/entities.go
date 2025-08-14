package domain

import "time"

type Price struct {
	Ticker   string    `db:"ticker" json:"ticker"`
	Time     time.Time `db:"time" json:"time"`
	Open     float64   `db:"open" json:"open"`
	High     float64   `db:"high" json:"high"`
	Low      float64   `db:"low" json:"low"`
	Close    float64   `db:"close" json:"close"`
	Volume   float64   `db:"volume" json:"volume"`
	Provider string    `db:"provider" json:"provider"`
}

type Fundamentals struct {
	Ticker     string    `db:"ticker" json:"ticker"`
	PeriodEnd  time.Time `db:"period_end" json:"period_end"`
	Revenue    float64   `db:"revenue" json:"revenue"`
	NetIncome  float64   `db:"net_income" json:"net_income"`
	EPS        float64   `db:"eps" json:"eps"`
	EBITDA     float64   `db:"ebitda" json:"ebitda"`
	Equity     float64   `db:"equity" json:"equity"`
	Debt       float64   `db:"debt" json:"debt"`
	Shares     float64   `db:"shares" json:"shares"`
	Provider   string    `db:"provider" json:"provider"`
}

type MacroPoint struct {
	Series string    `db:"series" json:"series"`
	Time   time.Time `db:"time" json:"time"`
	Value  float64   `db:"value" json:"value"`
	Source string    `db:"source" json:"source"`
}

type NewsItem struct {
	Source   string    `db:"source" json:"source"`
	Ticker   string    `db:"ticker" json:"ticker"`
	Time     time.Time `db:"time" json:"time"`
	Title    string    `db:"title" json:"title"`
	URL      string    `db:"url" json:"url"`
	Sentiment float64  `db:"sentiment" json:"sentiment"`
}
