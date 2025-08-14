package edgar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct { HTTP *http.Client }

func New() *Client { return &Client{HTTP: &http.Client{Timeout: 15 * time.Second}} }

// Recent filings for a CIK or ticker using SEC's submissions JSON
func (c *Client) SubmissionsCIK(cik string) (map[string]any, error) {
	u := url.URL{Scheme: "https", Host: "data.sec.gov", Path: "/submissions/CIK" + cik + ".json"}
	resp, err := c.HTTP.Get(u.String())
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode >= 400 { return nil, fmt.Errorf("edgar http %d", resp.StatusCode) }
	var v map[string]any
	return v, json.NewDecoder(resp.Body).Decode(&v)
}
