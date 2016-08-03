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
	return fmt.Sprintf("Title: %s", m.Title)
}

type Param struct {
	key, value string
	meta
}

func (p Param) String() string {
	return fmt.Sprintf("%s=%s", p.key, p.value)
}

type meta struct {
	required bool
	Query    QueryFunc
}

func IDParam(value string) (*Param, error) {
	if value == "" {
		return nil, fmt.Errorf("id is required")
	}
	p := Param{"i", value, meta{true, QueryByID}}
	return &p, nil
}

func TitleParam(value string) (*Param, error) {
	if value == "" {
		return nil, fmt.Errorf("title is required")
	}
	p := Param{"t", value, meta{true, QueryByTitle}}
	return &p, nil
}

func YearParam(value string) (*Param, error) {
	if value != "" {
		if _, err := strconv.Atoi(value); err != nil {
			return nil, fmt.Errorf("invalid year: %q", value)
		}
	}
	p := Param{"y", value, meta{}}
	return &p, nil
}

// QueryFunc type
type QueryFunc func(*Param, ...*Param) (*Movie, error)

// QueryByTitle omdb by title
func QueryByTitle(title *Param, extra ...*Param) (*Movie, error) {
	if title.key != "t" {
		return nil, fmt.Errorf("expected title param, got %q", title.key)
	}
	q := BaseURL + fmt.Sprintf("&%s=%s", title.key, url.QueryEscape(title.value))
	for _, p := range extra {
		if p.value == "" {
			continue
		}
		if p.required {
			return nil, fmt.Errorf("extra param %q must not be required", p.key)
		}
		q += fmt.Sprintf("&%s=%s", p.key, url.QueryEscape(p.value))
	}

	return get(q)
}

// QueryByID omdb by title
func QueryByID(id *Param, extra ...*Param) (*Movie, error) {
	if id.key != "i" {
		return nil, fmt.Errorf("expected id param, got %q", id.key)
	}
	q := BaseURL + fmt.Sprintf("&%s=%s", id.key, url.QueryEscape(id.value))
	for _, p := range extra {
		if p.value == "" {
			continue
		}
		if p.required {
			return nil, fmt.Errorf("extra param %q must not be required", p.key)
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
