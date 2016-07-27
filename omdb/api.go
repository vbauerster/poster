package omdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SearchURL is a base query url
const SearchURL = "http://www.omdbapi.com/"

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

// Query omdb by title and year
func Query(title, year string) (*Movie, error) {
	q := SearchURL + "?t=" + url.QueryEscape(title) + "&y=" + year
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
