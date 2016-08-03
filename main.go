package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"local/poster/omdb"
)

var title = flag.String("t", "", "search by movie title")
var id = flag.String("i", "", "search by movie imdb id (e.g. tt1285016)")
var year = flag.String("y", "", "year of release, optional")

// var plot = flag.String("plot", "", "short or full (short by default)")

func main() {
	flag.Parse()

	required, err := omdb.IDFlag(*id)

	if err != nil {
		required, err = omdb.TitleFlag(*title)
		if err != nil {
			fmt.Fprintf(os.Stderr, "poster: either -t or -i required\n")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	yearParam, err := omdb.YearFlag(*year)
	if err != nil {
		fmt.Fprintf(os.Stderr, "poster: %v\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename, n, err := fetch(required, yearParam)
	if err != nil {
		fmt.Fprintf(os.Stderr, "poster: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved => %s (%d bytes).\n", filename, n)
}

func fetch(required *omdb.UrlParam, extra ...*omdb.UrlParam) (filename string, n int64, err error) {
	movie, err := omdb.Query(required, extra...)
	if err != nil {
		return "", 0, err
	}
	if movie.Poster == "" || movie.Poster == "N/A" {
		return "", 0, fmt.Errorf("the movie has no poster")
	}
	fmt.Println(movie)
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
