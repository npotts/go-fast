# Go-fast  [![GoDoc](https://godoc.org/github.com/npotts/go-fast?status.svg)](https://godoc.org/github.com/npotts/go-fast) [![Build Status](https://travis-ci.org/npotts/go-fast.svg?branch=master)](https://travis-ci.org/npotts/go-fast) [![Go Report Card](https://goreportcard.com/badge/github.com/npotts/go-fast)](https://goreportcard.com/report/github.com/npotts/go-fast)
go-fast is a golang module for interacting with www.fast.com's API to extract believable speed benchmarks.

# Install

If you are familiar with Go's toolchain this is pretty much self-explanatory, otherwise you need to setup the ```$GOPATH``` enviromental variable

```sh
❯❯❯ export GOPATH=${HOME}/go
❯❯❯ go get -u github.com/npotts/go-fast/go-fast
❯❯❯ ${GOPATH}/bin/go-fast --help
usage: go-fast [<flags>]

A CLI interface to www.fast.com

Flags:
      --help        Show context-sensitive help (also try --help-long and --help-man).
  -w, --workers=3   Number of workers to start. Currently www.fast.com/Netflix only allows up to 3
  -m, --max="50MB"  Maximum worker download size

exit status 1
❯❯❯ go-fast
2016/07/20 21:46:45 Starting with 3 worker(s)
2016/07/20 21:47:43 3 worker(s) downloaded at an average of
  2809750.97 Bps
  2743.90 Kbps
  2.68 Mbps
```

# TODO
- Modify ```token.go``` to fill up to the requested number of URLs rather than the maximum of 3 the API defaults to