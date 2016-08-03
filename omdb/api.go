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

func (m *Movie) String() string {
	return fmt.Sprintf("ID: %s\nTitle: %s\nPlot: %s\n",
		m.ImdbID, m.Title, m.Plot)
}

type UrlParam struct {
	key, value string
	required   bool
}

func (p UrlParam) String() string {
	return fmt.Sprintf("%s=%s", p.key, p.value)
}

func IDFlag(value string) (*UrlParam, error) {
	if value == "" {
		return nil, fmt.Errorf("id is required")
	}
	p := UrlParam{"i", value, true}
	return &p, nil
}

func TitleFlag(value string) (*UrlParam, error) {
	if value == "" {
		return nil, fmt.Errorf("title is required")
	}
	p := UrlParam{"t", value, true}
	return &p, nil
}

func YearFlag(value string) (*UrlParam, error) {
	if value != "" {
		if _, err := strconv.Atoi(value); err != nil {
			return nil, fmt.Errorf("invalid year: %q", value)
		}
	}
	p := UrlParam{"y", value, false}
	return &p, nil
}

func Query(required *UrlParam, extra ...*UrlParam) (*Movie, error) {
	q := BaseURL + fmt.Sprintf("&%s=%s", required.key, url.QueryEscape(required.value))
	for _, p := range extra {
		if p.value == "" {
			continue
		}
		if p.required {
			return nil, fmt.Errorf("extra param %q must not be required", p)
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
