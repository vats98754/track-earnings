package fred

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/vats98754/track-earnings/internal/providers"
)

type Client struct {
	Key  string
	HTTP *http.Client
}

func New(key string) *Client {
	return &Client{Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) GetSeries(series string) ([]providers.MacroPoint, error) {
	q := url.Values{}
	q.Set("series_id", series)
	q.Set("api_key", c.Key)
	q.Set("file_type", "json")
	u := url.URL{Scheme: "https", Host: "api.stlouisfed.org", Path: "/fred/series/observations", RawQuery: q.Encode()}
	resp, err := c.HTTP.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("fred http %d", resp.StatusCode)
	}
	var raw struct {
		Observations []struct {
			Date  string `json:"date"`
			Value string `json:"value"`
		} `json:"observations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	out := make([]providers.MacroPoint, 0, len(raw.Observations))
	for _, o := range raw.Observations {
		var v float64
		_, _ = fmt.Sscanf(o.Value, "%f", &v)
		t, _ := time.Parse("2006-01-02", o.Date)
		out = append(out, providers.MacroPoint{Series: series, Time: t.Unix(), Value: v})
	}
	return out, nil
}
