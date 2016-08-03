poster
====

Movie poster download command line tool, written in Go.

Posters are downloaded from [omdbapi](http://www.omdbapi.com/)

Usage
-----

```sh
❯ poster -t "The Matrix"
ID: tt0133093
Title: The Matrix
Year: 1999
Plot: A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers.

Saved => The Matrix.jpg (34711 bytes).
```

Install
-------

```sh
go get -u https://github.com/vbauerster/poster
cd $GOPATH/src/github.com/vbauerster/poster
go install
```

License
-------

MIT
