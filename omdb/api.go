package omdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// BaseURL is a base query url
const BaseURL = "http://www.omdbapi.com/?v=1"

// Movie struct
type Movie struct {
	Title      string
	Year       string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Actors     string
	Plot       string
	Poster     string
	ImdbRating string
	ImdbID     string
	Type       string
}

func (m Movie) String() string {
	return fmt.Sprintf("ID: %s\nTitle: %s\nYear: %s\nPlot: %s\n",
		m.ImdbID, m.Title, m.Year, m.Plot)
}

// URLParam represents url param key value
type URLParam struct {
	key, value string
	required   bool
}

func (p URLParam) String() string {
	return fmt.Sprintf("%s=%s", p.key, p.value)
}

// IDFlag to URLParam
func IDFlag(value string) (*URLParam, error) {
	if value == "" {
		return nil, fmt.Errorf("id is required")
	}
	p := URLParam{"i", value, true}
	return &p, nil
}

// TitleFlag to URLParam
func TitleFlag(value string) (*URLParam, error) {
	if value == "" {
		return nil, fmt.Errorf("title is required")
	}
	p := URLParam{"t", value, true}
	return &p, nil
}

// YearFlag to URLParam
func YearFlag(value string) (*URLParam, error) {
	if value != "" {
		if _, err := strconv.Atoi(value); err != nil {
			return nil, fmt.Errorf("invalid year: %q", value)
		}
	}
	p := URLParam{"y", value, false}
	return &p, nil
}

// PlotFlag to URLParam
func PlotFlag(value string) (*URLParam, error) {
	p := new(URLParam)
	p.key = "plot"
	switch value {
	case "":
		p.value = "short"
		return p, nil
	case "short", "full":
		p.value = value
		return p, nil
	}
	return nil, fmt.Errorf("unrecognized plot: %q", value)
}

// Query func queries omdbapi
func Query(required *URLParam, extra ...*URLParam) (*Movie, error) {
	if !required.required {
		return nil, fmt.Errorf("Query: expected required URLParam, got: %q", required)
	}
	q := BaseURL + fmt.Sprintf("&%s=%s", required.key, url.QueryEscape(required.value))
	for _, p := range extra {
		if p.value == "" {
			continue
		}
		if p.required {
			return nil, fmt.Errorf("Query: expected not required URLParam, got: %q", p)
		}
		q += fmt.Sprintf("&%s=%s", p.key, url.QueryEscape(p.value))
	}

	return get(q)
}

func get(q string) (*Movie, error) {
	resp, err := http.Get(q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting %s: %s", q, resp.Status)
	}

	// Initialize a new Movie value,
	// store a pointer to it in the new variable result
	result := new(Movie)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling json: %v", err)
	}
	return result, nil
}
