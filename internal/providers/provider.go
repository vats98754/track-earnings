package providers

type PriceProvider interface {
	GetQuote(ticker string) (PriceQuote, error)
}

type FundamentalsProvider interface {
	GetFundamentals(ticker string) (Fundamentals, error)
}

type MacroProvider interface {
	GetSeries(series string) ([]MacroPoint, error)
}

type NewsProvider interface {
	GetNews(ticker string, limit int) ([]NewsItem, error)
}

type PriceQuote struct {
	Price   float64
	Open    float64
	High    float64
	Low     float64
	Volume  float64
	EV      float64 // optional if available
}

type Fundamentals struct {
	EPS         float64
	BookValuePS float64
	SalesPS     float64
	EBITDA      float64
	GrowthEPS   float64
	ROE         float64
	ROIC        float64
	FCF         float64
	MarketCap   float64
	Debt        float64
	Equity      float64
}

type MacroPoint struct {
	Series string
	Time   int64
	Value  float64
}

type NewsItem struct {
	Time      int64
	Title     string
	URL       string
	Sentiment float64 // -1..+1 if available
}
