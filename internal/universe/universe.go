package universe

import (
	"bufio"
	"os"
	"strings"
)

type Universe struct {
	Name    string
	Tickers []string
}

type Registry struct {
	groups map[string]Universe
}

func NewRegistry() *Registry { return &Registry{groups: map[string]Universe{}} }

func (r *Registry) Add(u Universe) { r.groups[strings.ToLower(u.Name)] = u }

func (r *Registry) Get(name string) (Universe, bool) { u, ok := r.groups[strings.ToLower(name)]; return u, ok }

func (r *Registry) List() []string {
	out := make([]string, 0, len(r.groups))
	for k := range r.groups { out = append(out, k) }
	return out
}

// LoadFromFile loads tickers from a simple newline-delimited list.
func LoadFromFile(name, path string) (Universe, error) {
	f, err := os.Open(path)
	if err != nil { return Universe{}, err }
	defer f.Close()
	s := bufio.NewScanner(f)
	var t []string
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") { continue }
		t = append(t, line)
	}
	return Universe{Name: name, Tickers: t}, s.Err()
}
