package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopl/ch4/poster/omdb"
)

var title = flag.String("t", "", "Movie title to search for")
var year = flag.String("y", "", "Year of release")

func main() {
	flag.Parse()
	if *title == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *year != "" {
		if _, err := strconv.Atoi(*year); err != nil {
			fmt.Fprintf(os.Stderr, "poster: invalid year: %q\n", *year)
			os.Exit(1)
		}
	}

	filename, n, err := fetchPoster(*title, *year)
	if err != nil {
		fmt.Fprintf(os.Stderr, "poster: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s => %s (%d bytes).\n", *title, filename, n)
}

func fetchPoster(title, year string) (filename string, n int64, err error) {
	movie, err := omdb.Query(title, year)
	if err != nil {
		return "", 0, err
	}
	if movie.Poster == "" || movie.Poster == "N/A" {
		return "", 0, fmt.Errorf("the movie has no poster")
	}
	resp, err := http.Get(movie.Poster)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	dot := strings.LastIndex(movie.Poster, ".")
	filename = movie.Title + movie.Poster[dot:]
	f, err := os.Create(filename)
	if err != nil {
		return "", 0, err
	}
	n, err = io.Copy(f, resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return
}
