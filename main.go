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

func main() {
	flag.Parse()

	required, err := omdb.IDParam(*id)

	if err != nil {
		required, err = omdb.TitleParam(*title)
		if err != nil {
			fmt.Fprintf(os.Stderr, "poster: either -t or -i required\n")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}
	fmt.Printf("%#v\n", required)

	yearParam, err := omdb.YearParam(*year)
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
	fmt.Printf("%s => %s (%d bytes).\n", *title, filename, n)
}

// func fetchByID(id, year string) (filename string, n int64, err error) {
// 	idParam, err := omdb.IDParam(id)
// 	if err != nil {
// 		return "", 0, err
// 	}
// 	yearParam, err := omdb.YearParam(year)
// 	if err != nil {
// 		return "", 0, err
// 	}
// 	return fetch(omdb.QueryByTitle, titleParam, yearParam)
// }

// func fetchByTitle(title, year string) (filename string, n int64, err error) {
// 	titleParam, err := omdb.TitleParam(title)
// 	if err != nil {
// 		return "", 0, err
// 	}
// 	yearParam, err := omdb.YearParam(year)
// 	if err != nil {
// 		return "", 0, err
// 	}
// 	return fetch(omdb.QueryByTitle, titleParam, yearParam)
// }

func fetch(required *omdb.Param, extra ...*omdb.Param) (filename string, n int64, err error) {
	movie, err := required.Query(required, extra...)
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
